/**
 * Placeholder contact details for flo-way.ru — real values are the
 * client's own (п. 3 ТЗ: "домен, контакты, юридические данные [ ]"), edited
 * via page_content once known. Unlike the school, decor has no walk-in
 * studio address — the service is on-site at the client's own property, so
 * there's no address/metro/directions here, just phone/messengers and the
 * service area (п. 1 ТЗ: Москва и Московская область).
 */
export const contactInfo = {
  phone: "+7 000 000 00 00",
  email: "info@flo-way.ru",
  telegramUrl: "https://t.me/flo_way",
  whatsappUrl: "https://wa.me/70000000000",
  maxUrl: "",
  serviceArea: "Москва и Московская область",
  legalEntity: "",
  inn: "",
  ogrn: "",
};

export interface SocialLink {
  label: string;
  href: string;
  disclaimer?: string;
}

// Fallback for useSocialLinks() when /api/v1/social-links is empty or
// unreachable — the real, admin-editable list lives in the social_links
// table (see /admin/social-links).
export const socialLinks: SocialLink[] = [
  { label: "Telegram", href: contactInfo.telegramUrl },
  {
    label: "Instagram",
    href: "",
    disclaimer:
      "Принадлежит компании Meta, признанной экстремистской организацией и запрещённой на территории РФ.",
  },
];
