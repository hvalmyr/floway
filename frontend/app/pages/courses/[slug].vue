<script setup lang="ts">
// TODO: Hero, таблица блоков курса, учебный план (аккордеон), форма записи.
// Вёрстка ждёт данных из Figma (desktop-course, node-id 0:203).
const route = useRoute()
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 py-16">
    <h1 class="text-3xl font-semibold">Курс: {{ route.params.slug }}</h1>
    <p class="mt-2 text-[var(--color-text-muted)]">Страница курса — в разработке.</p>
  </div>
</template>
