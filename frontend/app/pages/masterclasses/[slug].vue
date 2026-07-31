<script setup lang="ts">
// TODO: Hero, инфоблок (длительность/цена), описание, форма записи.
// Вёрстка ждёт данных из Figma (desktop-mc, node-id 0:306).
const route = useRoute();
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 py-16">
    <h1 class="text-3xl font-semibold">Мастер-класс: {{ route.params.slug }}</h1>
    <p class="mt-2 text-[var(--color-text-muted)]">Страница мастер-класса — в разработке.</p>
  </div>
</template>
