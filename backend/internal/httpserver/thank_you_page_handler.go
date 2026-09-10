package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// thankYouPageHandler manages the content shown after a lead form
// submission for a fixed variant (course, masterclass, trial_lesson).
// Mounted under a path carrying the {variant} URL param, e.g.
// r.Route("/thank-you-pages/{variant}", handler.routes) — get and
// updateSettings act on the variant's one settings row; the photo/faq-item
// routes additionally rely on {id}. Mirrors pageFAQHandler, extended with a
// second child list (photos).
type thankYouPageHandler struct {
	svc   *service.ThankYouPageService
	admin func(http.Handler) http.Handler
}

func newThankYouPageHandler(svc *service.ThankYouPageService, admin func(http.Handler) http.Handler) *thankYouPageHandler {
	return &thankYouPageHandler{svc: svc, admin: admin}
}

func (h *thankYouPageHandler) routes(r chi.Router) {
	r.Get("/", h.get)
	r.With(h.admin).Put("/", h.updateSettings)
	r.Route("/photos", func(r chi.Router) {
		r.With(h.admin).Post("/", h.createPhoto)
		r.With(h.admin).Put("/{id}", h.updatePhoto)
		r.With(h.admin).Delete("/{id}", h.deletePhoto)
	})
	r.Route("/faq-items", func(r chi.Router) {
		r.With(h.admin).Post("/", h.createFAQItem)
		r.With(h.admin).Put("/{id}", h.updateFAQItem)
		r.With(h.admin).Delete("/{id}", h.deleteFAQItem)
	})
}

// get returns the variant's settings (title/subtitle/description/show-*
// flags) plus its photos and FAQ items in one response — public, no auth.
func (h *thankYouPageHandler) get(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	result, err := h.svc.Get(r.Context(), variant)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type thankYouPageSettingsRequest struct {
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	Description     string `json:"description"`
	HeroImage       string `json:"heroImage"`
	ShowMessengers  bool   `json:"showMessengers"`
	ShowSocialLinks bool   `json:"showSocialLinks"`
	ShowBlogLink    bool   `json:"showBlogLink"`
	BlogLinkText    string `json:"blogLinkText"`
	BlogLinkURL     string `json:"blogLinkUrl"`
	ShowCarousel    bool   `json:"showCarousel"`
	ShowFAQ         bool   `json:"showFaq"`
	ShowCommunity   bool   `json:"showCommunity"`
	CommunityText   string `json:"communityText"`
	CommunityURL    string `json:"communityUrl"`
}

func (h *thankYouPageHandler) updateSettings(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	var req thankYouPageSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.UpdateSettings(r.Context(), model.ThankYouPage{
		Variant:         variant,
		Title:           req.Title,
		Subtitle:        req.Subtitle,
		Description:     req.Description,
		HeroImage:       req.HeroImage,
		ShowMessengers:  req.ShowMessengers,
		ShowSocialLinks: req.ShowSocialLinks,
		ShowBlogLink:    req.ShowBlogLink,
		BlogLinkText:    req.BlogLinkText,
		BlogLinkURL:     req.BlogLinkURL,
		ShowCarousel:    req.ShowCarousel,
		ShowFAQ:         req.ShowFAQ,
		ShowCommunity:   req.ShowCommunity,
		CommunityText:   req.CommunityText,
		CommunityURL:    req.CommunityURL,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

type thankYouPagePhotoRequest struct {
	Image     string `json:"image"`
	SortOrder int    `json:"sortOrder"`
}

func (h *thankYouPageHandler) createPhoto(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	var req thankYouPagePhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.CreatePhoto(r.Context(), model.ThankYouPagePhoto{
		Variant:   variant,
		Image:     req.Image,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *thankYouPageHandler) updatePhoto(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req thankYouPagePhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.UpdatePhoto(r.Context(), model.ThankYouPagePhoto{
		ID:        id,
		Variant:   variant,
		Image:     req.Image,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *thankYouPageHandler) deletePhoto(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.DeletePhoto(r.Context(), variant, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type thankYouPageFAQItemRequest struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder int    `json:"sortOrder"`
}

func (h *thankYouPageHandler) createFAQItem(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	var req thankYouPageFAQItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.CreateFAQItem(r.Context(), model.ThankYouPageFAQItem{
		Variant:   variant,
		Question:  req.Question,
		Answer:    req.Answer,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *thankYouPageHandler) updateFAQItem(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req thankYouPageFAQItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.UpdateFAQItem(r.Context(), model.ThankYouPageFAQItem{
		ID:        id,
		Variant:   variant,
		Question:  req.Question,
		Answer:    req.Answer,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *thankYouPageHandler) deleteFAQItem(w http.ResponseWriter, r *http.Request) {
	variant := chi.URLParam(r, "variant")

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.DeleteFAQItem(r.Context(), variant, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
