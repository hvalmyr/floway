package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// landingPageHandler serves both the admin-facing flat CRUD (list/create/
// update/delete, no parent scoping — a landing page belongs directly to an
// object type, not to an intermediate section like a course) and the two
// public aggregations (list, full-by-slug).
type landingPageHandler struct {
	svc     *service.LandingPageService
	catalog *service.LandingPageCatalogService
	admin   func(http.Handler) http.Handler
}

func newLandingPageHandler(svc *service.LandingPageService, catalog *service.LandingPageCatalogService, admin func(http.Handler) http.Handler) *landingPageHandler {
	return &landingPageHandler{svc: svc, catalog: catalog, admin: admin}
}

func (h *landingPageHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.Get("/visible", h.listVisible)
	r.Get("/{slug}/full", h.getFullBySlug)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type landingPageRequest struct {
	ObjectTypeID    int64  `json:"objectTypeId"`
	Slug            string `json:"slug"`
	H1              string `json:"h1"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	FAQTitle        string `json:"faqTitle"`
	FAQDescription  string `json:"faqDescription"`
	FAQVisible      bool   `json:"faqVisible"`
	Visible         bool   `json:"visible"`
	SortOrder       int    `json:"sortOrder"`
}

// list returns every landing page, hidden ones included — the admin screen
// (create/copy/rename/hide/delete, п. 7.5 ТЗ). Public callers use GET
// /visible for the homepage/nav instead.
func (h *landingPageHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// listVisible is the public homepage/nav endpoint: every visible landing
// page with its object type.
func (h *landingPageHandler) listVisible(w http.ResponseWriter, r *http.Request) {
	items, err := h.catalog.ListVisible(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// getFullBySlug is the public "landing page" endpoint: one page, its object
// type, its blocks, and its FAQ items.
func (h *landingPageHandler) getFullBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	detail, err := h.catalog.GetFullBySlug(r.Context(), slug)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

func (h *landingPageHandler) create(w http.ResponseWriter, r *http.Request) {
	var req landingPageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.LandingPage{
		ObjectTypeID:    req.ObjectTypeID,
		Slug:            req.Slug,
		H1:              req.H1,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		FAQTitle:        req.FAQTitle,
		FAQDescription:  req.FAQDescription,
		FAQVisible:      req.FAQVisible,
		Visible:         req.Visible,
		SortOrder:       req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *landingPageHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req landingPageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.LandingPage{
		ID:              id,
		ObjectTypeID:    req.ObjectTypeID,
		Slug:            req.Slug,
		H1:              req.H1,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		FAQTitle:        req.FAQTitle,
		FAQDescription:  req.FAQDescription,
		FAQVisible:      req.FAQVisible,
		Visible:         req.Visible,
		SortOrder:       req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *landingPageHandler) delete(w http.ResponseWriter, r *http.Request) {
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
