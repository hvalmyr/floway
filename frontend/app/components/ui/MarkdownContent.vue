<script setup lang="ts">
import MarkdownIt from "markdown-it";

/**
 * Renders admin-authored markdown (page_content.value) as HTML, styled with
 * the site's existing type tokens. `html: false` escapes any raw HTML in the
 * source instead of passing it through — content comes from the admin panel,
 * not public users, but there's no reason to allow arbitrary HTML/script
 * injection there either.
 *
 * @example
 * <MarkdownContent :source="text('legal_privacy_policy')" />
 */
const props = defineProps<{ source: string }>();

const md = new MarkdownIt({ html: false, linkify: true, breaks: true });
const html = computed(() => md.render(props.source));
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div class="markdown-content font-body text-body text-ink" v-html="html" />
</template>

<style scoped>
.markdown-content :deep(p) {
  margin-bottom: 1em;
}
.markdown-content :deep(p:last-child) {
  margin-bottom: 0;
}
.markdown-content :deep(h2) {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-family: var(--font-display);
  font-size: var(--text-h4);
}
.markdown-content :deep(h3) {
  margin-top: 1.25em;
  margin-bottom: 0.5em;
  font-family: var(--font-display);
  font-size: var(--text-h4);
}
.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin-bottom: 1em;
  padding-left: 1.5em;
}
.markdown-content :deep(ul) {
  list-style: disc;
}
.markdown-content :deep(ol) {
  list-style: decimal;
}
.markdown-content :deep(li) {
  margin-bottom: 0.25em;
}
.markdown-content :deep(strong) {
  font-weight: 700;
}
.markdown-content :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
}
</style>
