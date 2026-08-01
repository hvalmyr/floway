package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/service"
)

type pageContentHandler struct {
	svc   *service.PageContentService
	admin func(http.Handler) http.Handler
}

func newPageContentHandler(svc *service.PageContentService, admin func(http.Handler) http.Handler) *pageContentHandler {
	return &pageContentHandler{svc: svc, admin: admin}
}

func (h *pageContentHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Put("/{key}", h.update)
}

type pageContentUpdateRequest struct {
	Value string `json:"value"`
}

func (h *pageContentHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *pageContentHandler) update(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	var req pageContentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), key, req.Value)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
