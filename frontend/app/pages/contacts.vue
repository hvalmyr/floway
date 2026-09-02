<script setup lang="ts">
import type { Component } from "vue";
import IconInstagram from "~/components/ui/IconInstagram.vue";
import IconTelegram from "~/components/ui/IconTelegram.vue";
import IconVk from "~/components/ui/IconVk.vue";
import { contactInfo, socialLinks } from "~/constants/contact-info";

useSeoMeta({
  title: "Контакты — Фловей",
  description: "Телефон, почта, мессенджеры и адрес школы флористики “Фловей” в Москве.",
});

const { text } = await usePageContent();

const socialIcons: Record<string, Component> = {
  Telegram: IconTelegram,
  VK: IconVk,
  Instagram: IconInstagram,
};

const phone = computed(() => text("contact_phone", contactInfo.phone));
const phoneHref = computed(() => `tel:${phone.value.replace(/[^\d+]/g, "")}`);
const metroStations = computed(() =>
  text("contact_metro_stations", contactInfo.metroStations.join(", "))
    .split(",")
    .map((station) => station.trim())
    .filter(Boolean),
);
</script>

<template>
  <div class="container flex flex-col gap-64 py-48 sm:py-64 lg:py-80">
    <!-- Одна колонка на всю ширину — h1 (визуально уменьшен до размера
    остальных заголовков) теперь внутри той же стеклянной карточки, что и
    способы связи, а не отдельной строкой над ней. -->
    <div
      class="flex w-full flex-col gap-16 rounded-md bg-white/55 p-24 backdrop-blur backdrop-saturate-150"
    >
      <h1 class="font-display text-h2 text-ink">Контакты</h1>

      <p class="font-body text-body text-ink">
        Позвонить:
        <a :href="phoneHref" class="text-primary underline">{{ phone }}</a>
      </p>
      <p class="font-body text-body text-ink">
        Написать на почту:
        <a
          :href="`mailto:${text('contact_email', contactInfo.email)}`"
          class="text-primary underline"
          >{{ text("contact_email", contactInfo.email) }}</a
        >
      </p>
      <p class="font-body text-body text-ink">
        Написать в
        <a
          :href="text('contact_telegram_url', contactInfo.telegramUrl)"
          class="text-primary underline"
          >Telegram</a
        >,
        <a
          :href="text('contact_whatsapp_url', contactInfo.whatsappUrl)"
          class="text-primary underline"
          >Whatsapp</a
        >
        <template v-if="text('contact_max_url', contactInfo.maxUrl)">
          или
          <a :href="text('contact_max_url', contactInfo.maxUrl)" class="text-primary underline"
            >Max</a
          >
        </template>
      </p>
    </div>

    <div
      class="flex w-full flex-col gap-16 rounded-md bg-white/55 p-24 backdrop-blur backdrop-saturate-150"
    >
      <h2 class="font-display text-h2 text-primary">Адрес</h2>

      <p class="font-body text-body text-ink">{{ text("contact_address", contactInfo.address) }}</p>
      <p class="font-body text-body text-ink">
        Ближайшие станции метро:
        <strong>{{ metroStations.slice(0, -1).join(", ") }}</strong>
        и
        <strong>{{ metroStations.at(-1) }}</strong
        >.
      </p>
      <MarkdownContent :source="text('contact_directions', contactInfo.directions.join('\n\n'))" />
    </div>

    <iframe
      :src="text('contact_map_iframe_url', '')"
      title="Карта проезда до школы «Фловей»"
      loading="lazy"
      width="100%"
      height="607"
      frameborder="0"
      class="w-full rounded-lg border-0"
    />

    <div
      class="flex w-full flex-col gap-24 rounded-md bg-white/55 p-24 backdrop-blur backdrop-saturate-150"
    >
      <h2 class="font-display text-h2 text-ink">Соцсети</h2>
      <div class="flex flex-col gap-8">
        <div class="flex gap-16">
          <a
            v-for="social in socialLinks"
            :key="social.label"
            :href="social.href"
            target="_blank"
            rel="noopener noreferrer"
            :aria-label="social.label"
            class="grid size-[44px] place-items-center rounded-full bg-primary text-white hover:opacity-90"
          >
            <component :is="socialIcons[social.label]" class="size-[20px]" aria-hidden="true" />
          </a>
        </div>
        <p
          v-for="social in socialLinks.filter((s) => s.disclaimer)"
          :key="`${social.label}-disclaimer`"
          class="font-body text-body text-ink/70"
        >
          {{ social.disclaimer }}
        </p>
      </div>
    </div>
  </div>
</template>
