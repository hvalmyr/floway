<script setup lang="ts">
import { onMounted, ref } from "vue";

// The 3D background is purely decorative (fixed, -z-10, pointer-events-none
// — see AmbientTreeBackground.vue) and safe to start a moment late, but a
// real mobile Lighthouse audit found it costs ~13.6s of main-thread work
// on its own (downloading/parsing/executing the three.js bundle, then
// synchronously building the branch geometry and compiling shaders) —
// rendering it unconditionally made all of that part of the CRITICAL
// initial page load, competing with FCP/LCP/TBT for the same main thread.
// LazyAmbientTreeBackground defers the chunk import itself until v-if
// flips true; requestIdleCallback then waits for the main thread to
// actually go quiet (i.e. after real content has painted) before that
// happens, with a plain timeout fallback for Safari (no
// requestIdleCallback) and for the case where the thread never truly
// idles within a couple seconds.
const showBackground = ref(false);
onMounted(() => {
  const start = () => {
    showBackground.value = true;
  };
  if ("requestIdleCallback" in window) {
    requestIdleCallback(start, { timeout: 2000 });
  } else {
    setTimeout(start, 200);
  }
});
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <LazyAmbientTreeBackground v-if="showBackground" />
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
    <CookieConsentBanner />
  </div>
</template>
