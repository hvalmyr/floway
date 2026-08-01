package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"floway-backend/internal/service"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// writeError is for 4xx responses, where err.Error() is safe to show a
// client (validation messages, "not found", bad request). Never call this
// with a raw internal/driver error — use writeInternalError for that.
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

// writeInternalError is for 500s. The real error (which can carry SQL
// fragments, table/column names, or internal hostnames — see architecture
// review finding #5) is logged server-side against the request ID; the
// client gets a fixed, non-disclosing message plus that same request ID so
// a bug report can be correlated back to the log line.
func writeInternalError(w http.ResponseWriter, r *http.Request, err error) {
	reqID := middleware.GetReqID(r.Context())
	loggerFromContext(r.Context()).Error("internal error", "request_id", reqID, "error", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error":     "internal error",
		"requestId": reqID,
	})
}

func writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, service.ErrValidation) {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, service.ErrConflict) {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeInternalError(w, r, err)
}
