<script setup lang="ts">
import { contactInfo } from "~/constants/contact-info";
import { socialIcons } from "~/constants/social-icons";

useSeoMeta({
  title: "Контакты — flo-way",
  description: "Телефон, почта и мессенджеры flo-way — новогоднее оформление в Москве и области.",
});

const { text } = await usePageContent();
const { socialLinks } = await useSocialLinks();

const phone = computed(() => text("contact_phone", contactInfo.phone));
const phoneHref = computed(() => `tel:${phone.value.replace(/[^\d+]/g, "")}`);
</script>

<template>
  <div class="container flex flex-col gap-64 py-48 sm:py-64 lg:py-80">
    <div
      class="flex w-full flex-col gap-16 rounded-md bg-white/55 p-24 backdrop-blur backdrop-saturate-150"
    >
      <h1 class="font-display text-h2 text-ink">Контакты</h1>

      <p class="font-body text-body text-ink">
        Позвонить: <a :href="phoneHref" class="text-primary underline">{{ phone }}</a>
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
        >
        или
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
      <p class="font-body text-body text-ink">
        География работы: {{ text("contact_service_area", contactInfo.serviceArea) }}
      </p>
    </div>

    <div
      v-if="socialLinks.length"
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
