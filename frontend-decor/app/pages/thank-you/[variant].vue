<script setup lang="ts">
import { socialIcons } from "~/constants/social-icons";
import type { ThankYouPageVariant } from "~/types/api";

/**
 * Where ApplyForm.vue navigates to after a successful submission. Decor
 * site has exactly one variant ("decor", seeded by migration 00061) —
 * unlike the school's four, kept as a dynamic route only to reuse the same
 * ApplyForm.vue navigateTo pattern unchanged.
 */
const validVariants: ThankYouPageVariant[] = ["decor"];

const route = useRoute();
const variantParam = route.params.variant as string;
if (!validVariants.includes(variantParam as ThankYouPageVariant)) {
  throw createError({ statusCode: 404, statusMessage: "Страница не найдена" });
}
const variant = variantParam as ThankYouPageVariant;

const api = useApi();
const { data: page } = await useAsyncData(`thank-you-page-${variant}`, () =>
  api.getThankYouPage(variant),
);
if (!page.value) {
  throw createError({ statusCode: 404, statusMessage: "Страница не найдена" });
}

useSeoMeta({ title: () => `${page.value?.title ?? "Спасибо за заявку"} — flo-way` });

const { contactChannels } = await useContactChannels();
const { socialLinks } = await useSocialLinks();

const faqItems = computed(
  () => page.value?.faqItems.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);
</script>

<template>
  <UiGlassPage v-if="page">
    <div class="flex flex-col gap-24 lg:gap-64">
      <div class="flex flex-col gap-24 lg:flex-row lg:items-stretch lg:gap-64">
        <div class="order-2 flex flex-col items-start gap-24 lg:order-1 lg:w-1/2">
          <h1 class="font-display text-h1 text-ink">{{ page.title }}</h1>
          <MarkdownContent
            v-if="page.description"
            :source="page.description"
            class="w-full whitespace-pre-line font-body text-body text-ink"
          />
          <div class="mt-auto flex w-full flex-col gap-16">
            <UiButton variant="primary" to="/">Вернуться на главную</UiButton>
          </div>
        </div>
        <div class="order-1 lg:order-2 lg:w-1/2">
          <UiContentImage
            v-if="page.heroImage"
            :src="resolveOptimizedMediaUrl(page.heroImage)"
            :alt="page.title"
            class="aspect-square w-full rounded-lg object-cover"
            sizes="400:100vw lg:480px"
          />
          <UiMediaPlaceholder v-else aspect="1/1" />
        </div>
      </div>

      <p v-if="page.subtitle" class="text-center font-display text-h2 text-primary">
        {{ page.subtitle }}
      </p>

      <div class="mx-auto flex w-full max-w-[720px] flex-col items-center gap-24 text-center">
        <div
          v-if="page.showMessengers && contactChannels.length"
          class="flex flex-col items-center gap-12"
        >
          <p class="font-body text-body text-ink">Срочный вопрос? Напишите нам:</p>
          <div class="flex gap-16">
            <a
              v-for="channel in contactChannels"
              :key="channel.label"
              :href="channel.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="channel.label"
              class="grid size-[44px] place-items-center rounded-full bg-primary text-white hover:opacity-80"
            >
              <component :is="channel.icon" class="size-[20px]" aria-hidden="true" />
            </a>
          </div>
        </div>

        <div
          v-if="page.showSocialLinks && socialLinks.length"
          class="flex flex-col items-center gap-12"
        >
          <p class="font-body text-body text-ink">Пока ждёте — загляните в соцсети:</p>
          <div class="flex gap-16">
            <a
              v-for="social in socialLinks"
              :key="social.label"
              :href="social.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="social.label"
              class="grid size-[44px] place-items-center rounded-full bg-primary text-white hover:opacity-80"
            >
              <component :is="socialIcons[social.label]" class="size-[20px]" aria-hidden="true" />
            </a>
          </div>
        </div>
      </div>

      <FaqSection
        v-if="page.showFaq && faqItems.length"
        title="Частые вопросы"
        description=""
        :items="faqItems"
      />
    </div>
  </UiGlassPage>
</template>
