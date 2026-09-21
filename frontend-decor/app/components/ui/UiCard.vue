<script setup lang="ts">
/**
 * Base card shell (radius-md, 32px padding) with three built-in fill
 * variants, per the course/profile-courses mockups. `surface-primary` uses
 * WHITE text on the blue fill, matching the mockup exactly — contrast is
 * ~2.3:1 there (below WCAG AA), a deliberate tradeoff since the palette is
 * locked to exactly 4 colors with no darker "on-primary" shade available.
 *
 * `variant="custom"` applies no color classes at all — the caller sets its
 * own background/text color classes via the `class` attribute (fallthrough
 * onto this same root element), for callers with their own bg/text combos
 * that don't fit the 3 built-in variants (see CourseCard.vue's 4 display
 * styles).
 *
 * The default-slot wrapper is itself `flex flex-1 flex-col` (not a plain
 * block) so a caller stacking the root in an outer `flex flex-col` (again,
 * CourseCard) can push something inside the slot — its cover image — to the
 * bottom via `mt-auto`: that only resolves to a real value inside an actual
 * flex formatting context, and `flex-1` is what lets this wrapper actually
 * grow to fill the stretched card's leftover height instead of only its own
 * natural content height.
 *
 * @example
 * <UiCard variant="surface-primary">
 *   <template #title>Блок “Букеты”</template>
 *   7 занятий, 30 часов
 * </UiCard>
 * <UiCard variant="custom" class="bg-surface text-primary">...</UiCard>
 */
withDefaults(
  defineProps<{
    variant?: "surface-white" | "surface-primary" | "surface-ink" | "custom";
  }>(),
  { variant: "surface-white" },
);
</script>

<template>
  <div
    class="rounded-md p-32"
    :class="[
      variant === 'surface-white' && 'border-2 border-primary bg-white text-primary',
      variant === 'surface-primary' && 'bg-primary text-white',
      variant === 'surface-ink' && 'bg-ink text-white',
    ]"
  >
    <div v-if="$slots.media" class="mb-24">
      <slot name="media" />
    </div>
    <h3 v-if="$slots.title" class="mb-12 font-display text-h4">
      <slot name="title" />
    </h3>
    <div
      class="flex flex-1 flex-col font-body text-body"
      :class="variant === 'surface-white' ? 'text-ink' : ''"
    >
      <slot />
    </div>
  </div>
</template>
