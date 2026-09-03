package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type iconHandler struct {
	svc   *service.IconService
	admin func(http.Handler) http.Handler
}

func newIconHandler(svc *service.IconService, admin func(http.Handler) http.Handler) *iconHandler {
	return &iconHandler{svc: svc, admin: admin}
}

// No PUT: an icon's SVG is fixed once sanitized and uploaded — replacing it
// is delete-and-reupload, not edit-in-place, same as GalleryPhoto has no
// concept of replacing an image (just remove/re-add).
func (h *iconHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type iconRequest struct {
	Name string `json:"name"`
	SVG  string `json:"svg"`
}

func (h *iconHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *iconHandler) create(w http.ResponseWriter, r *http.Request) {
	var req iconRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.Icon{Name: req.Name, SVG: req.SVG})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *iconHandler) delete(w http.ResponseWriter, r *http.Request) {
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
