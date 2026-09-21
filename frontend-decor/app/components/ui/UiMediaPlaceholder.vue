<script setup lang="ts">
import { Image } from "lucide-vue-next";

/**
 * Decorative stand-in for a photo/video that hasn't been supplied yet
 * (no real photography assets exist at build time — Figma access wasn't
 * available either, see project notes). aria-hidden because it carries no
 * real content.
 *
 * TODO: заменить на <NuxtImg format="webp" :sizes="..." fit="cover" /> с
 * содержательным alt, когда появятся реальные фото (см. `gallery`/`coverImage`
 * на Course/Masterclass).
 *
 * @example
 * <UiMediaPlaceholder aspect="4/3" />
 */
withDefaults(
  defineProps<{
    aspect?: "4/3" | "3/4" | "1/1" | "16/9";
  }>(),
  { aspect: "4/3" },
);
</script>

<template>
  <div
    aria-hidden="true"
    class="flex items-center justify-center rounded-lg bg-surface text-ink/40"
    :class="{
      'aspect-[4/3]': aspect === '4/3',
      'aspect-[3/4]': aspect === '3/4',
      'aspect-square': aspect === '1/1',
      'aspect-video': aspect === '16/9',
    }"
  >
    <Image class="size-40" :stroke-width="1.5" />
  </div>
</template>
