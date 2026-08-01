package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type featureHandler struct {
	svc   *service.FeatureService
	admin func(http.Handler) http.Handler
}

func newFeatureHandler(svc *service.FeatureService, admin func(http.Handler) http.Handler) *featureHandler {
	return &featureHandler{svc: svc, admin: admin}
}

func (h *featureHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type featureRequest struct {
	Page        string `json:"page"`
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
}

// list returns every feature (both pages) by default, for the admin panel.
// Public callers pass ?page=home or ?page=masterclasses to get just the
// ones a given page renders.
func (h *featureHandler) list(w http.ResponseWriter, r *http.Request) {
	var items []model.Feature
	var err error
	if page := r.URL.Query().Get("page"); page != "" {
		items, err = h.svc.ListByPage(r.Context(), page)
	} else {
		items, err = h.svc.List(r.Context())
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *featureHandler) create(w http.ResponseWriter, r *http.Request) {
	var req featureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.Feature{
		Page:        req.Page,
		Icon:        req.Icon,
		Title:       req.Title,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *featureHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req featureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.Feature{
		ID:          id,
		Page:        req.Page,
		Icon:        req.Icon,
		Title:       req.Title,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *featureHandler) delete(w http.ResponseWriter, r *http.Request) {
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
