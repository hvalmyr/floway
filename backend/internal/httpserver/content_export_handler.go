package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type contentExportHandler struct {
	svc   *service.ContentExportService
	admin func(http.Handler) http.Handler
}

func newContentExportHandler(svc *service.ContentExportService, admin func(http.Handler) http.Handler) *contentExportHandler {
	return &contentExportHandler{svc: svc, admin: admin}
}

// routes are entirely admin-gated (bulk read/write of every content entity
// at once, not a public listing) — unlike most other handlers in this
// package, admin applies to the whole subrouter rather than per-method.
func (h *contentExportHandler) routes(r chi.Router) {
	r.Use(h.admin)
	r.Get("/export", h.export)
	r.Post("/import", h.importContent)
}

func (h *contentExportHandler) export(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.Export(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	// Content-Disposition makes the browser download this as a file (the
	// admin UI's "Экспорт" button is a plain link/fetch, not a form) instead
	// of navigating to it.
	w.Header().Set("Content-Disposition", `attachment; filename="floway-content-export.json"`)
	writeJSON(w, http.StatusOK, data)
}

// importRequest wraps the exported document with the mode, rather than
// reading the mode from a query param, so the whole request is one JSON
// body — simpler for the admin UI's fetch call than juggling a query string
// alongside a multipart/file body.
type importRequest struct {
	Mode service.ImportMode `json:"mode"`
	Data model.SiteContent  `json:"data"`
}

func (h *contentExportHandler) importContent(w http.ResponseWriter, r *http.Request) {
	var req importRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.svc.Import(r.Context(), req.Data, req.Mode)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
