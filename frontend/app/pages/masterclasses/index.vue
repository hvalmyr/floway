<script setup lang="ts">
import { masterclassDisplayStyleCycle } from "~/constants/display-style-colors";

useSeoMeta({
  title: "Мастер-классы по флористике в Москве — Фловей",
  description:
    "Мастер-классы по флористике в свободном графике: букеты и композиции, все материалы включены.",
});

const api = useApi();
const { text } = await usePageContent();
const { data: masterclasses } = await useAsyncData("masterclasses-list", () =>
  api.getMasterClasses(),
);

const { data: featuresData } = await useAsyncData("masterclasses-features", () =>
  api.getFeatures("masterclasses"),
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
        <UiButton variant="primary" to="#masterclasses-list">Мастер-классы</UiButton>
        <UiButton variant="outline" to="#apply">Оставить заявку</UiButton>
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
          <ApplyForm
            context="masterclass"
            :related-slug="selectedSlug"
            title="Оставить заявку на мастер-класс"
          />
        </div>
      </div>
    </section>
  </div>
</template>
