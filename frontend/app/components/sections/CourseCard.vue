<script setup lang="ts">
/**
 * Course / course-block card used in the "Основы флористики" and
 * "Профильные курсы" sections. lessonsCount/hours/price may be undefined
 * when the client hasn't filled them in yet — rendered as "?" per the
 * mockup (not a guessed value).
 *
 * @example
 * <CourseCard
 *   title="Блок «Букеты»"
 *   variant="surface-primary"
 *   :lessons-count="7"
 *   :hours="30"
 *   :price="38500"
 *   to="/courses/osnovy-floristiki"
 * />
 */
const props = withDefaults(
  defineProps<{
    title: string;
    description?: string;
    variant?: "surface-primary" | "surface-ink" | "surface-white";
    coverImage?: string;
    lessonsCount?: number;
    hours?: number;
    price?: number;
    to: string;
    ctaLabel?: string;
  }>(),
  {
    variant: "surface-white",
    ctaLabel: "Узнать больше",
  },
);

const ctaVariant = computed(() =>
  props.variant === "surface-white" ? "primary" : "outline-inverted",
);
</script>

<template>
  <UiCard :variant="variant" class="flex h-full flex-col">
    <template #media>
      <img
        v-if="coverImage"
        :src="resolveMediaUrl(coverImage)"
        :alt="title"
        class="mx-auto aspect-square w-full max-w-[200px] rounded-sm object-cover"
      />
      <div
        v-else
        class="mx-auto aspect-square w-full max-w-[200px] rounded-sm border-2 border-current"
      />
    </template>
    <template #title>{{ title }}</template>

    <p class="mb-16 font-body text-body">
      {{ lessonsCount ?? "?" }} занятий, {{ hours ?? "?" }} часов
    </p>
    <p v-if="description" class="mb-24 font-body text-body">{{ description }}</p>
    <p v-if="price" class="mb-24 font-display text-h4">{{ price.toLocaleString("ru-RU") }} ₽</p>

    <UiButton :variant="ctaVariant" :to="to" class="mt-auto w-full justify-center">{{
      ctaLabel
    }}</UiButton>
  </UiCard>
</template>
