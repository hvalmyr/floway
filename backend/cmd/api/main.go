package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/auth"
	"floway-backend/internal/config"
	"floway-backend/internal/httpserver"
	"floway-backend/internal/notify"
	"floway-backend/internal/repository"
	"floway-backend/internal/service"
	"floway-backend/internal/storage"
)

const adminSessionTTL = 12 * time.Hour

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	// pgxpool.New never connects (it's a lazy pool) — Ping so a bad
	// DATABASE_URL or an unreachable Postgres fails loudly at startup
	// instead of surfacing as a 500 on the first real request (architecture
	// review finding #6).
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	garageClient, err := storage.NewClient(cfg.GarageEndpoint, cfg.GarageRegion, cfg.GarageAccessKey, cfg.GarageSecretKey, cfg.GarageBucket)
	if err != nil {
		return err
	}
	if err := garageClient.EnsureBucket(ctx); err != nil {
		return fmt.Errorf("connect to garage: %w", err)
	}

	courseSectionRepo := repository.NewCourseSectionRepository(pool)
	courseRepo := repository.NewCourseRepository(pool)
	courseBlockRepo := repository.NewCourseBlockRepository(pool)
	lessonRepo := repository.NewLessonRepository(pool)
	courseFAQRepo := repository.NewCourseFAQRepository(pool)
	masterclassRepo := repository.NewMasterclassRepository(pool)
	leadRepo := repository.NewLeadRepository(pool)
	clientRepo := repository.NewClientRepository(pool)
	productTagRepo := repository.NewProductTagRepository(pool)
	clientTypeTagRepo := repository.NewClientTypeTagRepository(pool)

	// Both channels are optional and independent — SMTP_*/TELEGRAM_* are
	// deliberately absent from config.Load()'s required checks, so an
	// unconfigured channel is silently skipped rather than failing startup.
	var leadNotifyChannels []notify.Channel
	if cfg.SMTPHost != "" && cfg.NotifyEmailTo != "" {
		leadNotifyChannels = append(leadNotifyChannels, notify.NewEmailNotifier(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom, cfg.NotifyEmailTo, cfg.FrontendOrigin+"/admin/leads", cfg.SMTPUser, cfg.SMTPPassword))
	}
	if cfg.TelegramBotToken != "" && cfg.TelegramChatID != "" {
		leadNotifyChannels = append(leadNotifyChannels, notify.NewTelegramNotifier(cfg.TelegramBotToken, cfg.TelegramChatID))
	}
	var leadNotifier service.LeadNotifier
	if len(leadNotifyChannels) > 0 {
		leadNotifier = notify.NewMultiNotifier(leadNotifyChannels...)
	} else {
		logger.Info("lead notifications disabled: no SMTP or Telegram channel configured")
	}

	services := httpserver.Services{
		FAQ:           service.NewFAQService(repository.NewFAQRepository(pool)),
		Teacher:       service.NewTeacherService(repository.NewTeacherRepository(pool)),
		BlogPost:      service.NewBlogPostService(repository.NewBlogPostRepository(pool)),
		Masterclass:   service.NewMasterclassService(masterclassRepo),
		CourseSection: service.NewCourseSectionService(courseSectionRepo),
		Course:        service.NewCourseService(courseRepo),
		CourseBlock:   service.NewCourseBlockService(courseBlockRepo),
		Lesson:        service.NewLessonService(lessonRepo),
		CourseFAQ:     service.NewCourseFAQService(courseFAQRepo),
		CourseCatalog: service.NewCourseCatalogService(courseSectionRepo, courseRepo, courseRepo, courseBlockRepo, courseBlockRepo, lessonRepo, courseFAQRepo),
		Lead:          service.NewLeadService(leadRepo, clientRepo, leadNotifier, courseRepo, masterclassRepo),
		Client: service.NewClientService(
			clientRepo,
			leadRepo,
			productTagRepo,
			clientTypeTagRepo,
			repository.NewClientCommentRepository(pool),
			repository.NewReminderRepository(pool),
		),
		Tag:           service.NewTagService(productTagRepo, clientTypeTagRepo),
		AdminUser:     service.NewAdminUserService(repository.NewAdminUserRepository(pool)),
		PageContent:   service.NewPageContentService(repository.NewPageContentRepository(pool), garageClient),
		Feature:       service.NewFeatureService(repository.NewFeatureRepository(pool)),
		AboutItem:     service.NewAboutItemService(repository.NewAboutItemRepository(pool)),
		SocialLink:    service.NewSocialLinkService(repository.NewSocialLinkRepository(pool)),
		GalleryPhoto:  service.NewGalleryPhotoService(repository.NewGalleryPhotoRepository(pool)),
		Icon:          service.NewIconService(repository.NewIconRepository(pool)),
		ContentExport: service.NewContentExportService(pool, garageClient),

		Storage: garageClient,
		DB:      pool,

		Tokens:         auth.NewTokenManager(cfg.JWTSecret, adminSessionTTL),
		SecureCookies:  cfg.Env != "local",
		FrontendOrigin: cfg.FrontendOrigin,
		Logger:         logger,
	}

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           httpserver.NewRouter(services),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// shutdownDone is closed only once Shutdown() has actually finished
	// draining in-flight requests. Without waiting on it, ListenAndServe()
	// returning ErrServerClosed (which happens as soon as Shutdown closes
	// the listener, not when draining completes) would let run() return
	// immediately and fire the deferred pool.Close() out from under
	// requests still being handled (architecture review finding #7).
	shutdownDone := make(chan struct{})
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown error", "error", err)
		}
		close(shutdownDone)
	}()

	logger.Info("floway-backend listening", "port", cfg.HTTPPort, "env", cfg.Env)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-shutdownDone
	return nil
}
