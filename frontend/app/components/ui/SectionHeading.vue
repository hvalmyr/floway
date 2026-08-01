<script setup lang="ts">
/**
 * Centered H2 + optional lead paragraph, used to open every landing section.
 * Heading color alternates blue/dark per section in the mockups (no single
 * rule — set explicitly per section via `color`).
 *
 * @example
 * <SectionHeading color="primary">
 *   Почему стоит учиться в школе «ФлоВей»
 *   <template #lead>Никакой воды — только то, что действительно важно перед стартом.</template>
 * </SectionHeading>
 */
withDefaults(defineProps<{ color?: "primary" | "ink" }>(), { color: "ink" });
</script>

<template>
  <div class="mx-auto flex max-w-[720px] flex-col items-center gap-16 text-center">
    <h2 class="font-display text-h2" :class="color === 'primary' ? 'text-primary' : 'text-ink'">
      <slot />
    </h2>
    <p v-if="$slots.lead" class="font-body text-body-l text-ink">
      <slot name="lead" />
    </p>
  </div>
</template>
