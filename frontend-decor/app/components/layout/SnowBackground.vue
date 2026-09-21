<script setup lang="ts">
/**
 * Decorative falling-snow background (п. 5.2 ТЗ) — decor site's own
 * equivalent of the school's AmbientTreeBackground.vue, built on the same
 * architectural pattern (fixed, transparent, `pointer-events-none`, sits
 * behind page content) but much simpler: no WebGL/Three.js, no interaction —
 * just CSS-animated dots. Mounted client-only via the Lazy+v-if pattern in
 * layouts/default.vue, so randomized per-flake styles never have to match
 * between an SSR pass and hydration.
 *
 * prefers-reduced-motion needs no handling here — the global rule in
 * assets/css/main.css already forces every CSS animation to a single,
 * near-instant frame, which lands each flake at its own end-of-fall
 * position (past the bottom edge, invisible) instead of animating.
 */
const FLAKE_COUNT = 28;

interface Flake {
  left: string;
  size: string;
  duration: string;
  delay: string;
  drift: string;
  opacity: string;
}

function randomBetween(min: number, max: number): number {
  return min + Math.random() * (max - min);
}

const flakes: Flake[] = Array.from({ length: FLAKE_COUNT }, () => ({
  left: `${randomBetween(0, 100)}vw`,
  size: `${randomBetween(3, 9).toFixed(1)}px`,
  duration: `${randomBetween(9, 22).toFixed(1)}s`,
  delay: `-${randomBetween(0, 20).toFixed(1)}s`,
  drift: `${randomBetween(-40, 40).toFixed(0)}px`,
  opacity: randomBetween(0.35, 0.85).toFixed(2),
}));
</script>

<template>
  <div class="pointer-events-none fixed inset-0 -z-10 overflow-hidden" aria-hidden="true">
    <span
      v-for="(flake, i) in flakes"
      :key="i"
      class="snowflake absolute top-[-5vh] block rounded-full bg-white"
      :style="{
        left: flake.left,
        width: flake.size,
        height: flake.size,
        opacity: flake.opacity,
        animationDuration: flake.duration,
        animationDelay: flake.delay,
        '--drift': flake.drift,
      }"
    />
  </div>
</template>

<style scoped>
.snowflake {
  animation-name: snowfall;
  animation-timing-function: linear;
  animation-iteration-count: infinite;
}

@keyframes snowfall {
  from {
    transform: translate(0, 0);
  }
  to {
    transform: translate(var(--drift), 110vh);
  }
}
</style>
