import type { Component } from "vue";
import IconMax from "~/components/ui/IconMax.vue";
import IconTelegram from "~/components/ui/IconTelegram.vue";
import IconWhatsapp from "~/components/ui/IconWhatsapp.vue";
import { contactInfo } from "~/constants/contact-info";

export interface ContactChannel {
  label: string;
  href: string;
  icon: Component;
}

/**
 * "Связь со школой" — the messenger icons shown in AppFooter and (when
 * enabled) on a thank-you page (ThankYouPage.showMessengers). Same three
 * channels as ApplyForm's contactMethod options, sourced from the same
 * admin-editable page_content keys (/admin/page-content/info) — not a
 * separate per-surface list, so editing one place updates every surface
 * that shows them.
 *
 * @example
 * const { contactChannels } = await useContactChannels();
 */
export async function useContactChannels() {
  const { text } = await usePageContent();

  const contactChannels = computed<ContactChannel[]>(() =>
    [
      {
        label: "Telegram",
        href: text("contact_telegram_url", contactInfo.telegramUrl),
        icon: IconTelegram,
      },
      {
        label: "Whatsapp",
        href: text("contact_whatsapp_url", contactInfo.whatsappUrl),
        icon: IconWhatsapp,
      },
      { label: "Max", href: text("contact_max_url", contactInfo.maxUrl), icon: IconMax },
    ].filter((channel) => channel.href),
  );

  return { contactChannels };
}
