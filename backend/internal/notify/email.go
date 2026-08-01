package notify

import (
	"context"
	"fmt"
	"net"
	"net/smtp"

	"floway-backend/internal/model"
)

// EmailNotifier sends a plain-text notification via SMTP with no auth —
// matches the .env.example surface (SMTP_HOST/PORT/FROM, no SMTP_USER/
// PASSWORD), which is enough for Mailhog locally and for a relay that
// trusts by network/IP in production. Not independently unit-tested:
// net/smtp.SendMail is a package-level function, not swappable without a
// DI seam this thin wrapper doesn't need — same trade-off as
// internal/storage.Client, verified live (Mailhog locally).
type EmailNotifier struct {
	host, port, from, to string
}

func NewEmailNotifier(host, port, from, to string) *EmailNotifier {
	return &EmailNotifier{host: host, port: port, from: from, to: to}
}

func (n *EmailNotifier) NotifyNewLead(ctx context.Context, lead model.Lead) error {
	addr := net.JoinHostPort(n.host, n.port)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Новая заявка с сайта Floway\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		n.from, n.to, formatLeadText(lead),
	)
	if err := smtp.SendMail(addr, nil, n.from, []string{n.to}, []byte(msg)); err != nil {
		return fmt.Errorf("send lead notification email: %w", err)
	}
	return nil
}
