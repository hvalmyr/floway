<script setup lang="ts">
/**
 * Bottom-of-viewport notice about cookie/analytics usage, shown until the
 * visitor accepts (state in useCookieConsent.ts) — accept-only, no decline
 * option, so the banner just keeps showing and analytics stays off until
 * they do. Yandex.Metrika only loads on accept — immediately on mount for
 * a returning visitor who already accepted, or the moment the button is
 * clicked for a first-time one.
 */
const { status, accept } = useCookieConsent();
const { load: loadYandexMetrika } = useYandexMetrika();

onMounted(() => {
  if (status.value === "accepted") loadYandexMetrika();
});

watch(status, (value) => {
  if (value === "accepted") loadYandexMetrika();
});
</script>

<template>
  <div
    v-if="status === null"
    class="container fixed inset-x-0 bottom-16 z-50 animate-[cookie-banner-in_0.3s_ease-out] sm:bottom-24"
    role="region"
    aria-label="Уведомление об использовании cookie"
  >
    <div
      class="flex w-full flex-col items-start gap-12 rounded-md border-2 border-ink bg-white/40 p-16 backdrop-blur-lg backdrop-saturate-150 lg:flex-row lg:items-center lg:justify-between"
    >
      <p class="font-body text-body text-ink">
        Мы используем cookie и аналитику. Подробнее — в
        <NuxtLink to="/privacy" class="underline hover:no-underline"
          >политике обработки данных</NuxtLink
        >.
      </p>
      <div class="w-full shrink-0 lg:w-auto">
        <UiButton variant="primary" class="w-full lg:w-auto" @click="accept">Принять</UiButton>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes cookie-banner-in {
  from {
    transform: translateY(100%);
    opacity: 0;
  }
}
</style>
