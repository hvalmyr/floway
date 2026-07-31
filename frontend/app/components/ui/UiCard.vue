<script setup lang="ts">
/**
 * Base card shell (radius-md, 32px padding) with three fill variants.
 *
 * Note on `surface-primary` text/CTA color: docs/floway-design.md §5.2 says
 * "узнать подробнее" buttons use white text on both dark *and* colored
 * cards — but §2.3 already documents that white text on `primary-500`
 * (~2.3:1) fails contrast even for large bold text. We resolve that
 * inconsistency the same way §2.3 resolves it for the hero button: dark
 * `ink-900` text/border on the light-blue `surface-primary` card, white only
 * on the genuinely dark `surface-ink` card.
 *
 * @example
 * <UiCard variant="surface-primary">
 *   <template #title>Блок «Букеты»</template>
 *   7 занятий · 30 часов
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
      variant === 'surface-white' && 'border border-line bg-white text-ink-900',
      variant === 'surface-primary' && 'bg-primary-500 text-ink-900',
      variant === 'surface-ink' && 'bg-ink-900 text-white',
    ]"
  >
    <div v-if="$slots.media" class="mb-24">
      <slot name="media" />
    </div>
    <h3 v-if="$slots.title" class="mb-12 font-display text-h3">
      <slot name="title" />
    </h3>
    <div class="text-body">
      <slot />
    </div>
  </div>
</template>
