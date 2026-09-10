<script setup lang="ts">
import { socialIcons } from "~/constants/social-icons";
import type { ThankYouPageVariant } from "~/types/api";

/**
 * Where ApplyForm.vue navigates to after a successful submission — content
 * is fully admin-editable per variant (course/masterclass/trial_lesson) via
 * /admin/thank-you-page. Messenger icons and social links are the site's
 * existing global contact channels (see ThankYouPage's doc comment in
 * types/api.ts), just toggled on/off per variant.
 */
const validVariants: ThankYouPageVariant[] = ["course", "masterclass", "trial_lesson"];

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

useSeoMeta({ title: () => `${page.value?.title ?? "Спасибо за заявку"} — Фловей` });

const { contactChannels } = await useContactChannels();
const { socialLinks } = await useSocialLinks();

const carouselPhotos = computed(
  () => page.value?.photos.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);
const faqItems = computed(
  () => page.value?.faqItems.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);
</script>

<template>
  <UiGlassPage v-if="page">
    <div class="flex flex-col gap-24 lg:gap-64">
      <!-- Hero block matches the site's shared Hero.vue layout (title/
      description/CTAs left, 1:1 photo right, media-first on mobile) —
      reimplemented here rather than reusing that component since it's a
      standalone <section> with its own container/py padding meant to sit
      directly on the page background, while this page's content lives
      inside UiGlassPage's own container+padding (same reasoning as
      pages/blog/[slug].vue's hero). -->
      <div class="flex flex-col gap-24 lg:flex-row lg:items-stretch lg:gap-64">
        <div class="order-2 flex flex-col items-start gap-24 lg:order-1 lg:w-1/2">
          <h1 class="font-display text-h1 text-ink">{{ page.title }}</h1>
          <MarkdownContent
            v-if="page.description"
            :source="page.description"
            class="w-full whitespace-pre-line font-body text-body text-ink"
          />
          <div class="mt-auto flex w-full flex-col gap-16">
            <UiButton
              v-if="page.showBlogLink && page.blogLinkUrl"
              variant="outline"
              :to="page.blogLinkUrl"
            >
              {{ page.blogLinkText || "Читать блог" }}
            </UiButton>
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

      <!-- Subtitle is its own section, right after the hero. -->
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

        <UiButton
          v-if="page.showCommunity && page.communityUrl"
          variant="outline"
          :to="page.communityUrl"
        >
          {{ page.communityText || "Присоединиться к сообществу" }}
        </UiButton>
      </div>

      <PhotoCarousel v-if="page.showCarousel && carouselPhotos.length" :photos="carouselPhotos" />

      <FaqSection
        v-if="page.showFaq && faqItems.length"
        title="Частые вопросы"
        description=""
        :items="faqItems"
      />
    </div>
  </UiGlassPage>
</template>
