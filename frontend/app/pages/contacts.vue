<script setup lang="ts">
import { Instagram, Mail, MapPin, MessageCircle, Phone, Send } from "lucide-vue-next";
import type { Component } from "vue";
import IconVk from "~/components/ui/IconVk.vue";
import { contactInfo, socialLinks } from "~/constants/contact-info";

useSeoMeta({
  title: "Контакты — ФлоВей",
  description: "Телефон, почта, мессенджеры и адрес школы флористики «ФлоВей» в Москве.",
});

const socialIcons: Record<string, Component> = {
  Telegram: Send,
  VK: IconVk,
  Instagram,
};

const contactMethods = [
  { icon: Phone, label: "Позвонить", value: contactInfo.phone, href: contactInfo.phoneHref },
  { icon: Mail, label: "Написать на почту", value: contactInfo.email, href: `mailto:${contactInfo.email}` },
  { icon: Send, label: "Написать в Telegram", value: "Telegram", href: contactInfo.telegramUrl },
  { icon: MessageCircle, label: "Написать в Whatsapp", value: "Whatsapp", href: contactInfo.whatsappUrl },
];
</script>

<template>
  <div class="container flex flex-col gap-64 py-64 sm:py-96 lg:py-120">
    <h1 class="font-display text-h1 text-ink-900">Контакты</h1>

    <section class="grid grid-cols-1 gap-24 md:grid-cols-2 lg:grid-cols-4">
      <a
        v-for="method in contactMethods"
        :key="method.label"
        :href="method.href"
        class="flex flex-col items-start gap-16 rounded-md bg-surface p-32 text-ink-900 hover:bg-primary-50"
      >
        <component :is="method.icon" class="size-40 text-primary-500" aria-hidden="true" />
        <span class="font-display text-h4">{{ method.label }}</span>
        <span class="text-body text-ink-700">{{ method.value }}</span>
      </a>
      <!-- TODO: контакт в Max мессенджере ожидает ссылки/номера от заказчика. -->
      <div class="flex flex-col items-start gap-16 rounded-md border border-dashed border-line p-32 text-ink-400">
        <MessageCircle class="size-40" aria-hidden="true" />
        <span class="font-display text-h4">Написать в Max</span>
        <span class="text-body">Контакт уточняется</span>
      </div>
    </section>

    <section class="grid grid-cols-1 gap-32 lg:grid-cols-2 lg:gap-64">
      <div class="flex flex-col gap-16">
        <h2 class="font-display text-h2 text-ink-900">Адрес</h2>
        <p class="flex items-start gap-12 text-body-l text-ink-900">
          <MapPin class="mt-4 size-24 shrink-0 text-primary-500" aria-hidden="true" />
          {{ contactInfo.address }}
        </p>
        <!-- Названия станций метро в макете выделены цветом — это разовая
        маркировка, не токен дизайн-системы (см. docs/floway-design.md §2.1),
        поэтому здесь обычный текст, не цветной. Сами названия без доступа к
        Figma не считать — TODO ниже. -->
        <p class="text-body text-ink-700">{{ contactInfo.metroNote }}</p>
        <p class="text-body text-ink-700">
          <!-- TODO: точный текст инструкции "как добраться" уточнить у заказчика — ниже общее плейсхолдер-описание. -->
          От ближайших станций метро — 5–10 минут пешком. У входа ориентируйтесь на вывеску школы «ФлоВей»: ресепшн
          подскажет, как пройти в студию.
        </p>
      </div>
      <UiMediaPlaceholder aspect="4/3" />
    </section>

    <section class="flex flex-col gap-24">
      <h2 class="font-display text-h2 text-ink-900">Карта</h2>
      <!-- TODO: вставить Яндекс.Карты iframe. Получить код на
      https://yandex.ru/map-widget/ → «Конструктор» → указать адрес школы →
      скопировать сгенерированный <iframe> и заменить плейсхолдер ниже. -->
      <div class="flex aspect-[16/9] w-full items-center justify-center rounded-lg bg-surface text-ink-400" aria-hidden="true">
        <MapPin class="size-40" aria-hidden="true" />
      </div>
    </section>

    <section class="flex flex-col gap-24">
      <h2 class="font-display text-h2 text-ink-900">Соцсети</h2>
      <div class="flex flex-col gap-8">
        <div class="flex gap-16">
          <a
            v-for="social in socialLinks"
            :key="social.label"
            :href="social.href"
            target="_blank"
            rel="noopener noreferrer"
            :aria-label="social.label"
            class="grid size-[44px] place-items-center rounded-full border-[1.5px] border-primary-500 text-primary-700 hover:bg-primary-50"
          >
            <component :is="socialIcons[social.label]" class="size-[20px]" aria-hidden="true" />
          </a>
        </div>
        <p
          v-for="social in socialLinks.filter((s) => s.disclaimer)"
          :key="`${social.label}-disclaimer`"
          class="text-small text-ink-400"
        >
          {{ social.disclaimer }}
        </p>
      </div>
    </section>
  </div>
</template>
