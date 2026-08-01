package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/auth"
	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type blogPostHandler struct {
	svc    *service.BlogPostService
	tokens *auth.TokenManager
	admin  func(http.Handler) http.Handler
}

func newBlogPostHandler(svc *service.BlogPostService, tokens *auth.TokenManager, admin func(http.Handler) http.Handler) *blogPostHandler {
	return &blogPostHandler{svc: svc, tokens: tokens, admin: admin}
}

func (h *blogPostHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	// Public detail lookup by slug (published posts only), shares the
	// "/{id}" pattern with the admin-only Put/Delete below — same route
	// node, dispatched by method.
	r.Get("/{id}", h.getPublishedBySlug)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type blogPostRequest struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	CoverImage  string     `json:"coverImage"`
	Category    string     `json:"category"`
	Tags        []string   `json:"tags"`
	Author      string     `json:"author"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
}

// list is a public route (no requireAdminMiddleware — the admin panel's own
// requests hit this same node), so the response depends on who's asking:
// an authenticated admin session sees every post, drafts included; anyone
// else only ever sees published posts, regardless of query string. Drafts
// must never be reachable by an unauthenticated request (architecture
// review finding #4) — this used to default to "everything" and rely on
// the frontend voluntarily passing ?status=published, which anyone probing
// the API directly could simply not do.
func (h *blogPostHandler) list(w http.ResponseWriter, r *http.Request) {
	var items []model.BlogPost
	var err error
	if isAdminRequest(r, h.tokens) {
		items, err = h.svc.List(r.Context())
	} else {
		items, err = h.svc.ListPublished(r.Context())
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *blogPostHandler) getPublishedBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "id")

	item, err := h.svc.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *blogPostHandler) toModel(req blogPostRequest) model.BlogPost {
	return model.BlogPost{
		Slug:        req.Slug,
		Title:       req.Title,
		CoverImage:  req.CoverImage,
		Category:    req.Category,
		Tags:        req.Tags,
		Author:      req.Author,
		PublishedAt: req.PublishedAt,
		Content:     req.Content,
		Status:      model.BlogPostStatus(req.Status),
	}
}

func (h *blogPostHandler) create(w http.ResponseWriter, r *http.Request) {
	var req blogPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.svc.Create(r.Context(), h.toModel(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *blogPostHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req blogPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item := h.toModel(req)
	item.ID = id
	item, err = h.svc.Update(r.Context(), item)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *blogPostHandler) delete(w http.ResponseWriter, r *http.Request) {
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
