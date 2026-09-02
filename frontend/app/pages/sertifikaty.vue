<script setup lang="ts">
useSeoMeta({
  title: "Подарочные сертификаты — Фловей",
  description:
    "Подарочный сертификат школы флористики «Фловей» на любую сумму, курс или мастер-класс.",
});

const { text } = await usePageContent();
</script>

<template>
  <div class="container flex flex-col gap-24 py-48 sm:py-64">
    <h1 class="font-display text-h1 text-ink">Подарочные сертификаты</h1>
    <div
      class="flex w-full flex-col gap-16 rounded-md bg-white/55 p-24 backdrop-blur backdrop-saturate-150"
    >
      <p class="font-body text-body text-ink">
        {{
          text(
            "gift_certificate_page_body",
            "Подарите близким флористику: сертификат можно оформить на любую сумму или на конкретный курс/мастер-класс школы «Фловей». Чтобы приобрести сертификат, свяжитесь с нами любым удобным способом — подберём вариант и оформим за вас.",
          )
        }}
      </p>
      <UiButton to="/contacts" variant="primary" class="self-start">Связаться с нами</UiButton>
    </div>
  </div>
</template>
