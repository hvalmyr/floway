package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type masterclassCategoryHandler struct {
	svc   *service.MasterclassCategoryService
	admin func(http.Handler) http.Handler
}

func newMasterclassCategoryHandler(svc *service.MasterclassCategoryService, admin func(http.Handler) http.Handler) *masterclassCategoryHandler {
	return &masterclassCategoryHandler{svc: svc, admin: admin}
}

func (h *masterclassCategoryHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type masterclassCategoryRequest struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
}

// list returns every category — the public page uses it to render tabs, so
// it has no admin middleware (same convention as object types/features).
func (h *masterclassCategoryHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *masterclassCategoryHandler) create(w http.ResponseWriter, r *http.Request) {
	var req masterclassCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.MasterclassCategory{
		Slug:      req.Slug,
		Name:      req.Name,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *masterclassCategoryHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req masterclassCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.MasterclassCategory{
		ID:        id,
		Slug:      req.Slug,
		Name:      req.Name,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *masterclassCategoryHandler) delete(w http.ResponseWriter, r *http.Request) {
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
