<script setup lang="ts">
import type { Component } from "vue";
import IconInstagram from "~/components/ui/IconInstagram.vue";
import IconTelegram from "~/components/ui/IconTelegram.vue";
import IconVk from "~/components/ui/IconVk.vue";
import { contactInfo, socialLinks } from "~/constants/contact-info";

useSeoMeta({
  title: "Контакты — ФлоВей",
  description: "Телефон, почта, мессенджеры и адрес школы флористики “ФлоВей” в Москве.",
});

const { text } = await usePageContent();

const socialIcons: Record<string, Component> = {
  Telegram: IconTelegram,
  VK: IconVk,
  Instagram: IconInstagram,
};
</script>

<template>
  <div class="container flex flex-col gap-64 py-64 sm:py-96 lg:py-120">
    <div class="flex flex-col gap-16">
      <h1 class="font-display text-h1 text-ink">Контакты</h1>

      <p class="font-body text-body text-ink">
        Позвонить:
        <a :href="contactInfo.phoneHref" class="text-primary underline">{{ contactInfo.phone }}</a>
      </p>
      <p class="font-body text-body text-ink">
        Написать на почту:
        <a :href="`mailto:${contactInfo.email}`" class="text-primary underline">{{
          contactInfo.email
        }}</a>
      </p>
      <p class="font-body text-body text-ink">
        Написать в
        <a :href="contactInfo.telegramUrl" class="text-primary underline">Telegram</a>,
        <a :href="contactInfo.whatsappUrl" class="text-primary underline">Whatsapp</a>
        или
        <!-- TODO: реальная ссылка/контакт для Max ожидает данных от заказчика. -->
        <a href="#" class="text-primary underline">Max</a>
      </p>
    </div>

    <div class="grid grid-cols-1 gap-32 lg:grid-cols-2 lg:gap-64">
      <div class="flex flex-col gap-16">
        <h2 class="font-display text-h2 text-primary">Адрес</h2>

        <p class="font-body text-body text-ink">{{ contactInfo.address }}</p>
        <p class="font-body text-body text-ink">
          Ближайшие станции метро:
          <strong>{{ contactInfo.metroStations.slice(0, -1).join(", ") }}</strong>
          и
          <strong>{{ contactInfo.metroStations.at(-1) }}</strong
          >.
        </p>
        <p
          v-for="(paragraph, i) in contactInfo.directions"
          :key="i"
          class="font-body text-body text-ink"
        >
          {{ paragraph }}
        </p>
      </div>
      <img
        v-if="text('contacts_office_image')"
        :src="resolveMediaUrl(text('contacts_office_image'))"
        alt=""
        class="aspect-[4/3] w-full rounded-lg object-cover"
      />
      <UiMediaPlaceholder v-else aspect="4/3" />
    </div>

    <div class="flex flex-col gap-24">
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
