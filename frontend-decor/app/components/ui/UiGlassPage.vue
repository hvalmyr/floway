<script setup lang="ts">
 /**
 * Full-page glass card — white, blurred background so page content stays
 * legible over the animated tree background (AmbientTreeBackground.vue).
 * For mostly-text pages (legal documents, the blog, the thank-you page)
 * instead of a landing page built from alternating full-bleed sections;
 * those keep their own per-section bg-surface/55 glass panels instead.
 *
 * Doesn't impose a gap/flex layout on its content — nest your own
 * `flex flex-col gap-*` wrapper inside so each page picks its own spacing
 * (a document page's flat gap-24 vs. the blog listing's gap-48 before its
 * card grid).
 *
 * @example
 * <UiGlassPage>
 *   <div class="flex flex-col gap-24">
 *     <h1 class="font-display text-h1 text-ink">Заголовок</h1>
 *     <MarkdownContent :source="text('legal_terms', '...')" />
 *   </div>
 * </UiGlassPage>
 */
</script>

<template>
  <div class="container py-48 sm:py-64 lg:py-80">
    <div
      class="rounded-[30px] bg-white/55 p-24 backdrop-blur backdrop-saturate-150 sm:p-48 lg:p-64"
    >
      <slot />
    </div>
  </div>
</template>
