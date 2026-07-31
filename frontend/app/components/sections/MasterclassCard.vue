<script setup lang="ts">
import type { Masterclass } from "~/types/api";

/**
 * Photo + text block for a masterclass — used both as a repeating teaser on
 * /masterclasses and as the main content block of /masterclasses/[slug].
 * `linkToDetail` shows the "узнать подробнее" CTA (listing use); omit it on
 * the detail page itself, which has its own <ApplyForm> below instead.
 *
 * @example
 * <MasterclassCard :masterclass="mc" link-to-detail />
 */
const props = withDefaults(
  defineProps<{
    masterclass: Masterclass;
    linkToDetail?: boolean;
  }>(),
  { linkToDetail: false },
);

const infoItems = computed(() => {
  const items: { label: string; value: string }[] = [
    { label: "Длительность", value: props.masterclass.duration ?? "уточняется" },
  ];

  items.push({
    label: "Цена (группа)",
    value: props.masterclass.priceGroup != null ? `от ${props.masterclass.priceGroup.toLocaleString("ru-RU")} ₽` : "уточняется",
  });

  if (props.masterclass.priceIndividual != null) {
    items.push({
      label: "Индивидуально",
      value: `от ${props.masterclass.priceIndividual.toLocaleString("ru-RU")} ₽`,
    });
  }

  return items;
});
</script>

<template>
  <article class="grid grid-cols-1 gap-24 lg:grid-cols-[38%_1fr] lg:gap-40">
    <!-- TODO: заменить на один <NuxtImg> с `sizes` под брейкпоинты, когда появятся
    реальные фото — сейчас два плейсхолдера переключаются видимостью, чтобы не
    завязывать компонент на responsive-варианты одного prop`а. -->
    <UiMediaPlaceholder aspect="4/3" class="lg:hidden" />
    <UiMediaPlaceholder aspect="3/4" class="hidden lg:block" />

    <div class="flex flex-col items-start gap-16">
      <h2 class="font-display text-h2 text-ink-900">{{ masterclass.title }}</h2>
      <p class="text-body text-ink-700">{{ masterclass.shortDescription }}</p>
      <p v-if="masterclass.fullDescription" class="text-body text-ink-700">{{ masterclass.fullDescription }}</p>
      <p v-if="masterclass.endingText" class="text-body text-ink-700">{{ masterclass.endingText }}</p>

      <UiInfoRow :items="infoItems" class="w-full" />

      <UiButton v-if="linkToDetail" variant="outline" :to="`/masterclasses/${masterclass.slug}`">
        Узнать подробнее
      </UiButton>

      <p v-if="masterclass.priceDescription" class="text-small text-ink-400">{{ masterclass.priceDescription }}</p>
    </div>
  </article>
</template>
