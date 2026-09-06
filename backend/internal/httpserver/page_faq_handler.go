package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// pageFAQHandler manages the FAQ block for a fixed, non-course page
// (masterclasses, gift certificate). Mounted under a path carrying the
// {page} URL param, e.g. r.Route("/page-faq/{page}", handler.routes) — get
// and updateSettings act on the page's one settings row; the item routes
// additionally rely on {id}.
type pageFAQHandler struct {
	svc   *service.PageFAQService
	admin func(http.Handler) http.Handler
}

func newPageFAQHandler(svc *service.PageFAQService, admin func(http.Handler) http.Handler) *pageFAQHandler {
	return &pageFAQHandler{svc: svc, admin: admin}
}

func (h *pageFAQHandler) routes(r chi.Router) {
	r.Get("/", h.get)
	r.With(h.admin).Put("/", h.updateSettings)
	r.Route("/items", func(r chi.Router) {
		r.With(h.admin).Post("/", h.createItem)
		r.With(h.admin).Put("/{id}", h.updateItem)
		r.With(h.admin).Delete("/{id}", h.deleteItem)
	})
}

// get returns the page's settings (title/description/visible) plus its
// items in one response — public, no auth, mirrors the shape
// GetFullBySlug's FAQItems field gives course pages.
func (h *pageFAQHandler) get(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	result, err := h.svc.Get(r.Context(), page)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type pageFAQSettingsRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Visible     bool   `json:"visible"`
}

func (h *pageFAQHandler) updateSettings(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	var req pageFAQSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.UpdateSettings(r.Context(), model.PageFAQSettings{
		Page:        page,
		Title:       req.Title,
		Description: req.Description,
		Visible:     req.Visible,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

type pageFAQItemRequest struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder int    `json:"sortOrder"`
}

func (h *pageFAQHandler) createItem(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	var req pageFAQItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.CreateItem(r.Context(), model.PageFAQItem{
		Page:      page,
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

func (h *pageFAQHandler) updateItem(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req pageFAQItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.UpdateItem(r.Context(), model.PageFAQItem{
		ID:        id,
		Page:      page,
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

func (h *pageFAQHandler) deleteItem(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.DeleteItem(r.Context(), page, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
