/**
 * Single source of truth for contact details shown in both AppFooter and
 * /contacts. Placeholder values below are marked TODO where the brief didn't
 * give real data (no Figma access to read exact copy off the design) —
 * replace with real values, don't ship these to production as-is.
 */
export const contactInfo = {
  // TODO: подставить реальный номер школы.
  phone: "+7 (999) 000-00-00",
  phoneHref: "tel:+79990000000",
  // TODO: подставить реальный email школы.
  email: "hello@floway.ru",
  telegramUrl: "https://t.me/floway",
  whatsappUrl: "https://wa.me/79990000000",
  // TODO: уточнить контакт/ссылку для Max у заказчика — на момент вёрстки
  // публичного способа сослаться на диалог в Max не было под рукой.
  maxUrl: "",
  address: "г. Москва, Новинский бульвар, 18Б",
  // TODO: в макете названия станций метро есть, но раскрашены (см. заметку
  // дизайн-системы §2.1 — это разовая маркировка, не токен), а без доступа к
  // Figma текст самих названий станций не считать нельзя. Подставить из макета.
  metroNote: "Станции метро уточняются",
};

export interface SocialLink {
  label: string;
  href: string;
  /** Set for platforms that legally require a disclaimer in Russia (e.g. Meta-owned apps). */
  disclaimer?: string;
}

// TODO: сверить точный набор соцсетей с макетом «Контакты», когда будет
// доступ к Figma — в брифе явно названы Telegram и VK, третья иконка описана
// как "ещё одна — сверься с макетом"; Instagram взят как правдоподобный
// вариант для школы такого профиля, но это предположение, не факт из макета.
export const socialLinks: SocialLink[] = [
  { label: "Telegram", href: contactInfo.telegramUrl },
  { label: "VK", href: "https://vk.com/floway" },
  {
    label: "Instagram",
    href: "https://instagram.com/floway",
    disclaimer: "Instagram принадлежит компании Meta, признанной экстремистской организацией и запрещённой на территории РФ.",
  },
];
