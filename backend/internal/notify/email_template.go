package notify

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"strings"
	"time"

	"floway-backend/internal/model"
)

//go:embed templates/lead_email.html.tmpl
var leadEmailTemplateSource string

var leadEmailTemplate = template.Must(template.New("lead_email").Parse(leadEmailTemplateSource))

// Brand fonts, embedded as base64 data URIs directly in the email's
// @font-face rules (see templates/fonts/) rather than linked from the
// site's own origin — email clients that support @font-face at all
// (Outlook/Gmail don't, and fall back to Arial regardless) mostly block or
// don't reliably fetch external font requests, so a self-contained email is
// the only way to actually get the brand fonts to render. Same files as
// frontend/app/assets/styles/fonts.css; duplicated here because go:embed
// can't reach outside this module.
//
//go:embed templates/fonts/SoyuzGrotesk-Bold.woff
var soyuzGroteskBoldWoff []byte

//go:embed templates/fonts/NonBureau-Light.woff2
var nonBureauLightWoff2 []byte

//go:embed templates/fonts/NonBureau-Medium.woff2
var nonBureauMediumWoff2 []byte

//go:embed templates/fonts/NonBureau-Bold.woff2
var nonBureauBoldWoff2 []byte

var (
	soyuzGroteskBoldWoffBase64 = base64.StdEncoding.EncodeToString(soyuzGroteskBoldWoff)
	nonBureauLightWoff2Base64  = base64.StdEncoding.EncodeToString(nonBureauLightWoff2)
	nonBureauMediumWoff2Base64 = base64.StdEncoding.EncodeToString(nonBureauMediumWoff2)
	nonBureauBoldWoff2Base64   = base64.StdEncoding.EncodeToString(nonBureauBoldWoff2)
)

// moscowLocation falls back to UTC if the runtime has no tzdata (shouldn't
// happen — the backend image installs it — but a formatted local time is
// worth a graceful fallback rather than a panic).
var moscowLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.UTC
	}
	return loc
}()

// leadEmailData is the html/template view model for lead_email.html.tmpl.
// Every field is pre-escaped-safe plain text — html/template context-aware
// escapes it on render, so lead-supplied values (name, phone, email) can't
// break out of the markup even though they come straight from the public
// leads form.
type leadEmailData struct {
	Name                string
	Phone               string // display form, e.g. as the visitor typed it
	PhoneRaw            string // digits (and leading +) only, safe inside tel:
	EmailDisplay        string
	Program             string
	SubmittedAt         string
	SubmittedYear       string
	ContactMethodLabel  string
	SourceLabel         string
	PrimaryContactLabel string
	// template.URL, not string: this whole attribute value is one dynamic
	// template action (`href="{{.PrimaryContactHref}}"`), and html/template
	// only trusts a handful of schemes (http/https/mailto) in that shape —
	// anything else, including "tel:", gets replaced with the inert
	// "#ZgotmplZ" placeholder. template.URL opts out of that scheme check
	// (still escaped, just not scheme-filtered); safe here because every
	// value comes from primaryContactButton, never straight from lead input.
	PrimaryContactHref template.URL
	// ShowCallButton is false when the primary CTA already *is* the call
	// button (client picked "call", or "Max" which has no per-contact link
	// and falls back to a call too) — otherwise the same "позвонить"
	// button would render twice.
	ShowCallButton bool
	AdminURL       string

	SoyuzGroteskBoldWoffBase64 string
	NonBureauLightWoff2Base64  string
	NonBureauMediumWoff2Base64 string
	NonBureauBoldWoff2Base64   string
}

// digitsOnly strips everything but digits and a leading "+" — the form
// "tel:" links expect.
func digitsOnly(phone string) string {
	var b strings.Builder
	for i, r := range phone {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r == '+' && i == 0 {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// bareDigits strips the leading "+" too — the form wa.me and t.me
// phone-based deep links expect (t.me/+<digits> supplies its own "+").
func bareDigits(phone string) string {
	return strings.TrimPrefix(digitsOnly(phone), "+")
}

// primaryContactButton picks the CTA that matches what the client actually
// asked for, instead of hardcoding Telegram — a lead who chose WhatsApp
// should land the manager in WhatsApp, not in a channel the client never
// opted into. Max has no phone-addressable deep link (it's a static,
// company-wide URL, not a per-contact one), so it falls back to a call —
// still lets the manager act immediately.
func primaryContactButton(lead model.Lead) (label, href string) {
	switch lead.ContactMethod {
	case model.ContactMethodTelegram:
		return "Telegram", "https://t.me/+" + bareDigits(lead.Phone)
	case model.ContactMethodWhatsapp:
		return "WhatsApp", "https://wa.me/" + bareDigits(lead.Phone)
	default:
		return "позвонить", "tel:" + digitsOnly(lead.Phone)
	}
}

// usesCallAsPrimary reports whether primaryContactButton above already
// resolved to a phone call — true for both an explicit "call" preference
// and the Max fallback.
func usesCallAsPrimary(lead model.Lead) bool {
	return lead.ContactMethod != model.ContactMethodTelegram && lead.ContactMethod != model.ContactMethodWhatsapp
}

// programDisplay is what shows under the "Курс / мастер-класс" label —
// the resolved course/masterclass title (looked up by the caller, since
// this package has no DB access of its own) with no repeated "Курс:" /
// "Мастер-класс:" prefix, since the label above it already says that. Falls
// back to just the request-type label when there's no specific program
// (trial lessons) or the lookup came back empty (deleted/renamed course).
func programDisplay(lead model.Lead, programName string) string {
	if programName != "" {
		return programName
	}
	return capitalizeFirst(label(lead.RequestType, leadRequestTypeLabels))
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	return strings.ToUpper(string(r[0])) + string(r[1:])
}

// renderLeadEmailHTML renders the notification email. programName is the
// resolved course/masterclass title for lead.RelatedSlug (empty if there
// isn't one, or the lookup failed) — resolved by the caller since this
// package intentionally has no repository dependency.
func renderLeadEmailHTML(lead model.Lead, programName, adminURL string) (string, error) {
	primaryLabel, primaryHref := primaryContactButton(lead)

	emailDisplay := lead.Email
	if emailDisplay == "" {
		emailDisplay = "—"
	}

	submittedAt := lead.CreatedAt.In(moscowLocation)

	data := leadEmailData{
		Name:                lead.Name,
		Phone:               lead.Phone,
		PhoneRaw:            digitsOnly(lead.Phone),
		EmailDisplay:        emailDisplay,
		Program:             programDisplay(lead, programName),
		SubmittedAt:         submittedAt.Format("02.01, 15:04"),
		SubmittedYear:       submittedAt.Format("2006"),
		ContactMethodLabel:  capitalizeFirst(label(lead.ContactMethod, contactMethodLabels)),
		SourceLabel:         capitalizeFirst(label(lead.Source, leadSourceLabels)),
		PrimaryContactLabel: primaryLabel,
		PrimaryContactHref:  template.URL(primaryHref),
		ShowCallButton:      !usesCallAsPrimary(lead),
		AdminURL:            adminURL,

		SoyuzGroteskBoldWoffBase64: soyuzGroteskBoldWoffBase64,
		NonBureauLightWoff2Base64:  nonBureauLightWoff2Base64,
		NonBureauMediumWoff2Base64: nonBureauMediumWoff2Base64,
		NonBureauBoldWoff2Base64:   nonBureauBoldWoff2Base64,
	}

	var buf bytes.Buffer
	if err := leadEmailTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render lead email template: %w", err)
	}
	return buf.String(), nil
}
