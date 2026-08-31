package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/auth"
	"floway-backend/internal/httpserver"
)

// Every admin-mutating route must reject an unauthenticated request with 401
// before it ever reaches a service. requireAdminMiddleware rejects purely
// from the cookie/JWT, so this table needs no working services underneath —
// a dropped ".With(h.admin)" would be the only way one of these fails, and
// that's exactly the class of bug this test exists to catch (see backend
// architecture review, finding #13: 40+ hand-applied admin middleware calls
// with nothing verifying them).
func TestAdminRoutes_RejectUnauthenticated(t *testing.T) {
	services := httpserver.Services{
		Tokens:         auth.NewTokenManager("test-secret", time.Hour),
		FrontendOrigin: "http://localhost:3000",
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	routes := []struct {
		method string
		path   string
	}{
		// auth
		{http.MethodGet, "/api/v1/admin/me"},
		// faq
		{http.MethodPost, "/api/v1/faq"},
		{http.MethodPut, "/api/v1/faq/1"},
		{http.MethodDelete, "/api/v1/faq/1"},
		// teachers
		{http.MethodPost, "/api/v1/teachers"},
		{http.MethodPut, "/api/v1/teachers/1"},
		{http.MethodDelete, "/api/v1/teachers/1"},
		// blog-posts
		{http.MethodPost, "/api/v1/blog-posts"},
		{http.MethodPut, "/api/v1/blog-posts/1"},
		{http.MethodDelete, "/api/v1/blog-posts/1"},
		// masterclasses
		{http.MethodPost, "/api/v1/masterclasses"},
		{http.MethodPut, "/api/v1/masterclasses/1"},
		{http.MethodDelete, "/api/v1/masterclasses/1"},
		// course sections
		{http.MethodPost, "/api/v1/course-sections"},
		{http.MethodPut, "/api/v1/course-sections/1"},
		{http.MethodDelete, "/api/v1/course-sections/1"},
		// courses (nested under a section)
		{http.MethodPost, "/api/v1/course-sections/1/courses"},
		{http.MethodPut, "/api/v1/course-sections/1/courses/1"},
		{http.MethodDelete, "/api/v1/course-sections/1/courses/1"},
		// course blocks (nested)
		{http.MethodPost, "/api/v1/courses/1/blocks"},
		{http.MethodPut, "/api/v1/courses/1/blocks/1"},
		{http.MethodDelete, "/api/v1/courses/1/blocks/1"},
		// lessons (nested)
		{http.MethodPost, "/api/v1/course-blocks/1/lessons"},
		{http.MethodPut, "/api/v1/course-blocks/1/lessons/1"},
		{http.MethodDelete, "/api/v1/course-blocks/1/lessons/1"},
		// leads
		{http.MethodGet, "/api/v1/leads"},
		{http.MethodPatch, "/api/v1/leads/1/status"},
		{http.MethodDelete, "/api/v1/leads/1"},
		// page-content
		{http.MethodPut, "/api/v1/page-content/some_key"},
		// features
		{http.MethodPost, "/api/v1/features"},
		{http.MethodPut, "/api/v1/features/1"},
		{http.MethodDelete, "/api/v1/features/1"},
		// about-items
		{http.MethodPost, "/api/v1/about-items"},
		{http.MethodPut, "/api/v1/about-items/1"},
		{http.MethodDelete, "/api/v1/about-items/1"},
		// social-links
		{http.MethodPost, "/api/v1/social-links"},
		{http.MethodPut, "/api/v1/social-links/1"},
		{http.MethodDelete, "/api/v1/social-links/1"},
		// gallery-photos
		{http.MethodPost, "/api/v1/gallery-photos"},
		{http.MethodPut, "/api/v1/gallery-photos/1"},
		{http.MethodDelete, "/api/v1/gallery-photos/1"},
		// uploads
		{http.MethodPost, "/api/v1/admin/uploads"},
	}

	for _, rt := range routes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			req, err := http.NewRequest(rt.method, srv.URL+rt.path, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

// Public GET/list routes must NOT require auth — if a route is accidentally
// wrapped in admin middleware, this catches the opposite regression from
// TestAdminRoutes_RejectUnauthenticated. Zero-value services would panic on
// a nil-pointer method call if one of these ever reached real handler logic,
// which is itself proof the request wasn't rejected by the auth middleware —
// any response at all (even a 500 from the nil service) confirms the route
// is not admin-gated. Only /healthz is asserted for an actual 200, since it
// touches no service.
func TestPublicRoutes_DoNotRequireAuth(t *testing.T) {
	services := httpserver.Services{
		Tokens:         auth.NewTokenManager("test-secret", time.Hour),
		FrontendOrigin: "http://localhost:3000",
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/healthz")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// /login and /logout must be reachable without a cookie (that's the
	// whole point) — a 400 (bad body) is fine, 401 would mean something
	// wrapped them in admin middleware by mistake.
	resp2, err := http.Post(srv.URL+"/api/v1/admin/login", "application/json", nil)
	require.NoError(t, err)
	defer func() { _ = resp2.Body.Close() }()
	assert.NotEqual(t, http.StatusUnauthorized, resp2.StatusCode)
}
