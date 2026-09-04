package httpserver_test

import (
	"context"
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

type fakeLeadRepository struct{ nextID int64 }

func (f *fakeLeadRepository) ListWithClient(ctx context.Context) ([]model.LeadListItem, error) {
	return nil, nil
}

func (f *fakeLeadRepository) Create(ctx context.Context, item model.Lead) (model.Lead, error) {
	f.nextID++
	item.ID = f.nextID
	return item, nil
}

func (f *fakeLeadRepository) UpdateStatus(ctx context.Context, id int64, status model.LeadStatus) (model.Lead, error) {
	return model.Lead{}, service.ErrNotFound
}

func (f *fakeLeadRepository) DismissReview(ctx context.Context, id int64) (model.Lead, error) {
	return model.Lead{}, service.ErrNotFound
}

func (f *fakeLeadRepository) CountByStatus(ctx context.Context, statuses ...model.LeadStatus) (map[model.LeadStatus]int, error) {
	return nil, nil
}

func (f *fakeLeadRepository) Delete(ctx context.Context, id int64) error { return service.ErrNotFound }

type fakeClientRepository struct{ nextID int64 }

func (f *fakeClientRepository) FindByPhoneOrEmail(ctx context.Context, phoneNormalized, email string) (model.Client, error) {
	return model.Client{}, service.ErrNotFound
}

func (f *fakeClientRepository) Create(ctx context.Context, item model.Client) (model.Client, error) {
	f.nextID++
	item.ID = f.nextID
	return item, nil
}

func (f *fakeClientRepository) RefreshContactInfo(ctx context.Context, id int64, item model.Client) (model.Client, error) {
	item.ID = id
	return item, nil
}

// The public lead-submission endpoint costs a spammer nothing to hit —
// architecture review finding #8. A burst past the configured limit from
// the same client IP must eventually get 429, not silently create rows
// forever.
func TestLeadsBoundary_RateLimitedAfterBurst(t *testing.T) {
	services := httpserver.Services{
		Tokens:         auth.NewTokenManager("test-secret", time.Hour),
		FrontendOrigin: "http://localhost:3000",
		Lead:           service.NewLeadService(&fakeLeadRepository{}, &fakeClientRepository{}, nil, nil, nil),
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	body := `{"name":"Test","phone":"+79001234567","contactMethod":"call","source":"internet","requestType":"trial_lesson"}`

	var sawTooManyRequests bool
	for range 20 {
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/leads", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		// The test server has no real reverse proxy, so GetClientIP would
		// normally be empty (fail-open) — set the trusted XFF hop by hand so
		// this test actually exercises the limiter instead of bypassing it.
		req.Header.Set("X-Forwarded-For", "203.0.113.7")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		_ = resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			sawTooManyRequests = true
			break
		}
	}

	if !sawTooManyRequests {
		t.Fatal("20 rapid POSTs from the same IP never got rate-limited")
	}
}
