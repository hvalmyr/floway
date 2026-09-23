<script setup lang="ts">
import { onMounted, ref } from "vue";

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
// spinner indefinitely instead of a page that's mostly ready.
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
