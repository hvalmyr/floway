package httpserver

import (
	"net/http"

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
