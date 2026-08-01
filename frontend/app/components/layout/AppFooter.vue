<script setup lang="ts">
import type { Component } from "vue";
import IconInstagram from "~/components/ui/IconInstagram.vue";
import IconMax from "~/components/ui/IconMax.vue";
import IconTelegram from "~/components/ui/IconTelegram.vue";
import IconVk from "~/components/ui/IconVk.vue";
import IconWhatsapp from "~/components/ui/IconWhatsapp.vue";
import { contactInfo, socialLinks } from "~/constants/contact-info";

const year = new Date().getFullYear();

const socialIcons: Record<string, Component> = {
  Telegram: IconTelegram,
  VK: IconVk,
  Instagram: IconInstagram,
};

// "Связь со школой": те же мессенджеры, что и в форме заявки (contactMethod),
// плюс Telegram уже есть в "соцсети" ниже.
const contactChannels = [
  { label: "Telegram", href: contactInfo.telegramUrl, icon: IconTelegram },
  { label: "Whatsapp", href: contactInfo.whatsappUrl, icon: IconWhatsapp },
  { label: "Max", href: contactInfo.maxUrl || undefined, icon: IconMax },
];

const schoolLinks = [
  { to: "/#courses", label: "Курсы" },
  { to: "/masterclasses", label: "Мастерклассы" },
  // TODO: секции отзывов пока нет на главной (в брифе для главной страницы
  // отзывы не описаны отдельным блоком) — ссылка ведёт на якорь, которого
  // ещё не существует.
  { to: "/#reviews", label: "Отзывы" },
  { to: "/contacts", label: "Контакты" },
];
</script>

<template>
  <footer class="mt-auto bg-primary text-white">
    <div class="container flex flex-col gap-32 py-64">
      <NuxtLink to="/">
        <LogoFloway class="h-48 w-auto text-white" />
      </NuxtLink>

      <div class="grid grid-cols-1 gap-32 md:grid-cols-3">
        <div class="flex flex-col gap-16">
          <h3 class="font-display text-h4 text-white">связь с школой</h3>
          <div class="flex gap-16">
            <a
              v-for="channel in contactChannels"
              :key="channel.label"
              :href="channel.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="channel.label"
              class="grid size-[44px] place-items-center rounded-full bg-white text-primary hover:opacity-80"
            >
              <component :is="channel.icon" class="size-[20px]" aria-hidden="true" />
            </a>
          </div>

          <h3 class="mt-16 font-display text-h4 text-white">соцсети</h3>
          <div class="flex gap-16">
            <a
              v-for="social in socialLinks"
              :key="social.label"
              :href="social.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="social.label"
              class="grid size-[44px] place-items-center rounded-full bg-white text-primary hover:opacity-80"
            >
              <component :is="socialIcons[social.label]" class="size-[20px]" aria-hidden="true" />
            </a>
          </div>
          <p
            v-for="social in socialLinks.filter((s) => s.disclaimer)"
            :key="`${social.label}-disclaimer`"
            class="font-body text-small text-white/70"
          >
            {{ social.disclaimer }}
          </p>
        </div>

        <nav class="flex flex-col gap-16" aria-label="Школа">
          <h3 class="font-display text-h4 text-white">школа</h3>
          <NuxtLink
            v-for="link in schoolLinks"
            :key="link.to"
            :to="link.to"
            class="font-body text-body text-white hover:underline"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>

        <div class="flex flex-col gap-16">
          <h3 class="font-display text-h4 text-white">документы</h3>
          <NuxtLink to="/privacy" class="font-body text-body text-white hover:underline">
            политика конфиденциальности
          </NuxtLink>
          <p class="font-body text-body text-white">{{ contactInfo.legalEntity }}</p>
          <p class="font-body text-body text-white">ИНН: {{ contactInfo.inn }}</p>
          <p class="font-body text-body text-white">ОГРН: {{ contactInfo.ogrn }}</p>
        </div>
      </div>
    </div>

    <div class="border-t border-white/25">
      <p class="container py-24 font-body text-small text-white/80">
        &copy; {{ year }} ФлоВей — школа флористики
      </p>
    </div>
  </footer>
</template>
