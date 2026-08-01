// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  modules: ["@nuxtjs/tailwindcss", "@nuxt/image"],

  // tokens.css (CSS custom properties), fonts.css (@font-face for the local
  // brand fonts in public/fonts), then the Tailwind directives.
  css: ["~/assets/styles/tokens.css", "~/assets/styles/fonts.css", "~/assets/css/main.css"],

  // Explicit per-directory entries: each is scanned relative to itself, so
  // components register under their bare filename (<UiButton>, <Hero>,
  // <CourseCard>, <AppHeader>, ...) with no directory-name prefix anywhere.
  components: [
    { path: "~/components/ui", pathPrefix: false },
    { path: "~/components/sections", pathPrefix: false },
    { path: "~/components/layout", pathPrefix: false },
  ],

  // Публичные контентные страницы рендерятся на сервере (SSR) для SEO.
  // Админка — чистый client-side рендер, ей SSR не нужен.
  routeRules: {
    "/admin/**": { ssr: false },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8080",
      // Бэкенд пока не отдаёт публичные эндпоинты по slug (курс/мастер-класс
      // с деталями) — см. TODO в app/composables/useApi.ts. Пока это так,
      // useApi() отдаёт данные из app/mocks/*, чтобы вёрстка не блокировалась.
      useMocks: (process.env.NUXT_PUBLIC_USE_MOCKS ?? "true") === "true",
    },
  },

  // Прокси для dev-режима: запросы с localhost:3000/api/** уходят на Go-сервер
  // без ручной настройки CORS у каждого разработчика (бэкенд по умолчанию
  // доступен напрямую благодаря CORS-мидлваре, но прокси удобнее локально).
  nitro: {
    preset: "bun",
    devProxy: {
      "/api": {
        target: "http://localhost:8080/api",
        changeOrigin: true,
      },
    },
  },

  image: {
    format: ["webp"],
  },

  typescript: {
    strict: true,
  },
});
