package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"

	"floway-backend/internal/model"
)

const telegramAPIBase = "https://api.telegram.org"

type TelegramNotifier struct {
	botToken string
	chatID   string
	apiBase  string
	client   *http.Client
}

// NewTelegramNotifier talks to the real Telegram API, optionally through a
// forward proxy — needed because RU ISPs/RKN can throttle or block direct
// access to api.telegram.org even from a server that never touches the
// Telegram app. proxyURL is optional: empty means dial directly. Supported
// schemes: "http"/"https" (a forward proxy, e.g. tinyproxy/squid on a
// non-RU VPS) and "socks5" (e.g. an `ssh -D` dynamic tunnel to a non-RU
// host — no extra software needed on that end beyond sshd).
func NewTelegramNotifier(botToken, chatID, proxyURL string) (*TelegramNotifier, error) {
	client, err := telegramHTTPClient(proxyURL)
	if err != nil {
		return nil, err
	}
	return &TelegramNotifier{botToken: botToken, chatID: chatID, apiBase: telegramAPIBase, client: client}, nil
}

func telegramHTTPClient(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		return http.DefaultClient, nil
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("parse TELEGRAM_PROXY_URL: %w", err)
	}
	switch u.Scheme {
	case "http", "https":
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}, nil
	case "socks5":
		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("build socks5 dialer for TELEGRAM_PROXY_URL: %w", err)
		}
		return &http.Client{Transport: &http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}}, nil
	default:
		return nil, fmt.Errorf("unsupported TELEGRAM_PROXY_URL scheme %q: use http, https or socks5", u.Scheme)
	}
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
