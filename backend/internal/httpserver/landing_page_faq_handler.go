package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// landingPageFAQHandler manages a single landing page's FAQ Q&A items.
// Mounted under a path carrying the {landingPageId} URL param, e.g.
// r.Route("/landing-pages/{landingPageId}/faq-items", handler.routes). The
// block's own title/description/visible flag are plain LandingPage fields,
// edited through landingPageHandler, not here.
type landingPageFAQHandler struct {
	svc   *service.LandingPageFAQService
	admin func(http.Handler) http.Handler
}

func newLandingPageFAQHandler(svc *service.LandingPageFAQService, admin func(http.Handler) http.Handler) *landingPageFAQHandler {
	return &landingPageFAQHandler{svc: svc, admin: admin}
}

func (h *landingPageFAQHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type landingPageFAQCreateRequest struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder int    `json:"sortOrder"`
}

func (h *landingPageFAQHandler) list(w http.ResponseWriter, r *http.Request) {
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

func (h *landingPageFAQHandler) create(w http.ResponseWriter, r *http.Request) {
	landingPageID, err := strconv.ParseInt(chi.URLParam(r, "landingPageId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req landingPageFAQCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.LandingPageFAQItem{
		LandingPageID: landingPageID,
		Question:      req.Question,
		Answer:        req.Answer,
		SortOrder:     req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *landingPageFAQHandler) update(w http.ResponseWriter, r *http.Request) {
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

	var req landingPageFAQCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.LandingPageFAQItem{
		ID:            id,
		LandingPageID: landingPageID,
		Question:      req.Question,
		Answer:        req.Answer,
		SortOrder:     req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *landingPageFAQHandler) delete(w http.ResponseWriter, r *http.Request) {
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
