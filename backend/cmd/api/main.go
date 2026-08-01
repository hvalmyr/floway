package main

import (
	"context"
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

	garageClient, err := storage.NewClient(cfg.GarageEndpoint, cfg.GarageRegion, cfg.GarageAccessKey, cfg.GarageSecretKey, cfg.GarageBucket)
	if err != nil {
		return err
	}

	services := httpserver.Services{
		FAQ:         service.NewFAQService(repository.NewFAQRepository(pool)),
		Teacher:     service.NewTeacherService(repository.NewTeacherRepository(pool)),
		BlogPost:    service.NewBlogPostService(repository.NewBlogPostRepository(pool)),
		Masterclass: service.NewMasterclassService(repository.NewMasterclassRepository(pool)),
		Course:      service.NewCourseService(repository.NewCourseRepository(pool)),
		CourseBlock: service.NewCourseBlockService(repository.NewCourseBlockRepository(pool)),
		Lesson:      service.NewLessonService(repository.NewLessonRepository(pool)),
		Lead:        service.NewLeadService(repository.NewLeadRepository(pool)),
		AdminUser:   service.NewAdminUserService(repository.NewAdminUserRepository(pool)),
		PageContent: service.NewPageContentService(repository.NewPageContentRepository(pool)),
		Feature:     service.NewFeatureService(repository.NewFeatureRepository(pool)),
		AboutItem:   service.NewAboutItemService(repository.NewAboutItemRepository(pool)),
		SocialLink:  service.NewSocialLinkService(repository.NewSocialLinkRepository(pool)),

		Storage: garageClient,

		Tokens:         auth.NewTokenManager(cfg.JWTSecret, adminSessionTTL),
		SecureCookies:  cfg.Env != "local",
		FrontendOrigin: cfg.FrontendOrigin,
	}

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           httpserver.NewRouter(services),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("floway-backend listening on :%s (env=%s)", cfg.HTTPPort, cfg.Env)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
