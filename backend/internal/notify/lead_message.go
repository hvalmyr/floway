package notify

import (
	"fmt"
	"strings"

	"floway-backend/internal/model"
)

var contactMethodLabels = map[model.ContactMethod]string{
	model.ContactMethodCall:     "звонок",
	model.ContactMethodTelegram: "Telegram",
	model.ContactMethodWhatsapp: "WhatsApp",
	model.ContactMethodMax:      "Max",
}

var leadSourceLabels = map[model.LeadSource]string{
	model.LeadSourceReferral: "по рекомендации",
	model.LeadSourceAds:      "реклама",
	model.LeadSourceInternet: "интернет",
	model.LeadSourceSocial:   "соцсети",
	model.LeadSourceMaps:     "карты",
}

var leadRequestTypeLabels = map[model.LeadRequestType]string{
	model.LeadRequestTypeCourse:      "курс",
	model.LeadRequestTypeMasterclass: "мастер-класс",
	model.LeadRequestTypeTrialLesson: "пробное занятие",
}

// label falls back to the raw enum value for anything not in the map above,
// so an unrecognized value still shows up in the notification instead of
// silently disappearing.
func label[T ~string](value T, known map[T]string) string {
	if l, ok := known[value]; ok {
		return l
	}
	return string(value)
}

// formatLeadText renders a plain-text summary shared by every notification
// channel — email and Telegram both read fine as plain text, so one format
// covers both rather than maintaining near-duplicate templates. programName
// is the resolved course/masterclass title (empty for trial-lesson leads,
// or if the lookup came back empty) — the caller resolves it since this
// package has no repository dependency of its own.
func formatLeadText(lead model.Lead, programName string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Новая заявка с сайта Floway\n\n")
	fmt.Fprintf(&b, "Имя: %s\n", lead.Name)
	fmt.Fprintf(&b, "Телефон: %s\n", lead.Phone)
	if lead.Email != "" {
		fmt.Fprintf(&b, "Email: %s\n", lead.Email)
	}
	fmt.Fprintf(&b, "Способ связи: %s\n", label(lead.ContactMethod, contactMethodLabels))
	fmt.Fprintf(&b, "Источник: %s\n", label(lead.Source, leadSourceLabels))
	fmt.Fprintf(&b, "Запрос: %s\n", label(lead.RequestType, leadRequestTypeLabels))
	if programName != "" {
		fmt.Fprintf(&b, "Программа: %s\n", programName)
	} else if lead.RelatedSlug != "" {
		fmt.Fprintf(&b, "Страница: %s\n", lead.RelatedSlug)
	}
	return b.String()
}

// leadEmailSubject names what the lead is about and who submitted it, so it
// reads in an inbox list without opening the message — e.g. "Заявка с
// сайта: Пробное занятие — Иван Иванов".
func leadEmailSubject(lead model.Lead, programName string) string {
	what := programName
	if what == "" {
		what = capitalizeFirst(label(lead.RequestType, leadRequestTypeLabels))
	}
	return fmt.Sprintf("Заявка с сайта: %s — %s", what, lead.Name)
}
