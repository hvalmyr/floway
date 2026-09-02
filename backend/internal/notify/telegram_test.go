package notify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
)

func testLead() model.Lead {
	return model.Lead{
		ID:            42,
		Name:          "Иван Иванов",
		Phone:         "+79991234567",
		ContactMethod: model.ContactMethodTelegram,
		Source:        model.LeadSourceAds,
		RequestType:   model.LeadRequestTypeTrialLesson,
	}
}

func TestTelegramNotifier_NotifyNewLead_SendsExpectedRequest(t *testing.T) {
	var gotPath string
	var gotBody map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := newTelegramNotifierForTest("test-token", "12345", srv.URL, srv.Client())

	err := n.NotifyNewLead(context.Background(), testLead(), "")

	require.NoError(t, err)
	assert.Equal(t, "/bottest-token/sendMessage", gotPath)
	assert.Equal(t, "12345", gotBody["chat_id"])
	assert.Contains(t, gotBody["text"], "Иван Иванов")
	assert.Contains(t, gotBody["text"], "+79991234567")
	assert.Contains(t, gotBody["text"], "пробное занятие", "human-readable label, not the raw enum value")
}

func TestTelegramNotifier_NotifyNewLead_ReturnsErrorOnNonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	n := newTelegramNotifierForTest("test-token", "12345", srv.URL, srv.Client())

	err := n.NotifyNewLead(context.Background(), testLead(), "")

	require.Error(t, err)
}
