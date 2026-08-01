<script setup lang="ts">
/**
 * Base card shell (radius-md, 32px padding) with three fill variants, per
 * the course/profile-courses mockups. `surface-primary` uses WHITE text on
 * the blue fill, matching the mockup exactly — contrast is ~2.3:1 there
 * (below WCAG AA), a deliberate tradeoff since the palette is locked to
 * exactly 4 colors with no darker "on-primary" shade available.
 *
 * @example
 * <UiCard variant="surface-primary">
 *   <template #title>Блок «Букеты»</template>
 *   7 занятий, 30 часов
 * </UiCard>
 */
withDefaults(
  defineProps<{
    variant?: "surface-white" | "surface-primary" | "surface-ink";
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
    <h3 v-if="$slots.title" class="mb-12 font-display text-h3">
      <slot name="title" />
    </h3>
    <div class="font-body text-body" :class="variant === 'surface-white' ? 'text-ink' : ''">
      <slot />
    </div>
  </div>
</template>
