package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type homeSectionHandler struct {
	svc   *service.HomeSectionService
	admin func(http.Handler) http.Handler
}

func newHomeSectionHandler(svc *service.HomeSectionService, admin func(http.Handler) http.Handler) *homeSectionHandler {
	return &homeSectionHandler{svc: svc, admin: admin}
}

func (h *homeSectionHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Put("/{id}", h.update)
}

type homeSectionRequest struct {
	Visible   bool `json:"visible"`
	SortOrder int  `json:"sortOrder"`
}

// list returns every homepage section (visible or not, admin panel needs
// both) already ordered by sortOrder — the public homepage filters out the
// hidden ones itself.
func (h *homeSectionHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *homeSectionHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req homeSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.HomeSection{
		ID:        id,
		Visible:   req.Visible,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
