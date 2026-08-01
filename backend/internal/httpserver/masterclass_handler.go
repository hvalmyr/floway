package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type masterclassHandler struct {
	svc   *service.MasterclassService
	admin func(http.Handler) http.Handler
}

func newMasterclassHandler(svc *service.MasterclassService, admin func(http.Handler) http.Handler) *masterclassHandler {
	return &masterclassHandler{svc: svc, admin: admin}
}

func (h *masterclassHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	// Public detail lookup by slug, shares the "/{id}" pattern with the
	// admin-only Put/Delete below — same route node, dispatched by method,
	// this handler just treats the captured value as a slug instead of a
	// numeric id.
	r.Get("/{id}", h.getBySlug)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type masterclassRequest struct {
	Slug             string                  `json:"slug"`
	Title            string                  `json:"title"`
	ShortDesc        string                  `json:"shortDescription"`
	FullDesc         string                  `json:"fullDescription"`
	EndingText       string                  `json:"endingText"`
	Duration         string                  `json:"duration"`
	PriceGroup       int                     `json:"priceGroup"`
	PriceIndividual  int                     `json:"priceIndividual"`
	PriceDescription string                  `json:"priceDescription"`
	CoverImage       string                  `json:"coverImage"`
	Status           model.MasterclassStatus `json:"status"`
}

func (h *masterclassHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *masterclassHandler) getBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "id")

	item, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *masterclassHandler) create(w http.ResponseWriter, r *http.Request) {
	var req masterclassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.Masterclass{
		Slug:             req.Slug,
		Title:            req.Title,
		ShortDesc:        req.ShortDesc,
		FullDesc:         req.FullDesc,
		EndingText:       req.EndingText,
		Duration:         req.Duration,
		PriceGroup:       req.PriceGroup,
		PriceIndividual:  req.PriceIndividual,
		PriceDescription: req.PriceDescription,
		CoverImage:       req.CoverImage,
		Status:           req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *masterclassHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req masterclassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.Masterclass{
		ID:               id,
		Slug:             req.Slug,
		Title:            req.Title,
		ShortDesc:        req.ShortDesc,
		FullDesc:         req.FullDesc,
		EndingText:       req.EndingText,
		Duration:         req.Duration,
		PriceGroup:       req.PriceGroup,
		PriceIndividual:  req.PriceIndividual,
		PriceDescription: req.PriceDescription,
		CoverImage:       req.CoverImage,
		Status:           req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *masterclassHandler) delete(w http.ResponseWriter, r *http.Request) {
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
