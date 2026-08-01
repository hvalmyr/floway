package httpserver_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/auth"
	"floway-backend/internal/httpserver"
)

// A client-declared multipart Content-Type is just a string the caller
// chose — nothing stops it from naming an HTML/script payload "photo.png"
// with Content-Type: image/png. The upload endpoint must reject on the
// SNIFFED bytes, not the claimed header — architecture review's
// Content-Type-sniffing finding. This never reaches storage
// (services.Storage stays nil), so it needs no real Garage.
func TestUploadBoundary_RejectsFileWhoseBytesDontMatchClaimedContentType(t *testing.T) {
	tokens := auth.NewTokenManager("test-secret", time.Hour)
	services := httpserver.Services{
		Tokens:         tokens,
		FrontendOrigin: "http://localhost:3000",
		AdminUser:      newTestAdminUserService(),
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreatePart(map[string][]string{
		"Content-Disposition": {`form-data; name="file"; filename="photo.png"`},
		"Content-Type":        {"image/png"},
	})
	require.NoError(t, err)
	_, err = part.Write([]byte("<html><body><script>alert(1)</script></body></html>"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/admin/uploads", &body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	token, _, err := tokens.Issue(1, "admin", 0)
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "floway_admin_session", Value: token})

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "an HTML payload claiming to be image/png must be rejected on its real, sniffed content, not the client's declared header")
}
