<script setup lang="ts">
// Loading screen hides the page (and blocks scroll) until `load` fires (or
// LOADING_TIMEOUT_MS elapses), so visitors never see hero/portfolio images
// pop in piecemeal — mirrors frontend/app/layouts/default.vue's own
// reasoning. The 3D ambient background that file also renders
// (AmbientTreeBackground.vue) is school-specific and deliberately not
// copied here — decor's own decorative background (falling snow, п. 5.2
// ТЗ) lands in Phase 4, built on the same fixed/pointer-events-none/
// prefers-reduced-motion pattern, not this file.
const LOADING_TIMEOUT_MS = 3000;
const isLoading = ref(true);

onMounted(() => {
  document.documentElement.classList.add("overflow-hidden");
  let loadingTimeoutId: ReturnType<typeof setTimeout> | undefined;

  const revealPage = () => {
    clearTimeout(loadingTimeoutId);
    document.documentElement.classList.remove("overflow-hidden");
    isLoading.value = false;
  };
  const onLoad = () => {
    window.removeEventListener("load", onLoad);
    revealPage();
  };

  if (document.readyState === "complete") {
    revealPage();
    return;
  }
  window.addEventListener("load", onLoad);
  loadingTimeoutId = setTimeout(revealPage, LOADING_TIMEOUT_MS);
});
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppLoadingScreen :loading="isLoading" />
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
    <CookieConsentBanner />
  </div>
</template>
