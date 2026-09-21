/** Shared label maps for a Lead's structured fields — used by the request
 * list cards, the synthesized excerpt (see leadFilters.ts), and the client
 * detail page's request history so all three stay in sync. */
export const contactMethodLabels: Record<string, string> = {
  call: "Звонок",
  telegram: "Telegram",
  whatsapp: "WhatsApp",
  max: "MAX",
};

export const sourceLabels: Record<string, string> = {
  referral: "Рекомендация",
  ads: "Реклама",
  internet: "Интернет",
  social: "Соцсети",
  maps: "Карты",
};

// course/masterclass/trial_lesson/gift_certificate are the school's values —
// kept here (unused by this frontend's own LeadRequestType, which only ever
// sends "decor") purely so leadFilters.test.ts's copied coverage of
// formatLeadExcerpt keeps exercising the shared lib function unchanged.
export const requestTypeLabels: Record<string, string> = {
  course: "Курс",
  masterclass: "Мастер-класс",
  trial_lesson: "Пробный урок",
  gift_certificate: "Подарочный сертификат",
  decor: "Новогоднее оформление",
};

export const leadFormatLabels: Record<string, string> = {
  season: "На сезон",
  event: "На праздник",
};
