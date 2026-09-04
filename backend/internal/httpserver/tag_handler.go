package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type tagHandler struct {
	svc   *service.TagService
	admin func(http.Handler) http.Handler
}

func newTagHandler(svc *service.TagService, admin func(http.Handler) http.Handler) *tagHandler {
	return &tagHandler{svc: svc, admin: admin}
}

func (h *tagHandler) routes(r chi.Router) {
	r.With(h.admin).Get("/", h.search)
	r.With(h.admin).Patch("/{id}", h.setColor)
	r.With(h.admin).Delete("/{id}", h.delete)
}

// search backs both the filter dropdown (?type=product, no q) and the
// client-detail tag autocomplete (?type=product&q=свад).
func (h *tagHandler) search(w http.ResponseWriter, r *http.Request) {
	tagType := model.TagType(r.URL.Query().Get("type"))
	query := r.URL.Query().Get("q")

	tags, err := h.svc.Search(r.Context(), tagType, query)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tags)
}

type setTagColorRequest struct {
	Color string `json:"color"`
}

// setColor changes a tag's background color (?type=product|client_type) —
// the color is a property of the tag definition, so this affects every
// client the tag is assigned to at once.
func (h *tagHandler) setColor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	tagType := model.TagType(r.URL.Query().Get("type"))

	var req setTagColorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tag, err := h.svc.SetColor(r.Context(), tagType, id, req.Color)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, tag)
}

// delete removes the tag definition itself (?type=product|client_type) —
// not a per-client unassign, which is PUT .../tags/{type} with a shorter
// name list.
func (h *tagHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	tagType := model.TagType(r.URL.Query().Get("type"))

	if err := h.svc.Delete(r.Context(), tagType, id); err != nil {
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
