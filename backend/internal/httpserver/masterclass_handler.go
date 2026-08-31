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
	// Public detail lookup by slug, admin mutations by numeric id — same
	// URL shape ("/{something}"), dispatched by method, so the path itself
	// is unchanged for clients; only the chi param name here reflects what
	// each method actually expects, instead of both being called "id".
	r.Get("/{slug}", h.getBySlug)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type masterclassRequest struct {
	Slug         string                  `json:"slug"`
	Title        string                  `json:"title"`
	Description  string                  `json:"description"`
	Description2 string                  `json:"description2"`
	EndingText   string                  `json:"endingText"`
	Duration     string                  `json:"duration"`
	Price        string                  `json:"price"`
	CoverImage   string                  `json:"coverImage"`
	Status       model.MasterclassStatus `json:"status"`
}

func (h *masterclassHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *masterclassHandler) getBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	item, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeServiceError(w, r, err)
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
		Slug:         req.Slug,
		Title:        req.Title,
		Description:  req.Description,
		Description2: req.Description2,
		EndingText:   req.EndingText,
		Duration:     req.Duration,
		Price:        req.Price,
		CoverImage:   req.CoverImage,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, r, err)
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
		ID:           id,
		Slug:         req.Slug,
		Title:        req.Title,
		Description:  req.Description,
		Description2: req.Description2,
		EndingText:   req.EndingText,
		Duration:     req.Duration,
		Price:        req.Price,
		CoverImage:   req.CoverImage,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, r, err)
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
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
