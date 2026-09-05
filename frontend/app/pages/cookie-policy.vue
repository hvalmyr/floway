<script setup lang="ts">
useSeoMeta({ title: "Политика использования файлов cookie — Фловей" });

const { text } = await usePageContent();
</script>

<template>
  <div class="container flex flex-col gap-24 py-48 sm:py-64">
    <h1 class="font-display text-h1 text-ink">Политика использования файлов cookie</h1>
    <MarkdownContent
      :source="text('legal_cookie_policy', 'Текст политики находится на согласовании у заказчика.')"
    />
  </div>
</template>
