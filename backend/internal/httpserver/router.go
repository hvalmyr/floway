package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
}

func NewRouter(services Services) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", handleHealth)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/faq", newFAQHandler(services.FAQ).routes)
		r.Route("/teachers", newTeacherHandler(services.Teacher).routes)
		r.Route("/blog-posts", newBlogPostHandler(services.BlogPost).routes)
		r.Route("/masterclasses", newMasterclassHandler(services.Masterclass).routes)
		r.Route("/courses", func(r chi.Router) {
			newCourseHandler(services.Course).routes(r)
			r.Route("/{courseId}/blocks", newCourseBlockHandler(services.CourseBlock).routes)
		})
		r.Route("/course-blocks/{blockId}/lessons", newLessonHandler(services.Lesson).routes)
		r.Route("/leads", newLeadHandler(services.Lead).routes)
	})

	return r
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
