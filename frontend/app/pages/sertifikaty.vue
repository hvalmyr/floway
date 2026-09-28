<script setup lang="ts">
import type { SiteButtonVariant } from "~/types/api";

useSeoMeta({
  title: "Подарочный мастер-класс — Фловей",
  description:
    "Подарите мастер-класс по флористике: сертификат школы «Фловей» на любую сумму, курс или мастер-класс.",
});

const api = useApi();
// Fired together instead of one `await` per call — these 3 backend
// endpoints are independent of each other and of page-content, so awaiting
// them in sequence only added round-trips before SSR could respond.
const featuresAsync = useAsyncData("gift-certificate-features", () =>
  api.getFeatures("gift_certificate"),
);
const faqAsync = useAsyncData("gift-certificate-faq", () => api.getPageFaq("gift_certificate"));
const carouselPhotosAsync = useAsyncData("gift-certificate-carousel-photos", () =>
  api.getGiftCertificateCarouselPhotos(),
);
const siteButtonsAsync = useAsyncData("site-buttons", () => api.getSiteButtons());

const [{ text }] = await Promise.all([
  usePageContent(),
  featuresAsync,
  faqAsync,
  carouselPhotosAsync,
  siteButtonsAsync,
]);

const { data: featuresData } = featuresAsync;
const { data: faq } = faqAsync;
const { data: carouselPhotosData } = carouselPhotosAsync;
const carouselPhotos = computed(
  () => carouselPhotosData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const advantages = computed(
  () =>
    featuresData.value
      ?.slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((f) => ({ icon: f.icon, title: f.title, description: f.description })) ?? [],
);

// Static Hero CTA buttons — text/style/link editable at /admin/site-buttons.
const { data: siteButtonsData } = siteButtonsAsync;
const siteButtonsByKey = computed(
  () => new Map((siteButtonsData.value ?? []).map((b) => [b.key, b])),
);
function siteButton(
  key: string,
  fallback: { text: string; variant: SiteButtonVariant; url: string },
) {
  const found = siteButtonsByKey.value.get(key);
  return found ? { text: found.text, variant: found.variant, url: found.url } : fallback;
}
const applyButton = computed(() =>
  siteButton("sertifikaty_apply", { text: "Оставить заявку", variant: "primary", url: "#apply" }),
);
const masterclassesButton = computed(() =>
  siteButton("sertifikaty_masterclasses", {
    text: "Мастер-классы",
    variant: "outline",
    url: "/masterclasses",
  }),
);
</script>

<template>
  <div>
    <Hero>
      <template #title>{{
        text("gift_certificate_hero_title", "Подарите мастер-класс по флористике")
      }}</template>
      <template #lead>
        {{
          text(
            "gift_certificate_hero_lead",
            "Сертификат «Фловей» — оригинальный подарок для тех, кто любит цветы и творчество. Получатель сам выберет мастер-класс и удобную дату визита в школу.",
          )
        }}
      </template>
      <template #actions>
        <UiButton :variant="applyButton.variant" :to="applyButton.url">{{
          applyButton.text
        }}</UiButton>
        <UiButton :variant="masterclassesButton.variant" :to="masterclassesButton.url">{{
          masterclassesButton.text
        }}</UiButton>
      </template>
      <template v-if="text('gift_certificate_hero_image')" #media>
        <UiHeroPicture
          :src="resolveOptimizedMediaUrl(text('gift_certificate_hero_image'))"
          alt=""
        />
      </template>
    </Hero>

    <section v-if="carouselPhotos.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container">
        <PhotoCarousel :photos="carouselPhotos" />
      </div>
    </section>

    <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>
          {{ text("gift_certificate_features_heading", "Почему это отличный подарок") }}
          <template #lead>{{
            text(
              "gift_certificate_features_lead",
              "Дарите не вещь, а впечатление — тёплый опыт создания своими руками.",
            )
          }}</template>
        </SectionHeading>
        <FeatureGrid :items="advantages" />
      </div>
    </section>

    <section
      id="apply"
      class="scroll-mt-64 bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:scroll-mt-96 lg:py-80"
    >
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <LazyApplyForm
            context="gift_certificate"
            :title="
              text('gift_certificate_apply_form_title', 'Оставить заявку на подарочный сертификат')
            "
            :lead="
              text('gift_certificate_apply_form_lead', 'Свяжемся с вами и оформим сертификат.')
            "
            hydrate-on-visible
          />
        </div>
      </div>
    </section>

    <section v-if="faq?.visible && faq.items.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container">
        <FaqSection
          :title="faq.title || 'Вопросы и ответы'"
          :description="faq.description"
          :items="faq.items"
        />
      </div>
    </section>
  </div>
</template>
