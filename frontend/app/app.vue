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
  // The two fonts the hero/header need for first paint (headings + body
  // text, both above the fold) preload at the browser's default priority
  // for `as="font"` (effectively "Highest") — they're on the critical path,
  // they should win any contention.
  //
  // NonBureau-Light/Medium/Bold (badges, form labels — all further down the
  // page) are a different case: PageSpeed flags them as a ~1.1s critical-
  // path chain because nothing preloads them, but a first attempt at fixing
  // that (preloading them at the same default "Highest" priority as the two
  // above) backfired — confirmed live, it competed with the hero image's
  // own fetch for bandwidth and made desktop LCP measurably worse. The
  // actual problem was never "not preloaded", it was "preloaded too
  // aggressively" — `fetchpriority="low"` keeps the early-discovery win
  // (the browser's preload scanner starts the request immediately instead
  // of waiting to reach these fonts' rules during CSSOM matching) while
  // scheduling it behind anything more urgent, so it no longer competes
  // with the hero image. font-display: swap (fonts.css) already means none
  // of this ever blocks render either way.
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
    {
      rel: "preload",
      as: "font",
      type: "font/woff2",
      href: "/fonts/NonBureau-Light.woff2",
      crossorigin: "anonymous",
      fetchpriority: "low",
    },
    {
      rel: "preload",
      as: "font",
      type: "font/woff2",
      href: "/fonts/NonBureau-Medium.woff2",
      crossorigin: "anonymous",
      fetchpriority: "low",
    },
    {
      rel: "preload",
      as: "font",
      type: "font/woff2",
      href: "/fonts/NonBureau-Bold.woff2",
      crossorigin: "anonymous",
      fetchpriority: "low",
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
