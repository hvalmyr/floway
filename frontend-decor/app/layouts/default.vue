<script setup lang="ts">
// Loading screen hides the page (and blocks scroll) until `load` fires (or
// LOADING_TIMEOUT_MS elapses), so visitors never see hero/portfolio images
// pop in piecemeal — mirrors frontend/app/layouts/default.vue's own
// reasoning. SnowBackground.vue is cheap CSS (no deferred gating needed),
// but AmbientTreeBackground.vue is the same WebGL-heavy component the
// school uses (recolored to decor's palette — see flowery-branch.ts) and
// gets the same load-deferred showBackground treatment for the same
// reason: see frontend/app/layouts/default.vue's own comment for the real
// Lighthouse numbers behind BACKGROUND_TIMEOUT_MS.
const LOADING_TIMEOUT_MS = 3000;
const BACKGROUND_TIMEOUT_MS = 10_000;
const isLoading = ref(true);
const showSnow = ref(false);
const showBackground = ref(false);

// Same closeup treatment as the school's own list, narrowed to the pages
// that actually exist on this site (no /blog here).
const route = useRoute();
const CLOSEUP_TREE_PREFIXES = [
  "/privacy",
  "/terms",
  "/cookie-policy",
  "/pd-consent",
  "/legal-info",
  "/thank-you",
  "/contacts",
];
const isCloseupPage = computed(() =>
  CLOSEUP_TREE_PREFIXES.some(
    (prefix) => route.path === prefix || route.path.startsWith(`${prefix}/`),
  ),
);

onMounted(() => {
  showSnow.value = true;
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
    <LazyAmbientTreeBackground v-if="showBackground" :closeup="isCloseupPage" />
    <LazySnowBackground v-if="showSnow" />
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
    <CookieConsentBanner />
  </div>
</template>
