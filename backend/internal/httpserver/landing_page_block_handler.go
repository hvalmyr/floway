package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// landingPageBlockHandler is mounted under a path carrying the
// {landingPageId} URL param, e.g.
// r.Route("/landing-pages/{landingPageId}/blocks", handler.routes).
type landingPageBlockHandler struct {
	svc   *service.LandingPageBlockService
	admin func(http.Handler) http.Handler
}

func newLandingPageBlockHandler(svc *service.LandingPageBlockService, admin func(http.Handler) http.Handler) *landingPageBlockHandler {
	return &landingPageBlockHandler{svc: svc, admin: admin}
}

func (h *landingPageBlockHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type landingPageBlockRequest struct {
	Image     string `json:"image"`
	Text      string `json:"text"`
	SortOrder int    `json:"sortOrder"`
}

func (h *landingPageBlockHandler) list(w http.ResponseWriter, r *http.Request) {
	landingPageID, err := strconv.ParseInt(chi.URLParam(r, "landingPageId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	items, err := h.svc.ListByLandingPageID(r.Context(), landingPageID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *landingPageBlockHandler) create(w http.ResponseWriter, r *http.Request) {
	landingPageID, err := strconv.ParseInt(chi.URLParam(r, "landingPageId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req landingPageBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.LandingPageBlock{
		LandingPageID: landingPageID,
		Image:         req.Image,
		Text:          req.Text,
		SortOrder:     req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *landingPageBlockHandler) update(w http.ResponseWriter, r *http.Request) {
	landingPageID, err := strconv.ParseInt(chi.URLParam(r, "landingPageId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req landingPageBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.LandingPageBlock{
		ID:            id,
		LandingPageID: landingPageID,
		Image:         req.Image,
		Text:          req.Text,
		SortOrder:     req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *landingPageBlockHandler) delete(w http.ResponseWriter, r *http.Request) {
	landingPageID, err := strconv.ParseInt(chi.URLParam(r, "landingPageId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Delete(r.Context(), landingPageID, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
