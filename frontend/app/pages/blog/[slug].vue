<script setup lang="ts">
// TODO: обложка, заголовок, метаданные, текст статьи.
const route = useRoute();
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 py-16">
    <h1 class="text-3xl font-semibold">Статья: {{ route.params.slug }}</h1>
    <p class="mt-2 text-[var(--color-text-muted)]">Страница статьи — в разработке.</p>
  </div>
</template>
