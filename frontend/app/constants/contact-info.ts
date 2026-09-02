/**
 * Single source of truth for contact details shown in AppFooter and
 * /contacts. Phone/email/address/directions/legal entity are the real
 * values from the provided mockups (docs/../desktop-contacts.png). Telegram/
 * Whatsapp/Max/VK/Instagram handles weren't shown in the mockups — those
 * stay TODO placeholders.
 */
export const contactInfo = {
  phone: "+7 985 226 19 48",
  phoneHref: "tel:+79852261948",
  email: "floway-mos@mail.ru",
  // TODO: уточнить реальные ссылки на Telegram/Whatsapp/Max у заказчика —
  // в макете это просто подписанные текстовые ссылки без видимых URL.
  telegramUrl: "https://t.me/floway",
  whatsappUrl: "https://wa.me/79852261948",
  maxUrl: "",
  address: "г. Москва, Новинский бульвар, 18Б",
  // Названия станций в макете выделены цветом — по заметке дизайн-системы
  // (§2.1) это разовая маркировка, не токен, поэтому здесь обычный текст.
  metroStations: ["Смоленская", "Баррикадная", "Краснопресненская"],
  directions: [
    "Чтобы попасть к нам, зайдите в арку рядом с магазином “ВкусВилл”. Во дворе поверните налево и идите вдоль сквера с детской площадкой. В конце сквера поверните направо — перед вами будет узкая дорожка, ведущая ко второму подъезду.",
    "Подойдите ко второму подъезду, затем поверните налево. Наша дверь находится в самом углу жёлтого углового здания. На двери вы увидите вывеску “Фловей”.",
    "Если по пути возникнут вопросы или не получится найти вход, просто напишите или позвоните нам — мы с удовольствием подскажем и поможем найти дорогу.",
  ],
  legalEntity: 'ООО "МЭДЖИК ГАРДЕН"',
  inn: "9704145697",
  ogrn: "1227700362779",
};

export interface SocialLink {
  label: string;
  href: string;
  /** Set for platforms that legally require a disclaimer in Russia (e.g. Meta-owned apps). */
  disclaimer?: string;
}

// TODO: уточнить реальные ссылки на соцсети у заказчика (в макете иконки
// Telegram/VK/Instagram-подобная показаны без подписанных URL).
export const socialLinks: SocialLink[] = [
  { label: "Telegram", href: contactInfo.telegramUrl },
  { label: "VK", href: "https://vk.com/floway" },
  {
    label: "Instagram",
    href: "https://instagram.com/floway",
    // Ссылка на баннере не называет саму соцсеть — просто указывает на
    // компанию-владельца, признанную экстремистской и запрещённую в РФ.
    disclaimer:
      "Принадлежит компании Meta, признанной экстремистской организацией и запрещённой на территории РФ.",
  },
];
