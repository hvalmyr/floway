<script setup lang="ts">
import { masterclassDisplayStyleCycle } from "~/constants/display-style-colors";
import type { SiteButtonVariant } from "~/types/api";

useSeoMeta({
  title: "Мастер-классы по флористике в Москве — Фловей",
  description:
    "Мастер-классы по флористике в свободном графике: букеты и композиции, все материалы включены.",
});

const api = useApi();
// The hero image comes from page-content alone, but the other 3 fetches are
// independent of it and of each other — firing them together instead of
// awaiting one after another avoids a 4-deep sequential chain before SSR
// can respond.
const masterclassesAsync = useAsyncData("masterclasses-list", () => api.getMasterClasses());
const featuresAsync = useAsyncData("masterclasses-features", () =>
  api.getFeatures("masterclasses"),
);
const faqAsync = useAsyncData("masterclasses-faq", () => api.getPageFaq("masterclasses"));
const siteButtonsAsync = useAsyncData("site-buttons", () => api.getSiteButtons());

const [{ text }] = await Promise.all([
  usePageContent(),
  masterclassesAsync,
  featuresAsync,
  faqAsync,
  siteButtonsAsync,
]);

const { data: masterclasses } = masterclassesAsync;
const { data: featuresData } = featuresAsync;
const { data: faq } = faqAsync;
const { data: siteButtonsData } = siteButtonsAsync;
function siteButton(
  key: string,
  fallback: { text: string; variant: SiteButtonVariant; url: string },
) {
  const found = siteButtonsData.value?.find((b) => b.key === key);
  return found ? { text: found.text, variant: found.variant, url: found.url } : fallback;
}
const listButton = computed(() =>
  siteButton("masterclasses_hero_list", {
    text: "Мастер-классы",
    variant: "primary",
    url: "#masterclasses-list",
  }),
);
const applyButton = computed(() =>
  siteButton("masterclasses_hero_apply", {
    text: "Оставить заявку",
    variant: "outline",
    url: "#apply",
  }),
);

// Set by whichever MasterclassCard's "Записаться" was clicked last — there's
// one shared ApplyForm below the whole list (not one per card), so this is
// the only way the lead ends up tagged with which masterclass it was about.
const selectedSlug = ref<string | undefined>(undefined);
const features = computed(
  () =>
    featuresData.value
      ?.slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((f) => ({ icon: f.icon, title: f.title, description: f.description })) ?? [],
);
</script>

<template>
  <div>
    <Hero>
      <template #title>{{
        text("masterclasses_hero_title", "Мастер-классы по флористике в свободном графике")
      }}</template>
      <template #lead>
        {{
          text(
            "masterclasses_hero_lead",
            "Мастер классы по флористике посвящены созданию разных видов букетов и флористических композиций. На каждом занятии вы создаете собственную работу и осваиваете новые приемы и навыки флористики.",
          )
        }}
      </template>
      <template #actions>
        <UiButton :variant="listButton.variant" :to="listButton.url">{{
          listButton.text
        }}</UiButton>
        <UiButton :variant="applyButton.variant" :to="applyButton.url">{{
          applyButton.text
        }}</UiButton>
      </template>
      <template v-if="text('masterclasses_hero_image')" #media>
        <UiHeroPicture :src="resolveOptimizedMediaUrl(text('masterclasses_hero_image'))" alt="" />
      </template>
    </Hero>

    <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>
          {{ text("masterclasses_features_heading", "Почему стоит выбрать мастер-класс “Фловей”") }}
          <template #lead>{{
            text(
              "masterclasses_features_lead",
              "Разовое занятие без обязательств: приходите, когда удобно, и уходите с готовой работой в руках.",
            )
          }}</template>
        </SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <MasterclassesMarquee />

    <section id="masterclasses-list" class="scroll-mt-64 py-48 sm:py-64 lg:scroll-mt-96 lg:py-80">
      <div class="container flex flex-col gap-40 lg:gap-48">
        <MasterclassCard
          v-for="(mc, i) in masterclasses"
          :key="mc.id"
          :masterclass="mc"
          :display-style="masterclassDisplayStyleCycle[i % masterclassDisplayStyleCycle.length]"
          @apply="selectedSlug = mc.slug"
        />
      </div>
    </section>

    <section
      id="apply"
      class="scroll-mt-64 bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:scroll-mt-96 lg:py-80"
    >
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <LazyApplyForm
            context="masterclass"
            :related-slug="selectedSlug"
            :title="text('masterclasses_apply_form_title', 'Оставить заявку на мастер-класс')"
            :lead="text('masterclasses_apply_form_lead', '')"
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
