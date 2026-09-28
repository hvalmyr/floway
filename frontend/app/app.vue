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
  // Preloads only the two font files the hero/header need for first paint
  // (headings + body text, both above the fold) — deliberately NOT every
  // font the page's DOM matches somewhere. That was tried (preloading
  // NonBureau-Light/Medium/Bold too, for badges/form labels further down
  // the page) to close a ~1.1s font critical-path chain PageSpeed flagged;
  // confirmed live it backfired — those extra fonts are fetched at the
  // browser's "Highest" priority (same as these two), competing with the
  // hero image's own fetch (only "High" priority even with
  // fetchpriority="high") for the same bandwidth, and desktop LCP measurably
  // got WORSE (resource load delay alone exceeded the *previous total* LCP).
  // font-display: swap (see fonts.css) already means below-the-fold text
  // never blocks on its font either way — worst case with only these two
  // preloaded is a brief fallback-face flash on text nobody's scrolled to
  // yet, which is strictly cheaper than delaying the LCP image.
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
