package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"floway-backend/internal/model"
)

const telegramAPIBase = "https://api.telegram.org"

type TelegramNotifier struct {
	botToken string
	chatID   string
	apiBase  string
	client   *http.Client
}

func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{botToken: botToken, chatID: chatID, apiBase: telegramAPIBase, client: http.DefaultClient}
}

// newTelegramNotifierForTest points at a test double instead of the real
// Telegram API — unexported, used only by telegram_test.go in this package.
func newTelegramNotifierForTest(botToken, chatID, apiBase string, client *http.Client) *TelegramNotifier {
	return &TelegramNotifier{botToken: botToken, chatID: chatID, apiBase: apiBase, client: client}
}

func (n *TelegramNotifier) NotifyNewLead(ctx context.Context, lead model.Lead, programName string) error {
	payload, err := json.Marshal(map[string]string{
		"chat_id": n.chatID,
		"text":    formatLeadText(lead, programName),
	})
	if err != nil {
		return fmt.Errorf("encode telegram payload: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", n.apiBase, n.botToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send telegram notification: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram sendMessage: unexpected status %d", resp.StatusCode)
	}
	return nil
}
