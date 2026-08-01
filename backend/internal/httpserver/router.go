package httpserver

import (
	"context"
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
	FAQ         *service.FAQService
	Teacher     *service.TeacherService
	BlogPost    *service.BlogPostService
	Masterclass *service.MasterclassService
	Course      *service.CourseService
	CourseBlock *service.CourseBlockService
	Lesson      *service.LessonService
	Lead        *service.LeadService
	AdminUser   *service.AdminUserService
	PageContent *service.PageContentService
	Feature     *service.FeatureService
	AboutItem   *service.AboutItemService
	SocialLink  *service.SocialLinkService

	Storage *storage.Client
	// DB is used only by /readyz. Reaching into the pool directly here (not
	// through a repository) is a deliberate exception — a readiness probe
	// is an infra concern, not domain logic routed through
	// handler->service->repository.
	DB *pgxpool.Pool

	Tokens         *auth.TokenManager
	SecureCookies  bool
	FrontendOrigin string
}

func NewRouter(services Services) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	// Exactly one trusted reverse proxy in front of this service in every
	// deployment (Caddy in prod; the backend port is never published
	// directly). middleware.RealIP is deprecated/spoofable — this variant
	// trusts exactly one XFF hop and never mutates r.RemoteAddr.
	r.Use(middleware.ClientIPFromXFFTrustedProxies(1))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(limitBodySize)
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

	admin := requireAdminMiddleware(services.Tokens)
	uploads := newUploadHandler(services.Storage, admin)
	r.Get("/uploads/{key}", uploads.serve)

	loginLimiter := newIPRateLimiter(rate.Every(time.Minute/5), 5).middleware // 5/min/IP
	leadLimiter := newIPRateLimiter(rate.Every(time.Hour/3), 3).middleware    // 3/hour/IP

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/admin", newAuthHandler(services.AdminUser, services.Tokens, services.SecureCookies, admin, loginLimiter).routes)
		r.Route("/admin/uploads", uploads.adminRoutes)
		r.Route("/faq", newFAQHandler(services.FAQ, admin).routes)
		r.Route("/teachers", newTeacherHandler(services.Teacher, admin).routes)
		r.Route("/blog-posts", newBlogPostHandler(services.BlogPost, services.Tokens, admin).routes)
		r.Route("/masterclasses", newMasterclassHandler(services.Masterclass, admin).routes)
		r.Route("/courses", func(r chi.Router) {
			newCourseHandler(services.Course, services.CourseBlock, services.Lesson, admin).routes(r)
			r.Route("/{courseId}/blocks", newCourseBlockHandler(services.CourseBlock, admin).routes)
		})
		r.Route("/course-blocks/{blockId}/lessons", newLessonHandler(services.Lesson, admin).routes)
		r.Route("/leads", newLeadHandler(services.Lead, admin, leadLimiter).routes)
		r.Route("/page-content", newPageContentHandler(services.PageContent, admin).routes)
		r.Route("/features", newFeatureHandler(services.Feature, admin).routes)
		r.Route("/about-items", newAboutItemHandler(services.AboutItem, admin).routes)
		r.Route("/social-links", newSocialLinkHandler(services.SocialLink, admin).routes)
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
