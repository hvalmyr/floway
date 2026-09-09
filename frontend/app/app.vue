<script setup lang="ts">
// Права на сайт подтверждаются мета-тегами (Яндекс.Вебмастер, Google Search
// Console) — оба необязательны и просто не попадают в <head>, если
// соответствующая переменная окружения не задана.
const { public: publicConfig } = useRuntimeConfig();

useHead({
  htmlAttrs: {
    lang: "ru",
  },
  // Preloads the two font files every page's first paint actually needs
  // (headings + body text, both above the fold in the header/hero) so the
  // browser starts fetching them immediately instead of waiting to
  // discover the @font-face rules in fonts.css after CSSOM is built.
  // font-display: swap (see fonts.css) already avoids a render block, so
  // this only shortens how long text shows in the fallback face — not a
  // correctness fix, just less visible font-swap.
  link: [
    {
      rel: "preload",
      as: "font",
      type: "font/woff",
      href: "/fonts/SoyuzGrotesk-Bold.woff",
      crossorigin: "anonymous",
    },
    {
      rel: "preload",
      as: "font",
      type: "font/woff2",
      href: "/fonts/NonBureau-Regular.woff2",
      crossorigin: "anonymous",
    },
  ],
  meta: [
    ...(publicConfig.yandexVerification
      ? [{ name: "yandex-verification", content: publicConfig.yandexVerification }]
      : []),
    ...(publicConfig.googleSiteVerification
      ? [{ name: "google-site-verification", content: publicConfig.googleSiteVerification }]
      : []),
  ],
});
</script>

<template>
  <NuxtLayout>
    <NuxtRouteAnnouncer />
    <NuxtPage />
  </NuxtLayout>
</template>
