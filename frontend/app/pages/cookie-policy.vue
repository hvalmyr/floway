<script setup lang="ts">
useSeoMeta({ title: "Политика использования файлов cookie — Фловей" });

const { text } = await usePageContent();
</script>

<template>
  <UiGlassPage>
    <div class="flex flex-col gap-24">
      <h1 class="font-display text-h1 text-ink">Политика использования файлов cookie</h1>
      <MarkdownContent
        :source="
          text('legal_cookie_policy', 'Текст политики находится на согласовании у заказчика.')
        "
      />
    </div>
  </UiGlassPage>
</template>
