<script setup lang="ts">
/**
 * Centered H2 + optional lead paragraph, used to open every landing section.
 *
 * @example
 * <SectionHeading>
 *   Почему стоит учиться в школе «ФлоВей»
 *   <template #lead>Никакой воды — только то, что действительно важно перед стартом.</template>
 * </SectionHeading>
 */
</script>

<template>
  <div class="mx-auto flex max-w-[720px] flex-col items-center gap-16 text-center">
    <h2 class="font-display text-h2 text-ink-900">
      <slot />
    </h2>
    <p v-if="$slots.lead" class="text-body-l text-ink-700">
      <slot name="lead" />
    </p>
  </div>
</template>
