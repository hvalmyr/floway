package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type leadHandler struct {
	svc     *service.LeadService
	admin   func(http.Handler) http.Handler
	limiter func(http.Handler) http.Handler
}

func newLeadHandler(svc *service.LeadService, admin, limiter func(http.Handler) http.Handler) *leadHandler {
	return &leadHandler{svc: svc, admin: admin, limiter: limiter}
}

func (h *leadHandler) routes(r chi.Router) {
	r.With(h.admin).Get("/", h.list)
	// Static segment ahead of the /{id}/... routes below — safe in chi's
	// trie since there's no bare GET /{id} for leads to collide with.
	r.With(h.admin).Get("/stats/conversion", h.conversionStats)
	// Public and free to hit — rate-limited per IP (architecture review
	// finding #8), otherwise it's an open spam sink.
	r.With(h.limiter).Post("/", h.create)
	r.With(h.admin).Patch("/{id}/status", h.updateStatus)
	r.With(h.admin).Patch("/{id}/review", h.dismissReview)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type leadCreateRequest struct {
	Name          string                `json:"name"`
	Phone         string                `json:"phone"`
	Email         string                `json:"email"`
	ContactMethod model.ContactMethod   `json:"contactMethod"`
	Source        model.LeadSource      `json:"source"`
	RequestType   model.LeadRequestType `json:"requestType"`
	RelatedID     *int64                `json:"relatedId,omitempty"`
	RelatedSlug   string                `json:"relatedSlug"`
}

type leadUpdateStatusRequest struct {
	Status model.LeadStatus `json:"status"`
}

func (h *leadHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *leadHandler) create(w http.ResponseWriter, r *http.Request) {
	var req leadCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.Lead{
		Name:          req.Name,
		Phone:         req.Phone,
		Email:         req.Email,
		ContactMethod: req.ContactMethod,
		Source:        req.Source,
		RequestType:   req.RequestType,
		RelatedID:     req.RelatedID,
		RelatedSlug:   req.RelatedSlug,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *leadHandler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req leadUpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.UpdateStatus(r.Context(), id, req.Status)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *leadHandler) dismissReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.DismissReview(r.Context(), id)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

type conversionStatsResponse struct {
	ClosedWon      int      `json:"closedWon"`
	ClosedLost     int      `json:"closedLost"`
	ConversionRate *float64 `json:"conversionRate"`
}

func (h *leadHandler) conversionStats(w http.ResponseWriter, r *http.Request) {
	won, lost, rate, err := h.svc.ConversionRate(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, conversionStatsResponse{ClosedWon: won, ClosedLost: lost, ConversionRate: rate})
}

func (h *leadHandler) delete(w http.ResponseWriter, r *http.Request) {
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
