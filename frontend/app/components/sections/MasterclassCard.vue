<script setup lang="ts">
import type { Masterclass } from "~/types/api";

/**
 * Photo + text block for a masterclass — used both as a repeating teaser on
 * /masterclasses and as the main content block of /masterclasses/[slug].
 * `linkToDetail` shows the "узнать подробнее" CTA (listing use); omit it on
 * the detail page itself, which has its own <ApplyForm> below instead.
 * `accent` switches the title between blue and dark ink — the listing
 * alternates it per item, matching the mockup.
 *
 * @example
 * <MasterclassCard :masterclass="mc" link-to-detail accent />
 */
const props = withDefaults(
  defineProps<{
    masterclass: Masterclass;
    linkToDetail?: boolean;
    accent?: boolean;
  }>(),
  { linkToDetail: false, accent: false },
);

const infoItems = computed(() => {
  const duration = props.masterclass.duration ?? "уточняется";

  if (props.masterclass.priceGroup == null) {
    return [duration, "цена уточняется"];
  }
  if (props.masterclass.priceIndividual != null) {
    return [
      duration,
      `${props.masterclass.priceGroup.toLocaleString("ru-RU")}₽ или ${props.masterclass.priceIndividual.toLocaleString("ru-RU")}₽ [*]`,
    ];
  }
  return [duration, `${props.masterclass.priceGroup.toLocaleString("ru-RU")}₽`];
});
</script>

<template>
  <article class="grid grid-cols-1 gap-24 lg:grid-cols-[38%_1fr] lg:gap-40">
    <img
      v-if="masterclass.coverImage"
      :src="resolveMediaUrl(masterclass.coverImage)"
      :alt="masterclass.title"
      class="aspect-[4/3] w-full rounded-lg object-cover lg:aspect-[3/4]"
    />
    <UiMediaPlaceholder v-else aspect="4/3" class="lg:aspect-[3/4]" />

    <div class="flex flex-col items-start gap-16">
      <h2 class="font-display text-h2" :class="accent ? 'text-primary' : 'text-ink'">
        {{ masterclass.title }}
      </h2>
      <p class="w-full font-body text-body text-ink lg:w-4/5">{{ masterclass.shortDescription }}</p>
      <p v-if="masterclass.fullDescription" class="w-full font-body text-body text-ink lg:w-4/5">
        {{ masterclass.fullDescription }}
      </p>
      <p v-if="masterclass.endingText" class="w-full font-body text-body text-ink lg:w-4/5">
        {{ masterclass.endingText }}
      </p>

      <UiInfoRow :items="infoItems" class="w-full" />

      <UiButton v-if="linkToDetail" variant="primary" :to="`/masterclasses/${masterclass.slug}`">
        Узнать подробнее
      </UiButton>

      <p v-if="masterclass.priceDescription" class="font-body text-body text-ink/70">
        {{ masterclass.priceDescription }}
      </p>
    </div>
  </article>
</template>
