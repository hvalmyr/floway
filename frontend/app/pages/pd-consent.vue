<script setup lang="ts">
useSeoMeta({ title: "Согласие на обработку персональных данных — Фловей" });

const { text } = await usePageContent();
</script>

<template>
  <div class="container flex flex-col gap-24 py-48 sm:py-64">
    <h1 class="font-display text-h1 text-ink">Согласие на обработку персональных данных</h1>
    <MarkdownContent
      :source="text('legal_pd_consent', 'Текст согласия находится на согласовании у заказчика.')"
    />
  </div>
</template>
