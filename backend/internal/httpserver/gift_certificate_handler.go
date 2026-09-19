package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/certpdf"
	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// giftCertificateHandler manages issued gift certificates and their PDF
// downloads. Like notificationEmailHandler, there's no public read route —
// every method is admin-gated.
type giftCertificateHandler struct {
	svc   *service.GiftCertificateService
	admin func(http.Handler) http.Handler
}

func newGiftCertificateHandler(svc *service.GiftCertificateService, admin func(http.Handler) http.Handler) *giftCertificateHandler {
	return &giftCertificateHandler{svc: svc, admin: admin}
}

func (h *giftCertificateHandler) routes(r chi.Router) {
	r.Use(h.admin)
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Delete("/{id}", h.delete)
	r.Get("/{id}/pdf", h.pdf)
}

type giftCertificateRequest struct {
	Kind      model.GiftCertificateKind `json:"kind"`
	Value     string                    `json:"value"`
	Recipient string                    `json:"recipient"`
}

func (h *giftCertificateHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *giftCertificateHandler) create(w http.ResponseWriter, r *http.Request) {
	var req giftCertificateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.GiftCertificate{
		Kind:      req.Kind,
		Value:     req.Value,
		Recipient: req.Recipient,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *giftCertificateHandler) delete(w http.ResponseWriter, r *http.Request) {
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

// russianMonths formats a certificate's issue date as "1 сентября 2026" —
// Go's time package has no Russian locale, so this is a small lookup rather
// than pulling in a locale dependency for one string.
var russianMonths = [...]string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func formatRussianDate(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), russianMonths[t.Month()-1], t.Year())
}

func (h *giftCertificateHandler) pdf(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	cert, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}

	data, err := certpdf.Render(cert, formatRussianDate(cert.IssuedAt))
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="certificate-%s.pdf"`, cert.Number))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(data)
}
