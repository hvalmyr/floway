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
const showBackground = ref(false);

// Loading screen hides the page (and blocks scroll) until the same `load`
// event above fires, so visitors never see hero/course-card images pop in
// piecemeal — `load` already means everything requested up front (every
// non-lazy `<img>`) has finished. `isLoading` defaults to true both during
// SSR and on the client's first render, so there's no hydration mismatch or
// flash of unhidden content before onMounted runs.
const isLoading = ref(true);

onMounted(() => {
  document.documentElement.classList.add("overflow-hidden");
  const finishLoading = () => {
    document.documentElement.classList.remove("overflow-hidden");
    showBackground.value = true;
    isLoading.value = false;
  };
  if (document.readyState === "complete") {
    finishLoading();
    return;
  }
  window.addEventListener("load", finishLoading);
});
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppLoadingScreen :loading="isLoading" />
    <LazyAmbientTreeBackground v-if="showBackground" />
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
    <CookieConsentBanner />
  </div>
</template>
