package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type giftCertificateCarouselPhotoHandler struct {
	svc   *service.GiftCertificateCarouselPhotoService
	admin func(http.Handler) http.Handler
}

func newGiftCertificateCarouselPhotoHandler(svc *service.GiftCertificateCarouselPhotoService, admin func(http.Handler) http.Handler) *giftCertificateCarouselPhotoHandler {
	return &giftCertificateCarouselPhotoHandler{svc: svc, admin: admin}
}

func (h *giftCertificateCarouselPhotoHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type giftCertificateCarouselPhotoRequest struct {
	Image     string `json:"image"`
	SortOrder int    `json:"sortOrder"`
}

func (h *giftCertificateCarouselPhotoHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *giftCertificateCarouselPhotoHandler) create(w http.ResponseWriter, r *http.Request) {
	var req giftCertificateCarouselPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.GiftCertificateCarouselPhoto{
		Image:     req.Image,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *giftCertificateCarouselPhotoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req giftCertificateCarouselPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.GiftCertificateCarouselPhoto{
		ID:        id,
		Image:     req.Image,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *giftCertificateCarouselPhotoHandler) delete(w http.ResponseWriter, r *http.Request) {
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
