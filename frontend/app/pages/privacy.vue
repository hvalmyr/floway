<script setup lang="ts">
useSeoMeta({ title: "Политика обработки персональных данных — Фловей" });

const { text } = await usePageContent();
</script>

<template>
  <div class="container flex flex-col gap-24 py-48 sm:py-64">
    <h1 class="font-display text-h1 text-ink">Политика обработки персональных данных</h1>
    <MarkdownContent
      :source="
        text(
          'legal_privacy_policy',
          'Текст политики находится на согласовании у заказчика. Здесь будет размещён финальный документ о том, какие данные собирает школа «Фловей» через формы на сайте, как они хранятся и используются.',
        )
      "
    />
  </div>
</template>
