package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type courseBlockHandler struct {
	svc   *service.CourseBlockService
	admin func(http.Handler) http.Handler
}

func newCourseBlockHandler(svc *service.CourseBlockService, admin func(http.Handler) http.Handler) *courseBlockHandler {
	return &courseBlockHandler{svc: svc, admin: admin}
}

// routes is expected to be mounted under a path carrying the {courseId} URL
// param, e.g. r.Route("/courses/{courseId}/blocks", handler.routes).
// list/create rely on {courseId}; update/delete rely on the globally unique
// {id} instead.
func (h *courseBlockHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type courseBlockCreateRequest struct {
	BlockName    string                        `json:"blockName"`
	Description  string                        `json:"description"`
	BlockCover   string                        `json:"blockCover"`
	LessonCount  string                        `json:"lessonCount"`
	TimeLength   string                        `json:"timeLength"`
	Price        string                        `json:"price"`
	DisplayStyle model.CourseBlockDisplayStyle `json:"displayStyle"`
	Visible      bool                          `json:"visible"`
	SortOrder    int                           `json:"sortOrder"`
}

func (h *courseBlockHandler) list(w http.ResponseWriter, r *http.Request) {
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

func (h *courseBlockHandler) create(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req courseBlockCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.CourseBlock{
		CourseID:     courseID,
		BlockName:    req.BlockName,
		Description:  req.Description,
		BlockCover:   req.BlockCover,
		LessonCount:  req.LessonCount,
		TimeLength:   req.TimeLength,
		Price:        req.Price,
		DisplayStyle: req.DisplayStyle,
		Visible:      req.Visible,
		SortOrder:    req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *courseBlockHandler) update(w http.ResponseWriter, r *http.Request) {
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

	var req courseBlockCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.CourseBlock{
		ID:           id,
		CourseID:     courseID,
		BlockName:    req.BlockName,
		Description:  req.Description,
		BlockCover:   req.BlockCover,
		LessonCount:  req.LessonCount,
		TimeLength:   req.TimeLength,
		Price:        req.Price,
		DisplayStyle: req.DisplayStyle,
		Visible:      req.Visible,
		SortOrder:    req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *courseBlockHandler) delete(w http.ResponseWriter, r *http.Request) {
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
