package notify

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
)

func TestDigitsOnly(t *testing.T) {
	assert.Equal(t, "+79991234567", digitsOnly("+7 (999) 123-45-67"))
	assert.Equal(t, "89991234567", digitsOnly("8 999 123 45 67"))
}

func TestPrimaryContactButton(t *testing.T) {
	cases := []struct {
		method    model.ContactMethod
		wantLabel string
		wantHref  string
	}{
		{model.ContactMethodTelegram, "Telegram", "https://t.me/+79991234567"},
		{model.ContactMethodWhatsapp, "WhatsApp", "https://wa.me/79991234567"},
		{model.ContactMethodCall, "позвонить", "tel:+79991234567"},
		// Max has no per-contact deep link — falls back to a call so the
		// manager can still act immediately.
		{model.ContactMethodMax, "позвонить", "tel:+79991234567"},
	}
	for _, tc := range cases {
		t.Run(string(tc.method), func(t *testing.T) {
			lead := testLead()
			lead.ContactMethod = tc.method
			lead.Phone = "+79991234567"

			gotLabel, gotHref := primaryContactButton(lead)

			assert.Equal(t, tc.wantLabel, gotLabel)
			assert.Equal(t, tc.wantHref, gotHref)
		})
	}
}

func TestUsesCallAsPrimary(t *testing.T) {
	assert.False(t, usesCallAsPrimary(leadWithContactMethod(model.ContactMethodTelegram)))
	assert.False(t, usesCallAsPrimary(leadWithContactMethod(model.ContactMethodWhatsapp)))
	assert.True(t, usesCallAsPrimary(leadWithContactMethod(model.ContactMethodCall)))
	assert.True(t, usesCallAsPrimary(leadWithContactMethod(model.ContactMethodMax)), "Max falls back to a call, so it counts as \"already the call button\" too")
}

func leadWithContactMethod(m model.ContactMethod) model.Lead {
	lead := testLead()
	lead.ContactMethod = m
	return lead
}

func TestProgramDisplay(t *testing.T) {
	lead := testLead()
	lead.RequestType = model.LeadRequestTypeCourse

	// No resolved title (trial lesson, or the lookup failed) — falls back
	// to the plain request-type label, no raw slug leaking through.
	assert.Equal(t, "Курс", programDisplay(lead, ""))

	// A resolved title is used as-is — no "Курс: " prefix, since the field
	// label above it already says that.
	assert.Equal(t, "Актуальная флористика", programDisplay(lead, "Актуальная флористика"))
}

func TestRenderLeadEmailHTML_IncludesLeadDataAndChosenContactChannel(t *testing.T) {
	lead := testLead()
	lead.Phone = "+79991234567"
	lead.ContactMethod = model.ContactMethodWhatsapp

	html, err := renderLeadEmailHTML(lead, "Актуальная флористика", "https://floway.example/admin/leads")

	require.NoError(t, err)
	assert.Contains(t, html, "Иван Иванов")
	assert.Contains(t, html, "Актуальная флористика")
	assert.NotContains(t, html, "Курс:", "the field label already says \"Курс / мастер-класс\" — the value must not repeat it")
	// html/template HTML-escapes "+" to "&#43;" everywhere it renders
	// {{.Phone}}/{{.PhoneRaw}} — same visual result in a browser, just the
	// raw-source form to assert on.
	assert.Contains(t, html, "&#43;79991234567")
	assert.Contains(t, html, `href="https://wa.me/79991234567"`, "primary CTA follows the client's chosen channel, not a hardcoded one")
	assert.Contains(t, html, `href="tel:&#43;79991234567"`, "WhatsApp lead still gets a separate call button")
	assert.Contains(t, html, `href="https://floway.example/admin/leads"`)
	assert.NotContains(t, html, "t.me", "Telegram CTA must not appear when the client picked WhatsApp")
	assert.Contains(t, html, "data:font/woff;base64,", "brand fonts are embedded inline, not linked externally")
}

func TestRenderLeadEmailHTML_CallPreferenceShowsOnlyOneCallButton(t *testing.T) {
	lead := testLead()
	lead.Phone = "+79991234567"
	lead.ContactMethod = model.ContactMethodCall

	html, err := renderLeadEmailHTML(lead, "", "https://floway.example/admin/leads")

	require.NoError(t, err)
	callButtons := strings.Count(html, `href="tel:&#43;79991234567"`)
	assert.Equal(t, 1, callButtons, "must not render the primary CTA and the dedicated call button as two identical \"позвонить\" buttons")
}

func TestRenderLeadEmailHTML_CapitalizesSourceAndContactMethod(t *testing.T) {
	lead := testLead()
	lead.Source = model.LeadSourceInternet
	lead.ContactMethod = model.ContactMethodCall

	html, err := renderLeadEmailHTML(lead, "", "https://floway.example/admin/leads")

	require.NoError(t, err)
	assert.Contains(t, html, ">Интернет<")
	assert.Contains(t, html, ">Звонок<")
	assert.NotContains(t, html, ">интернет<")
}

func TestRenderLeadEmailHTML_FooterYearMatchesSubmission(t *testing.T) {
	lead := testLead()
	lead.CreatedAt = time.Date(2031, time.January, 1, 0, 0, 0, 0, time.UTC)

	html, err := renderLeadEmailHTML(lead, "", "https://floway.example/admin/leads")

	require.NoError(t, err)
	assert.Contains(t, html, "Floway · Москва · 2031")
	assert.NotContains(t, html, "2013", "footer year must be the submission year, not the school's founding year")
}

func TestRenderLeadEmailHTML_EscapesLeadSuppliedName(t *testing.T) {
	lead := testLead()
	lead.Name = `<script>alert(1)</script>`

	html, err := renderLeadEmailHTML(lead, "", "https://floway.example/admin/leads")

	require.NoError(t, err)
	assert.False(t, strings.Contains(html, "<script>"), "lead-supplied name must be HTML-escaped, not injected raw")
	assert.Contains(t, html, "&lt;script&gt;")
}

func TestRenderLeadEmailHTML_MissingEmailShowsPlaceholder(t *testing.T) {
	lead := testLead()
	lead.Email = ""

	html, err := renderLeadEmailHTML(lead, "", "https://floway.example/admin/leads")

	require.NoError(t, err)
	assert.Contains(t, html, "—")
}
