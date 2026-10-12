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

// Only set on leads whose source is "social" — see sourceLabels above.
// "instagram" is labeled as a disclaimer-safe generic name rather than
// "Instagram" itself, since Meta is recognized as extremist and banned in
// Russia (same convention as socialLinks in constants/contact-info.ts).
export const socialNetworkLabels: Record<string, string> = {
  telegram: "Telegram",
  vk: "ВКонтакте",
  pinterest: "Pinterest",
  instagram: "Запрещённая соц. сеть",
  youtube: "YouTube",
};

export const requestTypeLabels: Record<string, string> = {
  course: "Курс",
  masterclass: "Мастер-класс",
  trial_lesson: "Пробный урок",
  gift_certificate: "Подарочный сертификат",
};
