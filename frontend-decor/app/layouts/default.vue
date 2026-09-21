<script setup lang="ts">
// Loading screen hides the page (and blocks scroll) until `load` fires (or
// LOADING_TIMEOUT_MS elapses), so visitors never see hero/portfolio images
// pop in piecemeal — mirrors frontend/app/layouts/default.vue's own
// reasoning. The school's 3D ambient background (AmbientTreeBackground.vue)
// isn't copied here — SnowBackground.vue is decor's own equivalent (п. 5.2
// ТЗ), simple CSS rather than WebGL, so it doesn't need the same deferred-
// until-`load` gating; it's still mounted client-only (Lazy+v-if, flipped
// in onMounted) purely to keep its randomized per-flake styles out of SSR.
const LOADING_TIMEOUT_MS = 3000;
const isLoading = ref(true);
const showSnow = ref(false);

onMounted(() => {
  showSnow.value = true;
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
    <LazySnowBackground v-if="showSnow" />
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
    <CookieConsentBanner />
  </div>
</template>
