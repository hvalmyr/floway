package httpserver_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"floway-backend/internal/auth"
	"floway-backend/internal/httpserver"
	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeBlogPostRepository struct {
	items []model.BlogPost
}

func (f *fakeBlogPostRepository) List(ctx context.Context) ([]model.BlogPost, error) {
	return f.items, nil
}

func (f *fakeBlogPostRepository) ListPublished(ctx context.Context) ([]model.BlogPost, error) {
	var out []model.BlogPost
	for _, item := range f.items {
		if item.Status == model.BlogPostStatusPublished {
			out = append(out, item)
		}
	}
	return out, nil
}

func (f *fakeBlogPostRepository) FindPublishedBySlug(ctx context.Context, slug string) (model.BlogPost, error) {
	for _, item := range f.items {
		if item.Slug == slug && item.Status == model.BlogPostStatusPublished {
			return item, nil
		}
	}
	return model.BlogPost{}, service.ErrNotFound
}

func (f *fakeBlogPostRepository) Create(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	return item, nil
}

func (f *fakeBlogPostRepository) Update(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	return item, nil
}

func (f *fakeBlogPostRepository) Delete(ctx context.Context, id int64) error { return nil }

// A visitor with no admin session must never be able to read draft content
// through the public blog listing, regardless of query string — this was
// the default (architecture review finding #4): GET /api/v1/blog-posts
// with no query parameter returned every post, drafts included, to anyone.
func TestBlogPostBoundary_UnauthenticatedListNeverIncludesDrafts(t *testing.T) {
	tokens := auth.NewTokenManager("test-secret", time.Hour)
	repo := &fakeBlogPostRepository{items: []model.BlogPost{
		{ID: 1, Slug: "published-post", Title: "Published", Status: model.BlogPostStatusPublished},
		{ID: 2, Slug: "secret-draft", Title: "Unreleased announcement", Status: model.BlogPostStatusDraft},
	}}
	services := httpserver.Services{
		Tokens:         tokens,
		FrontendOrigin: "http://localhost:3000",
		BlogPost:       service.NewBlogPostService(repo),
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/blog-posts")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var posts []model.BlogPost
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&posts))

	for _, post := range posts {
		if post.Status != model.BlogPostStatusPublished {
			t.Fatalf("unauthenticated GET /blog-posts leaked a non-published post: %q (status=%q)", post.Slug, post.Status)
		}
	}
}

// The admin panel's own use of this same route (no ?status= param, relying
// on an authenticated session) must keep working — this is the "don't break
// the frontend contract" side of the fix.
func TestBlogPostBoundary_AuthenticatedListIncludesDrafts(t *testing.T) {
	tokens := auth.NewTokenManager("test-secret", time.Hour)
	repo := &fakeBlogPostRepository{items: []model.BlogPost{
		{ID: 1, Slug: "published-post", Title: "Published", Status: model.BlogPostStatusPublished},
		{ID: 2, Slug: "secret-draft", Title: "Unreleased announcement", Status: model.BlogPostStatusDraft},
	}}
	services := httpserver.Services{
		Tokens:         tokens,
		FrontendOrigin: "http://localhost:3000",
		BlogPost:       service.NewBlogPostService(repo),
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	req := newAuthenticatedRequest(t, tokens, http.MethodGet, srv.URL+"/api/v1/blog-posts")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var posts []model.BlogPost
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&posts))
	require.Len(t, posts, 2, "an authenticated admin request must still see drafts")
}
