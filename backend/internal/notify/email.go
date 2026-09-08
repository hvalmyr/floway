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
	"strings"

	"floway-backend/internal/model"
)

// NotificationEmailLister supplies the current recipient list at send
// time (not once at startup) so the admin-editable list in the
// notification_emails table (see NotificationEmailService) takes effect on
// the very next lead without a redeploy.
type NotificationEmailLister interface {
	List(ctx context.Context) ([]model.NotificationEmail, error)
}

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
	host, port, from, adminURL string
	user, password             string
	recipients                 NotificationEmailLister
}

func NewEmailNotifier(host, port, from, adminURL, user, password string, recipients NotificationEmailLister) *EmailNotifier {
	return &EmailNotifier{host: host, port: port, from: from, adminURL: adminURL, user: user, password: password, recipients: recipients}
}

// programName is the resolved course/masterclass title for lead.RelatedSlug
// (empty for trial lessons, or if the lookup failed) — see LeadService.
func (n *EmailNotifier) NotifyNewLead(ctx context.Context, lead model.Lead, programName string) error {
	recipients, err := n.recipients.List(ctx)
	if err != nil {
		return fmt.Errorf("list notification emails: %w", err)
	}
	if len(recipients) == 0 {
		// No recipients configured in the admin panel — nothing to send,
		// not an error (same "optional channel" treatment as an
		// unconfigured SMTP host).
		return nil
	}
	to := make([]string, len(recipients))
	for i, r := range recipients {
		to[i] = r.Email
	}

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
		n.from, strings.Join(to, ", "), subject, writer.Boundary(), parts.String(),
	)

	addr := net.JoinHostPort(n.host, n.port)
	var auth smtp.Auth
	if n.user != "" || n.password != "" {
		auth = smtp.PlainAuth("", n.user, n.password, n.host)
	}
	if err := smtp.SendMail(addr, auth, n.from, to, []byte(msg)); err != nil {
		return fmt.Errorf("send lead notification email: %w", err)
	}
	return nil
}
