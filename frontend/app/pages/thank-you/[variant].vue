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
  <div v-if="page">
    <section class="py-48 sm:py-64 lg:py-80">
      <div class="container">
        <div class="mx-auto flex max-w-[720px] flex-col items-center gap-24 text-center">
          <h1 class="font-display text-h1 text-ink">{{ page.title }}</h1>
          <p v-if="page.subtitle" class="font-display text-h2 text-primary">
            {{ page.subtitle }}
          </p>
          <MarkdownContent
            v-if="page.description"
            :source="page.description"
            class="w-full whitespace-pre-line font-body text-body text-ink"
          />

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

          <div class="mt-8 flex w-full flex-col items-center gap-16 sm:flex-row sm:justify-center">
            <UiButton
              v-if="page.showBlogLink && page.blogLinkUrl"
              variant="outline"
              :to="page.blogLinkUrl"
            >
              {{ page.blogLinkText || "Читать блог" }}
            </UiButton>
            <UiButton
              v-if="page.showCommunity && page.communityUrl"
              variant="outline"
              :to="page.communityUrl"
            >
              {{ page.communityText || "Присоединиться к сообществу" }}
            </UiButton>
            <UiButton variant="primary" to="/">Вернуться на главную</UiButton>
          </div>
        </div>
      </div>
    </section>

    <section v-if="page.showCarousel && carouselPhotos.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container">
        <PhotoCarousel :photos="carouselPhotos" />
      </div>
    </section>

    <section
      v-if="page.showFaq && faqItems.length"
      class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80"
    >
      <div class="container">
        <FaqSection title="Частые вопросы" description="" :items="faqItems" />
      </div>
    </section>
  </div>
</template>
