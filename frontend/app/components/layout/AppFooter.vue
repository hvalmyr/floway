<script setup lang="ts">
import { Instagram, Mail, Phone, Send } from "lucide-vue-next";
import type { Component } from "vue";
import { contactInfo, socialLinks } from "~/constants/contact-info";
import IconVk from "~/components/ui/IconVk.vue";

const year = new Date().getFullYear();

const socialIcons: Record<string, Component> = {
  Telegram: Send,
  VK: IconVk,
  Instagram,
};

const schoolLinks = [
  { to: "/#courses", label: "Курсы" },
  { to: "/masterclasses", label: "Мастер-классы" },
  // TODO: секции отзывов пока нет на главной (в брифе для главной страницы
  // отзывы не описаны отдельным блоком) — ссылка ведёт на якорь, которого
  // ещё не существует.
  { to: "/#reviews", label: "Отзывы" },
  { to: "/contacts", label: "Контакты" },
];

const documentLinks = [
  { to: "/privacy", label: "Политика конфиденциальности" },
  { to: "/legal-info", label: "Реквизиты ООО" },
];
</script>

<template>
  <footer class="mt-auto bg-primary-500 text-white">
    <div class="container flex flex-col gap-32 py-64">
      <NuxtLink to="/" class="font-display text-h3 text-white">ФлоВей</NuxtLink>

      <div class="grid grid-cols-1 gap-32 md:grid-cols-3">
        <nav class="flex flex-col gap-16" aria-label="Школа">
          <h3 class="font-display text-h4 text-white">Школа</h3>
          <NuxtLink v-for="link in schoolLinks" :key="link.to" :to="link.to" class="text-body text-white hover:underline">
            {{ link.label }}
          </NuxtLink>
        </nav>

        <nav class="flex flex-col gap-16" aria-label="Документы">
          <h3 class="font-display text-h4 text-white">Документы</h3>
          <NuxtLink
            v-for="link in documentLinks"
            :key="link.to"
            :to="link.to"
            class="text-body text-white hover:underline"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>

        <div class="flex flex-col gap-16">
          <h3 class="font-display text-h4 text-white">Связь со школой</h3>
          <a :href="contactInfo.phoneHref" class="flex items-center gap-8 text-body text-white hover:underline">
            <Phone class="size-[20px] shrink-0" aria-hidden="true" />
            {{ contactInfo.phone }}
          </a>
          <a :href="`mailto:${contactInfo.email}`" class="flex items-center gap-8 text-body text-white hover:underline">
            <Mail class="size-[20px] shrink-0" aria-hidden="true" />
            {{ contactInfo.email }}
          </a>

          <div class="flex gap-16 pt-8">
            <a
              v-for="social in socialLinks"
              :key="social.label"
              :href="social.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="social.label"
              class="grid size-[44px] place-items-center rounded-full border border-white/40 text-white hover:bg-white/10"
            >
              <component :is="socialIcons[social.label]" class="size-[20px]" aria-hidden="true" />
            </a>
          </div>
          <p v-for="social in socialLinks.filter((s) => s.disclaimer)" :key="`${social.label}-disclaimer`" class="text-small text-white/70">
            {{ social.disclaimer }}
          </p>
        </div>
      </div>
    </div>

    <div class="border-t border-white/25">
      <p class="container py-24 text-small text-white/80">&copy; {{ year }} ФлоВей — школа флористики</p>
    </div>
  </footer>
</template>
