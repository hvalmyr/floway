import type { Component } from "vue";
import IconInstagram from "~/components/ui/IconInstagram.vue";
import IconTelegram from "~/components/ui/IconTelegram.vue";
import IconVk from "~/components/ui/IconVk.vue";

/**
 * Icon per social_links.label (see useSocialLinks()) — used by AppFooter and
 * (when enabled) a thank-you page. A label an admin adds outside this fixed
 * set (Telegram/VK/Instagram) simply renders no icon.
 */
export const socialIcons: Record<string, Component> = {
  Telegram: IconTelegram,
  VK: IconVk,
  Instagram: IconInstagram,
};
