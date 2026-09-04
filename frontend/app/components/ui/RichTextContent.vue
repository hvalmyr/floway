<script setup lang="ts">
/**
 * Renders HTML produced by AdminRichTextEditor.vue (blog post content).
 * The source is already sanitized against a fixed tag/attribute whitelist
 * at save time (see richTextSanitize.ts) — same trust boundary as
 * MarkdownContent.vue, which likewise trusts admin-authored source without
 * re-sanitizing on render.
 *
 * @example
 * <RichTextContent :source="post.content" />
 */
defineProps<{ source: string }>();
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div class="rich-text-content font-body text-body text-ink" v-html="source" />
</template>

<style scoped>
.rich-text-content :deep(p) {
  margin-bottom: 1em;
}
.rich-text-content :deep(p:last-child) {
  margin-bottom: 0;
}
.rich-text-content :deep(h3) {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-family: var(--font-display);
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
  border-radius: 0.5rem;
}
.rich-text-content :deep(strong) {
  font-weight: 700;
}
.rich-text-content :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
}
</style>
