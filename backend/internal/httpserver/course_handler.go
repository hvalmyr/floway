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
	svc     *service.CourseService
	catalog *service.CourseCatalogService
	admin   func(http.Handler) http.Handler
}

func newCourseHandler(svc *service.CourseService, catalog *service.CourseCatalogService, admin func(http.Handler) http.Handler) *courseHandler {
	return &courseHandler{svc: svc, catalog: catalog, admin: admin}
}

// routes is expected to be mounted under a path carrying the {sectionId}
// URL param, e.g. r.Route("/course-sections/{sectionId}/courses", handler.routes)
// — the admin-facing nested CRUD. getFullBySlug is mounted separately at
// top-level /courses/{slug}/full (see router.go) since it's public and has
// nothing to do with the section-scoped admin routes.
func (h *courseHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type courseRequest struct {
	Slug           string                        `json:"slug"`
	Name           string                        `json:"name"`
	Description    string                        `json:"description"`
	CoverImage     string                        `json:"coverImage"`
	LessonCount    string                        `json:"lessonCount"`
	TimeLength     string                        `json:"timeLength"`
	Price          string                        `json:"price"`
	DisplayStyle   model.CourseBlockDisplayStyle `json:"displayStyle"`
	Visible        bool                          `json:"visible"`
	SortOrder      int                           `json:"sortOrder"`
	SingleCard     bool                          `json:"singleCard"`
	FAQTitle       string                        `json:"faqTitle"`
	FAQDescription string                        `json:"faqDescription"`
	FAQVisible     bool                          `json:"faqVisible"`
}

func (h *courseHandler) list(w http.ResponseWriter, r *http.Request) {
	sectionID, err := strconv.ParseInt(chi.URLParam(r, "sectionId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	items, err := h.svc.ListBySectionID(r.Context(), sectionID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// getFullBySlug is the public "course page" endpoint: one course, with its
// blocks and each block's lessons. Aggregation lives in CourseCatalogService
// (three queries total, not one-per-block).
func (h *courseHandler) getFullBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	detail, err := h.catalog.GetFullBySlug(r.Context(), slug)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (h *courseHandler) create(w http.ResponseWriter, r *http.Request) {
	sectionID, err := strconv.ParseInt(chi.URLParam(r, "sectionId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req courseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.Course{
		SectionID:      sectionID,
		Slug:           req.Slug,
		Name:           req.Name,
		Description:    req.Description,
		CoverImage:     req.CoverImage,
		LessonCount:    req.LessonCount,
		TimeLength:     req.TimeLength,
		Price:          req.Price,
		DisplayStyle:   req.DisplayStyle,
		Visible:        req.Visible,
		SortOrder:      req.SortOrder,
		SingleCard:     req.SingleCard,
		FAQTitle:       req.FAQTitle,
		FAQDescription: req.FAQDescription,
		FAQVisible:     req.FAQVisible,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *courseHandler) update(w http.ResponseWriter, r *http.Request) {
	sectionID, err := strconv.ParseInt(chi.URLParam(r, "sectionId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

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
		ID:             id,
		SectionID:      sectionID,
		Slug:           req.Slug,
		Name:           req.Name,
		Description:    req.Description,
		CoverImage:     req.CoverImage,
		LessonCount:    req.LessonCount,
		TimeLength:     req.TimeLength,
		Price:          req.Price,
		DisplayStyle:   req.DisplayStyle,
		Visible:        req.Visible,
		SortOrder:      req.SortOrder,
		SingleCard:     req.SingleCard,
		FAQTitle:       req.FAQTitle,
		FAQDescription: req.FAQDescription,
		FAQVisible:     req.FAQVisible,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *courseHandler) delete(w http.ResponseWriter, r *http.Request) {
	sectionID, err := strconv.ParseInt(chi.URLParam(r, "sectionId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Delete(r.Context(), sectionID, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
