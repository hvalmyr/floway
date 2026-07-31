package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"floway-backend/internal/auth"
	"floway-backend/internal/service"
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

	Tokens         *auth.TokenManager
	SecureCookies  bool
	FrontendOrigin string
}

func NewRouter(services Services) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
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

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/admin", newAuthHandler(services.AdminUser, services.Tokens, services.SecureCookies, admin).routes)
		r.Route("/faq", newFAQHandler(services.FAQ, admin).routes)
		r.Route("/teachers", newTeacherHandler(services.Teacher, admin).routes)
		r.Route("/blog-posts", newBlogPostHandler(services.BlogPost, admin).routes)
		r.Route("/masterclasses", newMasterclassHandler(services.Masterclass, admin).routes)
		r.Route("/courses", func(r chi.Router) {
			newCourseHandler(services.Course, admin).routes(r)
			r.Route("/{courseId}/blocks", newCourseBlockHandler(services.CourseBlock, admin).routes)
		})
		r.Route("/course-blocks/{blockId}/lessons", newLessonHandler(services.Lesson, admin).routes)
		r.Route("/leads", newLeadHandler(services.Lead, admin).routes)
	})

	return r
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
