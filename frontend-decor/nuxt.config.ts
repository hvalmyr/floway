// @nuxt/image's IPX provider fetches originals itself (server-side, either
// during SSR or when the browser later requests the /_ipx/** URL it
// generated) — the public apiBase domain either isn't known yet at build
// time (prod deploys a prebuilt image, the real domain only exists at
// deploy time) or isn't reachable from inside the frontend container (dev
// Docker). apiBaseInternal already solves exactly this for API calls (see
// useApiClient.ts) — reused here so image optimization gets a base that's
// always server-reachable, in every environment, without per-deploy config.
const mediaOptimizeBase =
  process.env.NUXT_API_BASE_INTERNAL || process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8080";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  modules: ["@nuxtjs/tailwindcss", "@nuxt/image", "nuxt-security"],

  // CSP only — every other header nuxt-security would set by default
  // (HSTS, X-Frame-Options, COOP, CORP, COEP, Permissions-Policy, ...) stays
  // owned by Caddy (see ansible/roles/deploy_app/templates/Caddyfile.j2),
  // set once for both this app and frontend (the school site). CSP moved
  // here instead because a real per-request nonce (needed to drop
  // 'unsafe-inline' from script-src) can only be generated and stamped onto
  // the SSR HTML by the app itself — Caddy has no way to do that. The
  // module's other features (CORS, rate limiting, request-size limits, XSS
  // body validation) are disabled below: they're not part of this change
  // and each has its own failure mode that could break existing behavior
  // (e.g. the request size limiter rejecting a legitimate admin image
  // upload) without being asked for.
  security: {
    nonce: true,
    headers: {
      contentSecurityPolicy: {
        "default-src": ["'self'"],
        // 'nonce-{{nonce}}' is the module's own per-request placeholder,
        // swapped for the real value by its Nitro plugins (see
        // node_modules/nuxt-security/dist/runtime/nitro/plugins/
        // 40-cspSsrNonce.js) — the mc.yandex.ru inline loader gets the same
        // nonce read client-side via useYandexMetrika.ts. Metrika's own
        // runtime code calls out to mc.yandex.com too (confirmed live in
        // local testing against frontend/'s identical setup, not
        // documented anywhere obvious) — both needed or webvisor/hit
        // tracking silently breaks.
        "script-src": ["'self'", "'nonce-{{nonce}}'", "https://mc.yandex.ru", "https://mc.yandex.com"],
        "script-src-attr": ["'none'"],
        // Vue's :style bindings compile to inline style="..." attributes —
        // style-src-attr has no nonce mechanism for those, so this stays
        // 'unsafe-inline' (Lighthouse's CSP/XSS audit only penalizes
        // script-src, not style-src).
        "style-src": ["'self'", "'unsafe-inline'"],
        "img-src": ["'self'", "https:", "data:"],
        "font-src": ["'self'"],
        // wss://mc.yandex.* is webvisor's session-replay socket (init sets
        // webvisor: true in useYandexMetrika.ts) — also only found by
        // actually loading the page under this CSP, not documented.
        "connect-src": [
          "'self'",
          "https://mc.yandex.ru",
          "https://mc.yandex.com",
          "wss://mc.yandex.ru",
          "wss://mc.yandex.com",
        ],
        // Unlike frontend/ (the school site), this stack embeds no iframes
        // at all — no Maps widget, no CMS map field — so nothing gets
        // allow-listed here beyond Metrika's own webvisor-related framing.
        "frame-src": ["'self'", "https://mc.yandex.com"],
        "frame-ancestors": ["'none'"],
        "base-uri": ["'self'"],
        "form-action": ["'self'"],
        "object-src": ["'none'"],
        // "vue" is Vue's own built-in passthrough policy (unconditionally
        // registered by @vue/runtime-dom whenever window.trustedTypes
        // exists) — v-html sinks sanitize their data before binding
        // (AppIcon.vue, RichTextContent.vue) rather than through a custom
        // policy, since Vue doesn't let a binding opt into a different one.
        // "default" is trusted-types.client.ts's catch-all for imperative
        // innerHTML writes (AdminRichTextEditor.vue).
        "trusted-types": ["vue", "default"],
        "require-trusted-types-for": ["'script'"],
      },
      strictTransportSecurity: false,
      xFrameOptions: false,
      crossOriginOpenerPolicy: false,
      crossOriginEmbedderPolicy: false,
      crossOriginResourcePolicy: false,
      xContentTypeOptions: false,
      referrerPolicy: false,
      originAgentCluster: false,
      xDNSPrefetchControl: false,
      xDownloadOptions: false,
      xPermittedCrossDomainPolicies: false,
      xXSSProtection: false,
      permissionsPolicy: false,
    },
    requestSizeLimiter: false,
    rateLimiter: false,
    xssValidator: false,
    corsHandler: false,
    allowedMethodsRestricter: false,
    sri: false,
    csrf: false,
    hidePoweredBy: true,
  },

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
    // Self-hosted fonts and sounds under public/ have no content hash in
    // their filename (unlike _nuxt/** build assets), so Nitro doesn't cache
    // them long by default — a real PageSpeed audit flagged ~417 KiB of
    // repeat-visit re-downloads because of this. They're static brand
    // assets that essentially never change; if one ever needs to, rename
    // the file (matching fonts.css/the sound's src) rather than relying on
    // this cache expiring.
    "/fonts/**": { headers: { "cache-control": "public, max-age=31536000, immutable" } },
    "/sounds/**": { headers: { "cache-control": "public, max-age=31536000, immutable" } },
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
      // Отдельная от apiBase база — специально для картинок, идущих через
      // @nuxt/image (см. комментарий у mediaOptimizeBase выше). В отличие
      // от apiBaseInternal она обязана быть public: NuxtImg/NuxtPicture
      // рендерят src и на сервере, и (при повторном рендере на клиенте)
      // в браузере, а результат должен совпадать один в один — иначе Vue
      // на каждой картинке ловит hydration mismatch. apiBase здесь не
      // подходит по той же причине, что и для API-запросов.
      mediaOptimizeBase: mediaOptimizeBase,
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
    // Without an explicit mount, useStorage("cache:ipx-manual") (see
    // ipx-cache.ts) falls back to unstorage's default in-memory driver —
    // fine for dev, but in prod it means the manually-cached IPX output
    // (avoiding the ~4s-per-image avif re-encode, see that file's own
    // comment) lives only in the container's RAM: wiped on every restart,
    // not just deploys, and growing unbounded for as long as the process
    // stays up. Pinning it to the filesystem instead makes it survive
    // restarts as long as ./.data/cache is mounted as a persistent volume
    // (see docker-compose.prod.yml) — relative path so it resolves under
    // whatever the process's cwd is in each environment (frontend/ in dev,
    // /app in the Docker image, per its WORKDIR).
    storage: {
      cache: { driver: "fs", base: "./.data/cache" },
    },
  },

  image: {
    format: ["webp"],
    quality: 80,
    // Allowlist for IPX's remote fetch (SSRF guard) — the host(s)
    // mediaOptimizeBase can actually resolve to across environments (see
    // its comment above). Static and stable across every deploy of this
    // repo's docker-compose setup, unlike the public domain.
    domains: [new URL(mediaOptimizeBase).host, "backend:8080", "localhost:8080"],
  },

  typescript: {
    strict: true,
  },
});
