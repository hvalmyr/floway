package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type faqHandler struct {
	svc *service.FAQService
}

func newFAQHandler(svc *service.FAQService) *faqHandler {
	return &faqHandler{svc: svc}
}

func (h *faqHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
}

type faqCreateRequest struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder int    `json:"sortOrder"`
}

func (h *faqHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *faqHandler) create(w http.ResponseWriter, r *http.Request) {
	var req faqCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), req.Question, req.Answer, req.SortOrder)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *faqHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req faqCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.FAQItem{
		ID:        id,
		Question:  req.Question,
		Answer:    req.Answer,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *faqHandler) delete(w http.ResponseWriter, r *http.Request) {
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

func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrValidation) {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeError(w, http.StatusInternalServerError, err)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
