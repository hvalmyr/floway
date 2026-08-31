package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type courseSectionHandler struct {
	svc     *service.CourseSectionService
	catalog *service.CourseCatalogService
	admin   func(http.Handler) http.Handler
}

func newCourseSectionHandler(svc *service.CourseSectionService, catalog *service.CourseCatalogService, admin func(http.Handler) http.Handler) *courseSectionHandler {
	return &courseSectionHandler{svc: svc, catalog: catalog, admin: admin}
}

func (h *courseSectionHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.Get("/full", h.listFull)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type courseSectionRequest struct {
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Visible     bool   `json:"visible"`
	SortOrder   int    `json:"sortOrder"`
}

// list returns the flat, unnested section rows — the admin sections table.
func (h *courseSectionHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// listFull returns every section with its courses and each course's blocks
// — the public homepage endpoint.
func (h *courseSectionHandler) listFull(w http.ResponseWriter, r *http.Request) {
	items, err := h.catalog.ListSections(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *courseSectionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req courseSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.CourseSection{
		Heading:     req.Heading,
		Description: req.Description,
		Visible:     req.Visible,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *courseSectionHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req courseSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.CourseSection{
		ID:          id,
		Heading:     req.Heading,
		Description: req.Description,
		Visible:     req.Visible,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *courseSectionHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
