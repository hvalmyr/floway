package httpserver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/httpserver"
	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

// A request must produce exactly one structured JSON log line carrying
// enough to correlate it with a client-visible request ID and diagnose it
// without re-running the request — architecture review finding #17: chi's
// default middleware.Logger printed colorized, unstructured text meant for
// a human watching a dev terminal, not a JSON line a log aggregator (or
// grep) can filter on.
func TestRequestLogging_EmitsStructuredJSONLine(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	services := httpserver.Services{
		Logger:         logger,
		FrontendOrigin: "http://localhost:3000",
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/healthz")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	line := strings.TrimSpace(buf.String())
	require.NotEmpty(t, line, "the request must have produced a log line")

	var entry map[string]any
	require.NoError(t, json.Unmarshal([]byte(line), &entry), "log output must be valid JSON, not chi's colorized text format")

	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/healthz", entry["path"])
	assert.Equal(t, float64(http.StatusOK), entry["status"])
	assert.NotEmpty(t, entry["request_id"])
}

type loggingTestFAQRepository struct{ err error }

func (f *loggingTestFAQRepository) List(ctx context.Context) ([]model.FAQItem, error) {
	return nil, f.err
}
func (f *loggingTestFAQRepository) Create(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	return item, nil
}
func (f *loggingTestFAQRepository) Update(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	return item, nil
}
func (f *loggingTestFAQRepository) Delete(ctx context.Context, id int64) error { return nil }

// The requestId a client gets back on a 500 must match the request_id in
// the server-side log line carrying the real, otherwise-hidden error — the
// whole point of writeInternalError logging server-side and returning only
// a correlation ID (architecture review finding #5, and #17's structured
// version of it).
func TestRequestLogging_InternalErrorCorrelatesWithClientRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	rawErr := errors.New("dial tcp 10.0.4.7:5432: connect: connection refused")
	services := httpserver.Services{
		Logger:         logger,
		FrontendOrigin: "http://localhost:3000",
		FAQ:            service.NewFAQService(&loggingTestFAQRepository{err: rawErr}),
	}
	srv := httptest.NewServer(httpserver.NewRouter(services))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/faq")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var body map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	clientRequestID := body["requestId"]
	require.NotEmpty(t, clientRequestID)

	var errorLine string
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var entry map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &entry))
		if entry["level"] == "ERROR" {
			errorLine = line
			var errEntry map[string]any
			require.NoError(t, json.Unmarshal([]byte(line), &errEntry))
			assert.Equal(t, clientRequestID, errEntry["request_id"], "server log's request_id must match the requestId returned to the client")
			assert.Contains(t, errEntry["error"], "10.0.4.7", "the real error must be in the server log even though it's hidden from the client")
		}
	}
	require.NotEmpty(t, errorLine, "the internal error must have produced an ERROR-level log line")
}
