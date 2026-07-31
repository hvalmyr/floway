// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss'],

  css: ['~/assets/css/main.css'],

  // Публичные контентные страницы рендерятся на сервере (SSR) для SEO.
  // Админка — чистый client-side рендер, ей SSR не нужен.
  routeRules: {
    '/admin/**': { ssr: false }
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
    }
  },

  typescript: {
    strict: true
  }
})
