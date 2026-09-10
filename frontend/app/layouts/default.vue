<script setup lang="ts">
import { onMounted, ref } from "vue";

// The 3D background is purely decorative (fixed, -z-10, pointer-events-none
// — see AmbientTreeBackground.vue) and safe to start a moment late.
// Deferred until the `load` event — i.e. until everything the page
// actually requested up front (hero image, course cards, ...) has
// finished — rather than starting immediately alongside them and
// competing for the same network/CPU budget.
//
// An earlier attempt deferred this via requestIdleCallback with a short
// timeout instead, and made things measurably worse (Time to Interactive
// went from 13.5s to 28.1s on a real audit) — Lighthouse's TTI algorithm
// penalizes any late-arriving long task by extending its search for a
// quiet window, and back then the deferred work was still ~23.8s of
// main-thread cost regardless of when it started, so moving it later only
// pushed that penalty out further. That's fixed now (see
// flowery-branch.ts — draw calls cut from ~900 to 4), so deferral no
// longer fights a still-expensive task; gating on `load` also means this
// never even competes with the images we now know were the real cost
// (~4s per avif encode server-side — see ipx-cache.ts).
//
// showBackground used to share isLoading's own (much shorter)
// LOADING_TIMEOUT_MS fallback, on the theory that `load` firing late meant
// a hung request. In practice, a real PageSpeed mobile audit (Slow 4G +
// CPU throttling) showed `load` legitimately taking 10s+ — meaning that
// shared fallback fired first on every slow connection, starting the 3D
// scene's (still nontrivial, WebGL-heavy) first mount right in the middle
// of the page's real critical-path work instead of only as a rare safety
// net. BACKGROUND_TIMEOUT_MS below is its own, much longer fallback, kept
// separate from the loading screen's so a slow `load` no longer drags the
// background in early.
const showBackground = ref(false);
const BACKGROUND_TIMEOUT_MS = 10_000;

// Pages that wrap their content in the full-page white glass card
// (UiGlassPage.vue — blog, legal documents, the thank-you page) get a
// closer, larger view of the ambient branch instead of the normal
// small-in-the-distance framing — see AmbientTreeBackground's `closeup`
// prop. Persists across client-side navigation without remounting the 3D
// scene (it's a layout-level singleton, outside <slot/>), so this just
// needs to stay reactive to the current route.
const route = useRoute();
const GLASS_PAGE_PREFIXES = [
  "/blog",
  "/privacy",
  "/terms",
  "/cookie-policy",
  "/pd-consent",
  "/legal-info",
  "/thank-you",
];
const isGlassPage = computed(() =>
  GLASS_PAGE_PREFIXES.some(
    (prefix) => route.path === prefix || route.path.startsWith(`${prefix}/`),
  ),
);

// Loading screen hides the page (and blocks scroll) until `load` fires (or
// LOADING_TIMEOUT_MS elapses), so visitors never see hero/course-card
// images pop in piecemeal — `load` already means everything requested up
// front (every non-lazy `<img>`) has finished. `isLoading` defaults to true
// both during SSR and on the client's first render, so there's no
// hydration mismatch or flash of unhidden content before onMounted runs.
//
// Capped at LOADING_TIMEOUT_MS: `load` only fires once every requested
// resource finishes, so one slow image, a flaky network, or a hung
// third-party request would otherwise leave visitors staring at the
// spinner indefinitely instead of a page that's mostly ready. This is a
// real UX blocker (scroll is locked) that showBackground's own timeout
// doesn't share — the 3D background is invisible-until-shown and blocks
// nothing, so it has no equivalent reason to cut `load` short.
const LOADING_TIMEOUT_MS = 3000;
const isLoading = ref(true);

onMounted(() => {
  document.documentElement.classList.add("overflow-hidden");
  let loadingTimeoutId: ReturnType<typeof setTimeout> | undefined;
  let backgroundTimeoutId: ReturnType<typeof setTimeout> | undefined;

  const revealPage = () => {
    clearTimeout(loadingTimeoutId);
    document.documentElement.classList.remove("overflow-hidden");
    isLoading.value = false;
  };
  const startBackground = () => {
    clearTimeout(backgroundTimeoutId);
    showBackground.value = true;
  };
  const onLoad = () => {
    window.removeEventListener("load", onLoad);
    revealPage();
    startBackground();
  };

  if (document.readyState === "complete") {
    revealPage();
    startBackground();
    return;
  }
  window.addEventListener("load", onLoad);
  loadingTimeoutId = setTimeout(revealPage, LOADING_TIMEOUT_MS);
  backgroundTimeoutId = setTimeout(startBackground, BACKGROUND_TIMEOUT_MS);
});
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppLoadingScreen :loading="isLoading" />
    <LazyAmbientTreeBackground v-if="showBackground" :closeup="isGlassPage" />
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
    <CookieConsentBanner />
  </div>
</template>
