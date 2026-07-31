<script setup lang="ts">
/**
 * Course / course-block card used in the "Основы флористики" and
 * "Профильные курсы" sections. lessonsCount/hours/price may be undefined
 * when the client hasn't filled them in yet — rendered as "—" rather than a
 * guessed value.
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

// Тёмная карточка (surface-ink) получает белый контур кнопки; светлая
// голубая (surface-primary) и белая — тёмный: белый текст на #82B1CC не
// проходит контраст (см. §2.3 дизайн-системы и комментарий в UiCard.vue).
const ctaVariant = computed(() => (props.variant === "surface-ink" ? "outline-inverted" : "outline"));
</script>

<template>
  <UiCard :variant="variant" class="flex h-full flex-col">
    <template #media>
      <UiMediaPlaceholder aspect="4/3" />
    </template>
    <template #title>{{ title }}</template>

    <p class="mb-16 text-small" :class="variant === 'surface-white' ? 'text-ink-700' : 'opacity-80'">
      {{ lessonsCount ?? "—" }} занятий · {{ hours ?? "—" }} часов
    </p>
    <p v-if="description" class="mb-24 text-body">{{ description }}</p>
    <p v-if="price" class="mb-24 font-display text-h4">{{ price.toLocaleString("ru-RU") }} ₽</p>

    <UiButton :variant="ctaVariant" :to="to" class="mt-auto">{{ ctaLabel }}</UiButton>
  </UiCard>
</template>
