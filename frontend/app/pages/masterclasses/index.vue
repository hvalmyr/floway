<script setup lang="ts">
import type { Component } from "vue";
import { featureIconComponent } from "~/constants/feature-icons";

useSeoMeta({
  title: "Мастер-классы по флористике в Москве — ФлоВей",
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
const features = computed(
  () =>
    featuresData.value
      ?.slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((f) => ({
        icon: featureIconComponent(f.icon),
        title: f.title,
        description: f.description,
      }))
      .filter((f): f is typeof f & { icon: Component } => f.icon !== undefined) ?? [],
);
</script>

<template>
  <div>
    <Hero>
      <template #title>Мастер-классы по флористике в свободном графике</template>
      <template #actions>
        <UiButton variant="primary" size="lg" to="#apply" class="w-full sm:w-auto"
          >Оставить заявку</UiButton
        >
      </template>
      <template v-if="text('masterclasses_hero_image')" #media>
        <img
          :src="resolveMediaUrl(text('masterclasses_hero_image'))"
          alt=""
          class="aspect-[4/3] w-full rounded-lg object-cover"
        />
      </template>
    </Hero>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">Почему стоит выбрать мастер-класс «ФлоВей»</SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section class="bg-surface py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-64">
        <MasterclassCard
          v-for="(mc, i) in masterclasses"
          :key="mc.id"
          :masterclass="mc"
          link-to-detail
          :accent="i % 2 === 0"
        />
      </div>
    </section>

    <section id="apply" class="scroll-mt-64 py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm context="masterclass" title="Оставить заявку на мастер-класс" />
        </div>
      </div>
    </section>
  </div>
</template>
