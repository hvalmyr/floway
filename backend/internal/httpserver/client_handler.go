package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/service"
)

type clientHandler struct {
	svc   *service.ClientService
	admin func(http.Handler) http.Handler
}

func newClientHandler(svc *service.ClientService, admin func(http.Handler) http.Handler) *clientHandler {
	return &clientHandler{svc: svc, admin: admin}
}

func (h *clientHandler) routes(r chi.Router) {
	// Every client route is admin-only — there's no public surface here.
	r.Use(h.admin)
	r.Get("/{id}", h.detail)
	r.Put("/{id}/tags/product", h.setProductTags)
	r.Put("/{id}/tags/client-type", h.setClientTypeTags)
	r.Post("/{id}/comments", h.addComment)
	r.Post("/{id}/reminders", h.addReminder)
	r.Patch("/{id}/reminders/{reminderId}/complete", h.completeReminder)
}

func (h *clientHandler) detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	detail, err := h.svc.GetDetail(r.Context(), id)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

type setTagsRequest struct {
	TagNames []string `json:"tagNames"`
}

func (h *clientHandler) setProductTags(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var req setTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tags, err := h.svc.SetProductTags(r.Context(), id, req.TagNames)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tags)
}

func (h *clientHandler) setClientTypeTags(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var req setTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tags, err := h.svc.SetClientTypeTags(r.Context(), id, req.TagNames)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tags)
}

type addCommentRequest struct {
	Text string `json:"text"`
}

func (h *clientHandler) addComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var req addCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	comment, err := h.svc.AddComment(r.Context(), id, req.Text)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

type addReminderRequest struct {
	Days int    `json:"days"`
	Note string `json:"note"`
}

func (h *clientHandler) addReminder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var req addReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reminder, err := h.svc.AddReminder(r.Context(), id, req.Days, req.Note)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, reminder)
}

func (h *clientHandler) completeReminder(w http.ResponseWriter, r *http.Request) {
	reminderID, err := strconv.ParseInt(chi.URLParam(r, "reminderId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.CompleteReminder(r.Context(), reminderID); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
