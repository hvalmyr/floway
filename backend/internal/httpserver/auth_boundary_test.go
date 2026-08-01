package httpserver_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
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

// statefulAdminUserRepository is a real in-memory stand-in — unlike
// fakeAdminUserRepositoryForAuth, IncrementTokenVersion here actually
// mutates state, so a login→logout→reuse flow through the real HTTP
// handlers exercises the genuine revocation contract end to end.
type statefulAdminUserRepository struct {
	items  []model.AdminUser
	nextID int64
}

func (f *statefulAdminUserRepository) FindByLogin(ctx context.Context, login string) (model.AdminUser, error) {
	for _, existing := range f.items {
		if existing.Login == login {
			return existing, nil
		}
	}
	return model.AdminUser{}, service.ErrNotFound
}

func (f *statefulAdminUserRepository) FindByID(ctx context.Context, id int64) (model.AdminUser, error) {
	for _, existing := range f.items {
		if existing.ID == id {
			return existing, nil
		}
	}
	return model.AdminUser{}, service.ErrNotFound
}

func (f *statefulAdminUserRepository) Create(ctx context.Context, login, passwordHash string) (model.AdminUser, error) {
	f.nextID++
	user := model.AdminUser{ID: f.nextID, Login: login, PasswordHash: passwordHash}
	f.items = append(f.items, user)
	return user, nil
}

func (f *statefulAdminUserRepository) IncrementTokenVersion(ctx context.Context, id int64) error {
	for i, existing := range f.items {
		if existing.ID == id {
			f.items[i].TokenVersion++
			return nil
		}
	}
	return service.ErrNotFound
}

// A logged-out session's cookie must stop working immediately, not merely
// once its 12h JWT naturally expires — architecture review finding #16:
// logout used to only clear the cookie client-side, so a copied/leaked
// token (or the same cookie replayed by hand, as this test does) stayed
// valid for the rest of its TTL regardless of logout.
func TestAuthBoundary_LogoutInvalidatesExistingSession(t *testing.T) {
	repo := &statefulAdminUserRepository{}
	adminSvc := service.NewAdminUserService(repo)
	_, err := adminSvc.Create(context.Background(), "admin", "supersecret")
	require.NoError(t, err)

	tokens := auth.NewTokenManager("test-secret", time.Hour)
	services := httpserver.Services{
		Tokens:         tokens,
		FrontendOrigin: "http://localhost:3000",
		AdminUser:      adminSvc,
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{Jar: jar}

	loginResp, err := client.Post(srv.URL+"/api/v1/admin/login", "application/json",
		strings.NewReader(`{"login":"admin","password":"supersecret"}`))
	require.NoError(t, err)
	defer func() { _ = loginResp.Body.Close() }()
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	// The now-authenticated cookie is captured (not just carried by the
	// jar) so it can be replayed by hand after logout clears the jar's copy.
	staleCookies := loginResp.Cookies()
	require.NotEmpty(t, staleCookies)

	meResp, err := client.Get(srv.URL + "/api/v1/admin/me")
	require.NoError(t, err)
	defer func() { _ = meResp.Body.Close() }()
	require.Equal(t, http.StatusOK, meResp.StatusCode, "the freshly issued session must work")

	logoutResp, err := client.Post(srv.URL+"/api/v1/admin/logout", "application/json", nil)
	require.NoError(t, err)
	defer func() { _ = logoutResp.Body.Close() }()
	require.Equal(t, http.StatusNoContent, logoutResp.StatusCode)

	replay, err := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/admin/me", nil)
	require.NoError(t, err)
	for _, c := range staleCookies {
		replay.AddCookie(c)
	}
	replayResp, err := http.DefaultClient.Do(replay)
	require.NoError(t, err)
	defer func() { _ = replayResp.Body.Close() }()

	if replayResp.StatusCode != http.StatusUnauthorized {
		var body map[string]any
		_ = json.NewDecoder(replayResp.Body).Decode(&body)
		t.Fatalf("replaying the pre-logout cookie: want 401, got %d (body=%v) — logout must invalidate the JWT, not just clear the client's cookie", replayResp.StatusCode, body)
	}
}
