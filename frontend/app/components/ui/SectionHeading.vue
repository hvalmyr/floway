<script setup lang="ts">
/**
 * Centered H2 + optional lead paragraph, used to open every landing section.
 * Heading color alternates blue/dark per section in the mockups (no single
 * rule — set explicitly per section via `color`).
 *
 * The lead sits in a glass container so it stays legible over the animated
 * tree background. When the section itself already has that glass fill
 * (`bg-surface/55` beige sections), pass `on-glass` to keep the container's
 * shape/spacing but drop its own fill — otherwise glass-on-glass double-tints.
 *
 * @example
 * <SectionHeading color="primary">
 *   Почему стоит учиться в школе “ФлоВей”
 *   <template #lead>Никакой воды — только то, что действительно важно перед стартом.</template>
 * </SectionHeading>
 */
withDefaults(defineProps<{ color?: "primary" | "ink"; onGlass?: boolean }>(), {
  color: "ink",
  onGlass: false,
});
</script>

<template>
  <div class="mx-auto flex w-full flex-col items-center gap-16 text-center">
    <h2 class="font-display text-h2" :class="color === 'primary' ? 'text-primary' : 'text-ink'">
      <slot />
    </h2>
    <p
      v-if="$slots.lead"
      class="w-full whitespace-pre-line rounded-md px-24 py-16 font-body text-body text-ink lg:w-4/5"
      :class="onGlass ? '' : 'bg-white/55 backdrop-blur backdrop-saturate-150'"
    >
      <slot name="lead" />
    </p>
  </div>
</template>
