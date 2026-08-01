package httpserver_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"floway-backend/internal/auth"
	"floway-backend/internal/httpserver"
	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// fakeFAQRepository mirrors the real repository's CONTRACT (see
// internal/repository/errors.go): a missing row on Update/Delete comes back
// as service.ErrNotFound, already translated — the fake stands in for the
// repository interface, so it must speak the same boundary language the
// real repository does, not raw driver errors.
type fakeFAQRepository struct {
	items map[int64]model.FAQItem
}

func newFakeFAQRepository(seed ...model.FAQItem) *fakeFAQRepository {
	items := make(map[int64]model.FAQItem)
	for _, item := range seed {
		items[item.ID] = item
	}
	return &fakeFAQRepository{items: items}
}

func (f *fakeFAQRepository) List(ctx context.Context) ([]model.FAQItem, error) { return nil, nil }

func (f *fakeFAQRepository) Create(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	return item, nil
}

func (f *fakeFAQRepository) Update(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	if _, ok := f.items[item.ID]; !ok {
		return model.FAQItem{}, service.ErrNotFound
	}
	f.items[item.ID] = item
	return item, nil
}

func (f *fakeFAQRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := f.items[id]; !ok {
		return service.ErrNotFound
	}
	delete(f.items, id)
	return nil
}

func newAuthenticatedRequest(t *testing.T, tokens *auth.TokenManager, method, url string) *http.Request {
	t.Helper()
	token, _, err := tokens.Issue(1, "admin")
	require.NoError(t, err)

	req, err := http.NewRequest(method, url, nil)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "floway_admin_session", Value: token})
	return req
}

func TestFAQBoundary_UpdateMissingID_Returns404NotInternalError(t *testing.T) {
	tokens := auth.NewTokenManager("test-secret", time.Hour)
	repo := newFakeFAQRepository(model.FAQItem{ID: 1, Question: "q", Answer: "a"})
	services := httpserver.Services{
		Tokens:         tokens,
		FrontendOrigin: "http://localhost:3000",
		FAQ:            service.NewFAQService(repo),
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	req := newAuthenticatedRequest(t, tokens, http.MethodPut, srv.URL+"/api/v1/faq/999")
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(`{"question":"q","answer":"a","sortOrder":0}`))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("PUT on a nonexistent FAQ item: want 404, got %d — a repository-layer \"row not found\" must never surface as a 500 with raw driver error text", resp.StatusCode)
	}
}

func TestFAQBoundary_DeleteMissingID_Returns404NotSilent204(t *testing.T) {
	tokens := auth.NewTokenManager("test-secret", time.Hour)
	repo := newFakeFAQRepository() // empty: id 999 never existed
	services := httpserver.Services{
		Tokens:         tokens,
		FrontendOrigin: "http://localhost:3000",
		FAQ:            service.NewFAQService(repo),
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	req := newAuthenticatedRequest(t, tokens, http.MethodDelete, srv.URL+"/api/v1/faq/999")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("DELETE of a nonexistent FAQ item: want 404, got %d — deleting nothing must not report success", resp.StatusCode)
	}
}
