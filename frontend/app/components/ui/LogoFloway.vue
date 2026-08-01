<script setup lang="ts">
 /**
 * "фловей" wordmark as an SVG (not plain text), so it scales crisply at any
 * size. `textLength`/`lengthAdjust` pin the rendered width explicitly, so it
 * never clips or leaves stray whitespace regardless of how the browser
 * measures the custom font. Color comes from `currentColor` — set text-*
 * on the wrapping element/class. Size via height utility classes (width
 * follows automatically from the 168:48 viewBox aspect ratio).
 *
 * @example
 * <NuxtLink to="/"><LogoFloway class="h-32 w-auto text-primary" /></NuxtLink>
 * <NuxtLink to="/"><LogoFloway class="h-48 w-auto text-white" /></NuxtLink>
 */
</script>

<template>
  <svg viewBox="0 0 168 48" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="ФлоВей">
    <text
      x="0"
      y="38"
      font-family="'Soyuz Grotesk', sans-serif"
      font-weight="700"
      font-size="40"
      fill="currentColor"
      textLength="168"
      lengthAdjust="spacingAndGlyphs"
    >
      фловей
    </text>
  </svg>
</template>
