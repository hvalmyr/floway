package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type courseHandler struct {
	svc       *service.CourseService
	blockSvc  *service.CourseBlockService
	lessonSvc *service.LessonService
	admin     func(http.Handler) http.Handler
}

func newCourseHandler(svc *service.CourseService, blockSvc *service.CourseBlockService, lessonSvc *service.LessonService, admin func(http.Handler) http.Handler) *courseHandler {
	return &courseHandler{svc: svc, blockSvc: blockSvc, lessonSvc: lessonSvc, admin: admin}
}

func (h *courseHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.Get("/{slug}/full", h.getFullBySlug)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

// courseModuleResponse embeds CourseBlock so its JSON fields (id, courseId,
// title, ...) stay inline, plus the block's lessons nested under it —
// matches the frontend's CourseModule shape (see app/types/api.ts).
type courseModuleResponse struct {
	model.CourseBlock
	Lessons []model.Lesson `json:"lessons"`
}

// courseDetailResponse embeds Course the same way, plus its blocks (as
// "modules") each with their own nested lessons — matches CourseDetail.
type courseDetailResponse struct {
	model.Course
	Modules []courseModuleResponse `json:"modules"`
}

// getFullBySlug is the public "course page" endpoint: one course, with its
// blocks and each block's lessons, aggregated from three repositories so the
// frontend doesn't have to make a waterfall of requests.
func (h *courseHandler) getFullBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	course, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	blocks, err := h.blockSvc.ListByCourseID(r.Context(), course.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	modules := make([]courseModuleResponse, 0, len(blocks))
	for _, block := range blocks {
		lessons, err := h.lessonSvc.ListByCourseBlockID(r.Context(), block.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		modules = append(modules, courseModuleResponse{CourseBlock: block, Lessons: lessons})
	}

	writeJSON(w, http.StatusOK, courseDetailResponse{Course: course, Modules: modules})
}

type courseRequest struct {
	Slug       string             `json:"slug"`
	Title      string             `json:"title"`
	ShortDesc  string             `json:"shortDescription"`
	FullDesc   string             `json:"fullDescription"`
	Status     model.CourseStatus `json:"status"`
	CoverImage string             `json:"coverImage"`
	Gallery    []string           `json:"gallery"`
	SortOrder  int                `json:"sortOrder"`
}

func (h *courseHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *courseHandler) create(w http.ResponseWriter, r *http.Request) {
	var req courseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.Course{
		Slug:       req.Slug,
		Title:      req.Title,
		ShortDesc:  req.ShortDesc,
		FullDesc:   req.FullDesc,
		Status:     req.Status,
		CoverImage: req.CoverImage,
		Gallery:    req.Gallery,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *courseHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req courseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.Course{
		ID:         id,
		Slug:       req.Slug,
		Title:      req.Title,
		ShortDesc:  req.ShortDesc,
		FullDesc:   req.FullDesc,
		Status:     req.Status,
		CoverImage: req.CoverImage,
		Gallery:    req.Gallery,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *courseHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
