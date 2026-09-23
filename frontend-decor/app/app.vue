<script setup lang="ts">
// Права на сайт подтверждаются мета-тегами (Яндекс.Вебмастер, Google Search
// Console) — оба необязательны и просто не попадают в <head>, если
// соответствующая переменная окружения не задана.
const { public: publicConfig } = useRuntimeConfig();

// useNonce() only resolves during SSR (it reads the per-request value off
// the SSR event context — see nuxt-security's composables/nonce.ts), so it
// has to be read here and stamped into a meta tag rather than called from
// client-only code. useYandexMetrika.ts reads it back off this tag when it
// injects its own inline <script> client-side.
const nonce = useNonce();

// getRequestURL(event).origin is what sitemap.xml.ts uses server-side for
// the same reason (the request's own Host header, not config.public.apiBase
// — the two only coincide in prod); useRequestURL() is its universal
// (SSR+client) composable equivalent. route.path excludes the query string,
// which is exactly what makes this useful: ad click-tracking params
// (Yandex Direct's ?etext=..., ?utm_*) point back at the plain page instead
// of search engines treating the tagged URL as a separate, duplicate one.
const route = useRoute();
const requestUrl = useRequestURL();
const canonicalHref = computed(() => `${requestUrl.origin}${route.path}`);

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
    { rel: "canonical", href: canonicalHref },
  ],
  meta: [
    ...(publicConfig.yandexVerification
      ? [{ name: "yandex-verification", content: publicConfig.yandexVerification }]
      : []),
    ...(publicConfig.googleSiteVerification
      ? [{ name: "google-site-verification", content: publicConfig.googleSiteVerification }]
      : []),
    ...(nonce ? [{ name: "csp-nonce", content: nonce }] : []),
  ],
});
</script>

<template>
  <NuxtLayout>
    <NuxtRouteAnnouncer />
    <NuxtPage />
  </NuxtLayout>
</template>
