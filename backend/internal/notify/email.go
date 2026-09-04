package notify

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"

	"floway-backend/internal/model"
)

// EmailNotifier sends a multipart/alternative notification (plain text +
// HTML) via SMTP. Auth is optional: with SMTP_USER/SMTP_PASSWORD set it uses
// smtp.PlainAuth (needed by providers like mail.ru/Yandex); left empty it
// sends with no auth, which is enough for Mailhog locally and for a relay
// that trusts by network/IP.
// Not independently unit-tested: net/smtp.SendMail is a package-level
// function, not swappable without a DI seam this thin wrapper doesn't
// need — same trade-off as internal/storage.Client, verified live
// (Mailhog locally). The HTML template rendering itself (renderLeadEmailHTML)
// is pure and tested separately.
type EmailNotifier struct {
	host, port, from, to, adminURL string
	user, password                 string
}

func NewEmailNotifier(host, port, from, to, adminURL, user, password string) *EmailNotifier {
	return &EmailNotifier{host: host, port: port, from: from, to: to, adminURL: adminURL, user: user, password: password}
}

// programName is the resolved course/masterclass title for lead.RelatedSlug
// (empty for trial lessons, or if the lookup failed) — see LeadService.
func (n *EmailNotifier) NotifyNewLead(ctx context.Context, lead model.Lead, programName string) error {
	htmlBody, err := renderLeadEmailHTML(lead, programName, n.adminURL)
	if err != nil {
		return fmt.Errorf("render lead notification email: %w", err)
	}

	var parts bytes.Buffer
	writer := multipart.NewWriter(&parts)

	textPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=utf-8"},
		"Content-Transfer-Encoding": {"8bit"},
	})
	if err != nil {
		return fmt.Errorf("build lead notification email: %w", err)
	}
	if _, err := textPart.Write([]byte(formatLeadText(lead, programName))); err != nil {
		return fmt.Errorf("build lead notification email: %w", err)
	}

	htmlPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=utf-8"},
		"Content-Transfer-Encoding": {"8bit"},
	})
	if err != nil {
		return fmt.Errorf("build lead notification email: %w", err)
	}
	if _, err := htmlPart.Write([]byte(htmlBody)); err != nil {
		return fmt.Errorf("build lead notification email: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("build lead notification email: %w", err)
	}

	// Cyrillic subjects need RFC 2047 encoding — a raw UTF-8 Subject header
	// renders blank or garbled in clients that don't assume UTF-8 (Mail.ru
	// webmail among them).
	subject := mime.QEncoding.Encode("utf-8", leadEmailSubject(lead, programName))
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=%s\r\n\r\n%s",
		n.from, n.to, subject, writer.Boundary(), parts.String(),
	)

	addr := net.JoinHostPort(n.host, n.port)
	var auth smtp.Auth
	if n.user != "" || n.password != "" {
		auth = smtp.PlainAuth("", n.user, n.password, n.host)
	}
	if err := smtp.SendMail(addr, auth, n.from, []string{n.to}, []byte(msg)); err != nil {
		return fmt.Errorf("send lead notification email: %w", err)
	}
	return nil
}
