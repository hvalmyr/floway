<script setup lang="ts">
useSeoMeta({ title: "Пользовательское соглашение — Фловей" });

const { text } = await usePageContent();
</script>

<template>
  <UiGlassPage>
    <div class="flex flex-col gap-24">
      <h1 class="font-display text-h1 text-ink">Пользовательское соглашение</h1>
      <MarkdownContent
        :source="text('legal_terms', 'Текст соглашения находится на согласовании у заказчика.')"
      />
    </div>
  </UiGlassPage>
</template>
