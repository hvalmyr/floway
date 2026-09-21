<script setup lang="ts">
/**
 * Decorative fir branch (п. 5.2 ТЗ: "хвойная веточка" — для углов блоков,
 * разделителей и подвала). A curved stem with alternating needle clusters,
 * inherits `currentColor` like every other Icon*.vue. Designed to sit in a
 * corner — the stem runs from the bottom-left up and to the right.
 *
 * @example
 * <IconFirBranch class="h-64 w-auto text-primary" />
 */
const NEEDLE_CLUSTERS = [
  { x: 18, y: 46, angle: -35 },
  { x: 34, y: 36, angle: -30 },
  { x: 50, y: 27, angle: -25 },
  { x: 64, y: 19, angle: -20 },
  { x: 76, y: 12, angle: -15 },
];
</script>

<template>
  <svg viewBox="0 0 120 60" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path
      d="M6 54 C 30 50, 55 40, 96 8"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      fill="none"
    />
    <g
      v-for="cluster in NEEDLE_CLUSTERS"
      :key="cluster.x"
      :transform="`translate(${cluster.x} ${cluster.y}) rotate(${cluster.angle})`"
    >
      <path
        d="M0 0 L -9 -10 M0 0 L -9 2 M0 0 L -9 12"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
      />
      <path
        d="M0 0 L 9 -10 M0 0 L 9 2 M0 0 L 9 12"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
      />
    </g>
    <circle cx="96" cy="8" r="2.5" fill="currentColor" />
  </svg>
</template>
