package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"floway-backend/internal/auth"
	"floway-backend/internal/service"
	"floway-backend/internal/storage"
)

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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{services.FrontendOrigin},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", handleHealth)

	admin := requireAdminMiddleware(services.Tokens)
	uploads := newUploadHandler(services.Storage, admin)
	r.Get("/uploads/{key}", uploads.serve)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/admin", newAuthHandler(services.AdminUser, services.Tokens, services.SecureCookies, admin).routes)
		r.Route("/admin/uploads", uploads.adminRoutes)
		r.Route("/faq", newFAQHandler(services.FAQ, admin).routes)
		r.Route("/teachers", newTeacherHandler(services.Teacher, admin).routes)
		r.Route("/blog-posts", newBlogPostHandler(services.BlogPost, admin).routes)
		r.Route("/masterclasses", newMasterclassHandler(services.Masterclass, admin).routes)
		r.Route("/courses", func(r chi.Router) {
			newCourseHandler(services.Course, services.CourseBlock, services.Lesson, admin).routes(r)
			r.Route("/{courseId}/blocks", newCourseBlockHandler(services.CourseBlock, admin).routes)
		})
		r.Route("/course-blocks/{blockId}/lessons", newLessonHandler(services.Lesson, admin).routes)
		r.Route("/leads", newLeadHandler(services.Lead, admin).routes)
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
