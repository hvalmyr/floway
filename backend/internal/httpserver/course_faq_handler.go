package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// courseFAQHandler manages a single course's FAQ Q&A items. Mounted under a
// path carrying the {courseId} URL param, e.g.
// r.Route("/courses/{courseId}/faq-items", handler.routes) — list/create
// rely on {courseId}; update/delete rely on both {courseId} and {id}. The
// block's own title/description/visible flag are plain Course fields, edited
// through courseHandler, not here.
type courseFAQHandler struct {
	svc   *service.CourseFAQService
	admin func(http.Handler) http.Handler
}

func newCourseFAQHandler(svc *service.CourseFAQService, admin func(http.Handler) http.Handler) *courseFAQHandler {
	return &courseFAQHandler{svc: svc, admin: admin}
}

func (h *courseFAQHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type courseFAQCreateRequest struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder int    `json:"sortOrder"`
}

func (h *courseFAQHandler) list(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	items, err := h.svc.ListByCourseID(r.Context(), courseID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *courseFAQHandler) create(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req courseFAQCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.CourseFAQItem{
		CourseID:  courseID,
		Question:  req.Question,
		Answer:    req.Answer,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *courseFAQHandler) update(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req courseFAQCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.CourseFAQItem{
		ID:        id,
		CourseID:  courseID,
		Question:  req.Question,
		Answer:    req.Answer,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *courseFAQHandler) delete(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Delete(r.Context(), courseID, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
