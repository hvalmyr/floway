<script setup lang="ts">
import { Gift } from "lucide-vue-next";

/**
 * Full-width scrolling promo strip for gift certificates, sat between the
 * header and page content on every page (see layouts/default.vue) — same
 * spot conceptually as a "free shipping" ticker.
 *
 * The whole strip is one link (single accessible name via aria-label); the
 * scrolling text itself is aria-hidden since it's the same sentence
 * repeated many times purely for the continuous-scroll effect, which would
 * otherwise be read out loud once per repeat.
 *
 * Seamless loop: REPEAT_COUNT copies of one unit (icon + text + dot) sit in
 * a `w-max` row, animated from translateX(0) to translateX(-50%). Because
 * every copy is identical, sliding by exactly half the track's width lands
 * on a pixel-identical frame to the start — this only self-aligns when
 * REPEAT_COUNT is even, so it must stay even if ever changed.
 */
const REPEAT_COUNT = 8;

const { text } = await usePageContent();
const message = computed(() =>
  text("gift_certificate_marquee_text", "Подарочные сертификаты — порадуйте близких цветами"),
);
</script>

<template>
  <NuxtLink
    to="/sertifikaty"
    :aria-label="message"
    class="group relative block w-full cursor-pointer overflow-hidden bg-primary py-12 transition-colors duration-300 hover:bg-primary/85"
  >
    <div
      aria-hidden="true"
      class="flex w-max animate-[gift-marquee_16s_linear_infinite] items-center gap-48 whitespace-nowrap motion-reduce:animate-none group-hover:[animation-duration:32s]"
    >
      <span v-for="n in REPEAT_COUNT" :key="n" class="flex shrink-0 items-center gap-12">
        <Gift class="size-24 shrink-0 text-white" aria-hidden="true" />
        <span class="font-display text-h4 text-white">{{ message }}</span>
        <span class="font-display text-h4 text-white/50">•</span>
      </span>
    </div>
  </NuxtLink>
</template>

<style scoped>
@keyframes gift-marquee {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-50%);
  }
}
</style>
