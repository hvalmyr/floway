package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"

	"floway-backend/internal/auth"
	"floway-backend/internal/service"
	"floway-backend/internal/storage"
)

// maxRequestBodyBytes bounds every request body, not just uploads (which
// already have their own tighter check) — architecture review finding #8:
// nothing stopped a client from sending an arbitrarily large body to any
// JSON endpoint.
const maxRequestBodyBytes = 10 << 20 // 10 MB, matches upload_handler's own cap

func limitBodySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		next.ServeHTTP(w, r)
	})
}

type Services struct {
	FAQ           *service.FAQService
	Teacher       *service.TeacherService
	BlogPost      *service.BlogPostService
	Masterclass   *service.MasterclassService
	CourseSection *service.CourseSectionService
	Course        *service.CourseService
	CourseBlock   *service.CourseBlockService
	Lesson        *service.LessonService
	CourseCatalog *service.CourseCatalogService
	Lead          *service.LeadService
	AdminUser     *service.AdminUserService
	PageContent   *service.PageContentService
	Feature       *service.FeatureService
	AboutItem     *service.AboutItemService
	SocialLink    *service.SocialLinkService
	GalleryPhoto  *service.GalleryPhotoService
	ContentExport *service.ContentExportService

	Storage *storage.Client
	// DB is used only by /readyz. Reaching into the pool directly here (not
	// through a repository) is a deliberate exception — a readiness probe
	// is an infra concern, not domain logic routed through
	// handler->service->repository.
	DB *pgxpool.Pool

	Tokens         *auth.TokenManager
	SecureCookies  bool
	FrontendOrigin string

	// Logger defaults to slog.Default() when nil (tests mostly leave this
	// unset) — production wiring in cmd/api/main.go always sets it to a
	// JSON handler.
	Logger *slog.Logger
}

func NewRouter(services Services) http.Handler {
	r := chi.NewRouter()

	logger := services.Logger
	if logger == nil {
		logger = slog.Default()
	}

	r.Use(middleware.RequestID)
	// Exactly one trusted reverse proxy in front of this service in every
	// deployment (Caddy in prod; the backend port is never published
	// directly). middleware.RealIP is deprecated/spoofable — this variant
	// trusts exactly one XFF hop and never mutates r.RemoteAddr.
	r.Use(middleware.ClientIPFromXFFTrustedProxies(1))
	r.Use(injectLogger(logger))
	r.Use(requestLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(limitBodySize)
	// CSRF: no anti-CSRF token by design, not by omission. The session
	// cookie is SameSite=Lax (never sent on a cross-site POST/PUT/DELETE)
	// plus this CORS policy allows exactly one origin with credentials —
	// together they close the same hole a CSRF token would, without a
	// second token to issue/store/validate. Revisit if a second frontend
	// origin or a SameSite=None requirement ever shows up.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{services.FrontendOrigin},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// /healthz: static liveness — "is the process up" (Recoverer keeps a
	// panic from taking it down entirely). /readyz: real dependency check —
	// "can this instance actually serve traffic" (architecture review
	// finding #6). Docker's HEALTHCHECK and any depends_on/service_healthy
	// should point at /readyz, not /healthz.
	r.Get("/healthz", handleHealth)
	r.Get("/readyz", handleReady(services.DB))

	admin := requireAdminMiddleware(services.Tokens, services.AdminUser)
	uploads := newUploadHandler(services.Storage, admin)
	r.Get("/uploads/{key}", uploads.serve)

	loginLimiter := newIPRateLimiter(rate.Every(time.Minute/5), 5).middleware // 5/min/IP
	leadLimiter := newIPRateLimiter(rate.Every(time.Hour/3), 3).middleware    // 3/hour/IP

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/admin", newAuthHandler(services.AdminUser, services.Tokens, services.SecureCookies, admin, loginLimiter).routes)
		r.Route("/admin/uploads", uploads.adminRoutes)
		r.Route("/faq", newFAQHandler(services.FAQ, admin).routes)
		r.Route("/teachers", newTeacherHandler(services.Teacher, admin).routes)
		r.Route("/blog-posts", newBlogPostHandler(services.BlogPost, services.Tokens, services.AdminUser, admin).routes)
		r.Route("/masterclasses", newMasterclassHandler(services.Masterclass, admin).routes)
		courseHdl := newCourseHandler(services.Course, services.CourseCatalog, admin)
		r.Route("/course-sections", func(r chi.Router) {
			newCourseSectionHandler(services.CourseSection, services.CourseCatalog, admin).routes(r)
			r.Route("/{sectionId}/courses", courseHdl.routes)
		})
		r.Route("/courses", func(r chi.Router) {
			r.Get("/{slug}/full", courseHdl.getFullBySlug)
			r.Route("/{courseId}/blocks", newCourseBlockHandler(services.CourseBlock, admin).routes)
			r.Route("/{courseId}/lessons", newCourseLessonHandler(services.Lesson, admin).routes)
		})
		r.Route("/course-blocks/{blockId}/lessons", newLessonHandler(services.Lesson, admin).routes)
		r.Route("/gallery-photos", newGalleryPhotoHandler(services.GalleryPhoto, admin).routes)
		r.Route("/leads", newLeadHandler(services.Lead, admin, leadLimiter).routes)
		r.Route("/page-content", newPageContentHandler(services.PageContent, admin).routes)
		r.Route("/features", newFeatureHandler(services.Feature, admin).routes)
		r.Route("/about-items", newAboutItemHandler(services.AboutItem, admin).routes)
		r.Route("/social-links", newSocialLinkHandler(services.SocialLink, admin).routes)
		r.Route("/admin/content", newContentExportHandler(services.ContentExport, admin).routes)
	})

	return r
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func handleReady(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "db unreachable"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
