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
  // Preloads every self-hosted font file the initial SSR HTML actually
  // matches somewhere in the DOM (not just the hero's own headings/body
  // text) so the browser starts fetching all of them immediately instead
  // of discovering each one only once it parses the @font-face rules
  // (fonts.css, now inlined — see nuxt.config.ts's features.inlineStyles).
  // Used to be just SoyuzGrotesk-Bold + NonBureau-Regular on the theory
  // that only the hero needed fonts fast; a real PageSpeed audit showed
  // the other three NonBureau weights (badges, form labels, ... —
  // rendered lower on the page but still part of the same SSR payload)
  // showing up as a ~1.1s critical-path chain instead, because the
  // browser doesn't need visibility to discover a matched @font-face, only
  // a DOM match, and the whole page's DOM is already there from SSR.
  // font-display: swap (see fonts.css) already avoids a render block
  // either way — this only shortens how long matched text shows in the
  // fallback face.
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
    },
    {
      rel: "preload",
      as: "font",
      type: "font/woff2",
      href: "/fonts/NonBureau-Medium.woff2",
      crossorigin: "anonymous",
    },
    {
      rel: "preload",
      as: "font",
      type: "font/woff2",
      href: "/fonts/NonBureau-Bold.woff2",
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
