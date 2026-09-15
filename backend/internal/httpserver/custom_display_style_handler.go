package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type customDisplayStyleHandler struct {
	svc   *service.CustomDisplayStyleService
	admin func(http.Handler) http.Handler
}

func newCustomDisplayStyleHandler(svc *service.CustomDisplayStyleService, admin func(http.Handler) http.Handler) *customDisplayStyleHandler {
	return &customDisplayStyleHandler{svc: svc, admin: admin}
}

// list is public (no admin middleware) — the public site itself needs to
// resolve a course's customDisplayStyleId to actual colors when rendering
// its card, same as gallery-photos' list route.
func (h *customDisplayStyleHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type customDisplayStyleRequest struct {
	Name      string `json:"name"`
	BgColor   string `json:"bgColor"`
	TextColor string `json:"textColor"`
	SortOrder int    `json:"sortOrder"`
}

func (h *customDisplayStyleHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *customDisplayStyleHandler) create(w http.ResponseWriter, r *http.Request) {
	var req customDisplayStyleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), model.CustomDisplayStyle{
		Name:      req.Name,
		BgColor:   req.BgColor,
		TextColor: req.TextColor,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *customDisplayStyleHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req customDisplayStyleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Update(r.Context(), model.CustomDisplayStyle{
		ID:        id,
		Name:      req.Name,
		BgColor:   req.BgColor,
		TextColor: req.TextColor,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *customDisplayStyleHandler) delete(w http.ResponseWriter, r *http.Request) {
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
