package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/auth"
	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type blogPostHandler struct {
	svc     *service.BlogPostService
	tokens  *auth.TokenManager
	checker tokenVersionChecker
	admin   func(http.Handler) http.Handler
}

func newBlogPostHandler(svc *service.BlogPostService, tokens *auth.TokenManager, checker tokenVersionChecker, admin func(http.Handler) http.Handler) *blogPostHandler {
	return &blogPostHandler{svc: svc, tokens: tokens, checker: checker, admin: admin}
}

func (h *blogPostHandler) routes(r chi.Router) {
	r.Get("/", h.list)
	// Public detail lookup by slug (published posts only), admin mutations
	// by numeric id — same URL shape ("/{something}"), dispatched by
	// method, so the path itself is unchanged for clients; only the chi
	// param name here reflects what each method actually expects, instead
	// of both being called "id".
	r.Get("/{slug}", h.getPublishedBySlug)
	r.With(h.admin).Post("/", h.create)
	r.With(h.admin).Put("/{id}", h.update)
	r.With(h.admin).Delete("/{id}", h.delete)
}

type blogPostRequest struct {
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	MetaTitle       string   `json:"metaTitle"`
	MetaDescription string   `json:"metaDescription"`
	CoverImage      string   `json:"coverImage"`
	DisplayStyle    string   `json:"displayStyle"`
	Category        string   `json:"category"`
	Tags            []string `json:"tags"`
	Author          string   `json:"author"`
	// A plain string, not *time.Time — the admin form (a free-text field,
	// not a date picker) lets an admin type anything, and *time.Time's
	// strict RFC3339-only UnmarshalJSON used to fail the whole request
	// (every other field along with it) on the slightest typo. Parsed
	// leniently in toModel() instead: an unparseable value just saves as
	// "no publish date" rather than blocking the save entirely.
	PublishedAt *string `json:"publishedAt,omitempty"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
}

func parsePublishedAt(raw *string) *time.Time {
	if raw == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil
	}
	return &t
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
	if isAdminRequest(r, h.tokens, h.checker) {
		items, err = h.svc.List(r.Context())
	} else {
		items, err = h.svc.ListPublished(r.Context())
	}
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *blogPostHandler) getPublishedBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	item, err := h.svc.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *blogPostHandler) toModel(req blogPostRequest) model.BlogPost {
	return model.BlogPost{
		Slug:            req.Slug,
		Title:           req.Title,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		CoverImage:      req.CoverImage,
		DisplayStyle:    model.CourseBlockDisplayStyle(req.DisplayStyle),
		Category:        req.Category,
		Tags:            req.Tags,
		Author:          req.Author,
		PublishedAt:     parsePublishedAt(req.PublishedAt),
		Content:         req.Content,
		Status:          model.BlogPostStatus(req.Status),
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
		writeServiceError(w, r, err)
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
		writeServiceError(w, r, err)
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
		writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
