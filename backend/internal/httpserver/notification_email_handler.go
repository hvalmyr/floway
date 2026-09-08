package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// notificationEmailHandler manages the recipient list for new-lead email
// notifications. Unlike most content entities here, there's no public read
// route — this list is never shown on the site, only used server-side by
// notify.EmailNotifier — so every route is admin-gated.
type notificationEmailHandler struct {
	svc   *service.NotificationEmailService
	admin func(http.Handler) http.Handler
}

func newNotificationEmailHandler(svc *service.NotificationEmailService, admin func(http.Handler) http.Handler) *notificationEmailHandler {
	return &notificationEmailHandler{svc: svc, admin: admin}
}

func (h *notificationEmailHandler) routes(r chi.Router) {
	r.Use(h.admin)
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Delete("/{id}", h.delete)
}

type notificationEmailRequest struct {
	Email string `json:"email"`
}

func (h *notificationEmailHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *notificationEmailHandler) create(w http.ResponseWriter, r *http.Request) {
	var req notificationEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.NotificationEmail{Email: req.Email})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *notificationEmailHandler) delete(w http.ResponseWriter, r *http.Request) {
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
