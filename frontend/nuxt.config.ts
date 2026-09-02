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
    // Админ-панель — отдельная зона (не дизайн-система), поэтому свой
    // каталог вместо смешивания с components/ui.
    { path: "~/components/admin", pathPrefix: false },
  ],

  // Публичные контентные страницы рендерятся на сервере (SSR) для SEO.
  // Админка — чистый client-side рендер, ей SSR не нужен.
  routeRules: {
    "/admin/**": { ssr: false },
  },

  runtimeConfig: {
    // Сервер (SSR) и браузер видят backend по-разному в Docker: браузер идёт
    // по публичному адресу (apiBase), а SSR-запросы из контейнера фронтенда
    // должны идти по внутренней docker-сети (например http://backend:8080),
    // иначе "localhost" внутри контейнера указывает сам на себя, а не на
    // backend. Локально без Docker (bun run dev) оба адреса совпадают, так
    // что apiBaseInternal по умолчанию просто берёт публичный.
    apiBaseInternal:
      process.env.NUXT_API_BASE_INTERNAL ||
      process.env.NUXT_PUBLIC_API_BASE ||
      "http://localhost:8080",
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8080",
      // Бэкенд теперь отдаёт публичные эндпоинты по slug (GET
      // /api/v1/courses/{slug}/full, GET /api/v1/masterclasses/{slug}) —
      // моки из app/mocks/* используются только если явно включить флагом,
      // например для вёрстки без поднятого бэкенда.
      useMocks: process.env.NUXT_PUBLIC_USE_MOCKS === "true",
      // Аналитика и подтверждение прав на сайт — все три необязательны и по
      // умолчанию пустые (счётчик/мета-теги просто не рендерятся в dev).
      // Читаются из окружения при старте контейнера (не на этапе сборки),
      // как и остальные NUXT_PUBLIC_* — см. .env.example.
      yandexMetrikaId: process.env.NUXT_PUBLIC_YANDEX_METRIKA_ID || "",
      yandexVerification: process.env.NUXT_PUBLIC_YANDEX_VERIFICATION || "",
      googleSiteVerification: process.env.NUXT_PUBLIC_GOOGLE_SITE_VERIFICATION || "",
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
