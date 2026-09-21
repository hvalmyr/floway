<script setup lang="ts">
/**
 * Six-fold symmetric decorative snowflake (п. 5.2 ТЗ: "сложная, детальная
 * снежинка" — для акцентов и фонов). Pure stroke art (no fill) so it reads
 * as a linear ornament rather than a solid icon; inherits `currentColor`,
 * same convention as every other Icon*.vue.
 *
 * @example
 * <IconSnowflake class="size-32 text-primary" />
 */
const ARM_ROTATIONS = [0, 60, 120, 180, 240, 300];
</script>

<template>
  <svg viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
    <g
      v-for="rotation in ARM_ROTATIONS"
      :key="rotation"
      :transform="`rotate(${rotation} 50 50)`"
      stroke="currentColor"
      stroke-linecap="round"
    >
      <line x1="50" y1="50" x2="50" y2="6" stroke-width="2.5" />
      <line x1="50" y1="26" x2="41" y2="17" stroke-width="1.75" />
      <line x1="50" y1="26" x2="59" y2="17" stroke-width="1.75" />
      <line x1="50" y1="16" x2="44" y2="10" stroke-width="1.5" />
      <line x1="50" y1="16" x2="56" y2="10" stroke-width="1.5" />
      <circle cx="50" cy="6" r="2.5" fill="currentColor" stroke="none" />
    </g>
    <circle cx="50" cy="50" r="4" fill="currentColor" stroke="none" />
  </svg>
</template>
