package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/service"
)

type siteButtonHandler struct {
	svc   *service.SiteButtonService
	admin func(http.Handler) http.Handler
}

func newSiteButtonHandler(svc *service.SiteButtonService, admin func(http.Handler) http.Handler) *siteButtonHandler {
	return &siteButtonHandler{svc: svc, admin: admin}
}

func (h *siteButtonHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Put("/{key}", h.update)
}

type siteButtonRequest struct {
	Text    string `json:"text"`
	Variant string `json:"variant"`
	URL     string `json:"url"`
}

func (h *siteButtonHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *siteButtonHandler) update(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	var req siteButtonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), key, req.Text, req.Variant, req.URL)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
