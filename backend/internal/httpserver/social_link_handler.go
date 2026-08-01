package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type socialLinkHandler struct {
	svc   *service.SocialLinkService
	admin func(http.Handler) http.Handler
}

func newSocialLinkHandler(svc *service.SocialLinkService, admin func(http.Handler) http.Handler) *socialLinkHandler {
	return &socialLinkHandler{svc: svc, admin: admin}
}

func (h *socialLinkHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type socialLinkRequest struct {
	Label      string `json:"label"`
	Href       string `json:"href"`
	Disclaimer string `json:"disclaimer"`
	SortOrder  int    `json:"sortOrder"`
}

func (h *socialLinkHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *socialLinkHandler) create(w http.ResponseWriter, r *http.Request) {
	var req socialLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.SocialLink{
		Label:      req.Label,
		Href:       req.Href,
		Disclaimer: req.Disclaimer,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *socialLinkHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req socialLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.SocialLink{
		ID:         id,
		Label:      req.Label,
		Href:       req.Href,
		Disclaimer: req.Disclaimer,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *socialLinkHandler) delete(w http.ResponseWriter, r *http.Request) {
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
