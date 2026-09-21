<script setup lang="ts">
/**
 * Renders HTML produced by AdminRichTextEditor.vue (blog post content).
 * The source is already sanitized against a fixed tag/attribute whitelist
 * at save time (see richTextSanitize.ts) — same trust boundary as
 * MarkdownContent.vue, which likewise trusts admin-authored source without
 * re-sanitizing on render.
 *
 * Clicking any inline image opens it in a fullscreen lightbox over a brown
 * (`bg-ink`) backdrop — same visual language as PhotoCarousel.vue's
 * lightbox, just without that component's thumbnail strip/prev-next
 * sequencing, since a single content image has no siblings to page
 * through.
 *
 * @example
 * <RichTextContent :source="post.content" />
 */
import { X } from "lucide-vue-next";
import { nextTick, onUnmounted, ref, watch } from "vue";

defineProps<{ source: string }>();

const lightboxSrc = ref<string | null>(null);
const lightboxAlt = ref("");
const dialogRef = ref<HTMLElement | null>(null);

function onContentClick(event: MouseEvent) {
  const target = event.target as HTMLElement;
  if (target.tagName !== "IMG") return;
  const img = target as HTMLImageElement;
  lightboxSrc.value = img.currentSrc || img.src;
  lightboxAlt.value = img.alt;
}

function closeLightbox() {
  lightboxSrc.value = null;
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === "Escape") closeLightbox();
}

// Focused on open so Escape works immediately without requiring a prior
// click inside the dialog — a plain <div> isn't focusable by default,
// hence the tabindex="-1" in the template.
watch(lightboxSrc, (src) => {
  document.body.style.overflow = src ? "hidden" : "";
  if (src) nextTick(() => dialogRef.value?.focus());
});
onUnmounted(() => {
  document.body.style.overflow = "";
});
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div
    class="rich-text-content font-body text-body text-ink"
    data-no-orbit
    v-html="source"
    @click="onContentClick"
  />

  <Teleport to="body">
    <div
      v-if="lightboxSrc"
      ref="dialogRef"
      class="fixed inset-0 z-50 flex items-center justify-center bg-ink/90 p-24 outline-none"
      role="dialog"
      aria-modal="true"
      aria-label="Просмотр изображения"
      tabindex="-1"
      data-no-orbit
      @click.self="closeLightbox"
      @keydown="onKeydown"
    >
      <button
        type="button"
        aria-label="Закрыть"
        class="absolute right-16 top-16 grid size-[44px] place-items-center rounded-full bg-white/10 text-white hover:bg-white/20"
        @click="closeLightbox"
      >
        <X class="size-24" aria-hidden="true" />
      </button>
      <img
        :src="lightboxSrc"
        :alt="lightboxAlt"
        class="h-[90vh] max-h-[90vh] w-auto max-w-[94vw] cursor-zoom-out rounded-lg object-contain"
        @click="closeLightbox"
      />
    </div>
  </Teleport>
</template>

<style scoped>
.rich-text-content :deep(p) {
  margin-bottom: 1em;
}
.rich-text-content :deep(p:last-child) {
  margin-bottom: 0;
}
.rich-text-content :deep(h2),
.rich-text-content :deep(h3) {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-family: var(--font-display);
}
.rich-text-content :deep(h2) {
  font-size: var(--text-h2);
}
.rich-text-content :deep(h3) {
  font-size: var(--text-h4);
}
.rich-text-content :deep(blockquote) {
  margin: 1.5em 0;
  padding-left: 1em;
  border-left: 3px solid var(--color-primary);
  font-style: italic;
  color: var(--color-text-muted);
}
.rich-text-content :deep(figure) {
  margin: 1.5em 0;
}
.rich-text-content :deep(img) {
  max-width: 100%;
  /* Caps tall/portrait photos to roughly the column's own width instead of
     letting them stretch edge-to-edge at full width — same trick
     telegra.ph's article CSS uses (max-height ≈ content column width) so a
     portrait shot reads as a framed photo instead of dominating the page.
     `display: block` + auto side margins (rather than the figure's
     text-align) center it — Tailwind's preflight already forces `img` to
     `display: block`, which text-align can't center. */
  max-height: 720px;
  border-radius: 0.5rem;
  display: block;
  margin-left: auto;
  margin-right: auto;
  cursor: zoom-in;
  transition: opacity 0.15s;
}
.rich-text-content :deep(img:hover) {
  opacity: 0.9;
}
.rich-text-content :deep(strong) {
  font-weight: 700;
}
.rich-text-content :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
}
</style>
