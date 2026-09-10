<script setup lang="ts">
useSeoMeta({ title: "Согласие на обработку персональных данных — Фловей" });

const { text } = await usePageContent();
</script>

<template>
  <UiGlassPage>
    <div class="flex flex-col gap-24">
      <h1 class="font-display text-h1 text-ink">Согласие на обработку персональных данных</h1>
      <MarkdownContent
        :source="text('legal_pd_consent', 'Текст согласия находится на согласовании у заказчика.')"
      />
    </div>
  </UiGlassPage>
</template>
