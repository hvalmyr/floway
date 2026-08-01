package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/auth"
	"floway-backend/internal/config"
	"floway-backend/internal/httpserver"
	"floway-backend/internal/repository"
	"floway-backend/internal/service"
	"floway-backend/internal/storage"
)

const adminSessionTTL = 12 * time.Hour

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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

	courseRepo := repository.NewCourseRepository(pool)
	courseBlockRepo := repository.NewCourseBlockRepository(pool)
	lessonRepo := repository.NewLessonRepository(pool)

	services := httpserver.Services{
		FAQ:          service.NewFAQService(repository.NewFAQRepository(pool)),
		Teacher:      service.NewTeacherService(repository.NewTeacherRepository(pool)),
		BlogPost:     service.NewBlogPostService(repository.NewBlogPostRepository(pool)),
		Masterclass:  service.NewMasterclassService(repository.NewMasterclassRepository(pool)),
		Course:       service.NewCourseService(courseRepo),
		CourseDetail: service.NewCourseDetailService(courseRepo, courseBlockRepo, lessonRepo),
		CourseBlock:  service.NewCourseBlockService(courseBlockRepo),
		Lesson:       service.NewLessonService(lessonRepo),
		Lead:         service.NewLeadService(repository.NewLeadRepository(pool)),
		AdminUser:    service.NewAdminUserService(repository.NewAdminUserRepository(pool)),
		PageContent:  service.NewPageContentService(repository.NewPageContentRepository(pool)),
		Feature:      service.NewFeatureService(repository.NewFeatureRepository(pool)),
		AboutItem:    service.NewAboutItemService(repository.NewAboutItemRepository(pool)),
		SocialLink:   service.NewSocialLinkService(repository.NewSocialLinkRepository(pool)),

		Storage: garageClient,
		DB:      pool,

		Tokens:         auth.NewTokenManager(cfg.JWTSecret, adminSessionTTL),
		SecureCookies:  cfg.Env != "local",
		FrontendOrigin: cfg.FrontendOrigin,
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
			log.Printf("graceful shutdown error: %v", err)
		}
		close(shutdownDone)
	}()

	log.Printf("floway-backend listening on :%s (env=%s)", cfg.HTTPPort, cfg.Env)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-shutdownDone
	return nil
}
