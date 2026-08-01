package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type lessonHandler struct {
	svc   *service.LessonService
	admin func(http.Handler) http.Handler
}

func newLessonHandler(svc *service.LessonService, admin func(http.Handler) http.Handler) *lessonHandler {
	return &lessonHandler{svc: svc, admin: admin}
}

// routes registers lesson routes. It expects to be mounted under a path
// containing a {blockId} URL parameter, e.g. /course-blocks/{blockId}/lessons.
func (h *lessonHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type lessonCreateRequest struct {
	Number        int    `json:"number"`
	Title         string `json:"title"`
	Topics        string `json:"topics"`
	Outcomes      string `json:"outcomes"`
	DurationHours int    `json:"durationHours"`
}

func (h *lessonHandler) list(w http.ResponseWriter, r *http.Request) {
	blockID, err := strconv.ParseInt(chi.URLParam(r, "blockId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	items, err := h.svc.ListByCourseBlockID(r.Context(), blockID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *lessonHandler) create(w http.ResponseWriter, r *http.Request) {
	blockID, err := strconv.ParseInt(chi.URLParam(r, "blockId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req lessonCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.Lesson{
		CourseBlockID: blockID,
		Number:        req.Number,
		Title:         req.Title,
		Topics:        req.Topics,
		Outcomes:      req.Outcomes,
		DurationHours: req.DurationHours,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *lessonHandler) update(w http.ResponseWriter, r *http.Request) {
	blockID, err := strconv.ParseInt(chi.URLParam(r, "blockId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req lessonCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.Lesson{
		ID:            id,
		CourseBlockID: blockID,
		Number:        req.Number,
		Title:         req.Title,
		Topics:        req.Topics,
		Outcomes:      req.Outcomes,
		DurationHours: req.DurationHours,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *lessonHandler) delete(w http.ResponseWriter, r *http.Request) {
	blockID, err := strconv.ParseInt(chi.URLParam(r, "blockId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Delete(r.Context(), blockID, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
