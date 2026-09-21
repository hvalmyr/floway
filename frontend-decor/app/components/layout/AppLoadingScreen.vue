<script setup lang="ts">
defineProps<{
  loading: boolean;
}>();
</script>

<template>
  <Transition name="loading-fade">
    <div
      v-if="loading"
      class="fixed inset-0 z-[200] flex items-center justify-center bg-surface"
      aria-hidden="true"
    >
      <div class="flex gap-12">
        <span class="loading-dot size-12 rounded-pill bg-primary" style="animation-delay: 0ms" />
        <span class="loading-dot size-12 rounded-pill bg-primary" style="animation-delay: 150ms" />
        <span class="loading-dot size-12 rounded-pill bg-primary" style="animation-delay: 300ms" />
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.loading-dot {
  animation: loading-dot-pulse 1.1s ease-in-out infinite;
}

@keyframes loading-dot-pulse {
  0%,
  80%,
  100% {
    opacity: 0.3;
    transform: scale(0.75);
  }
  40% {
    opacity: 1;
    transform: scale(1);
  }
}

.loading-fade-enter-active,
.loading-fade-leave-active {
  transition: opacity 0.4s ease;
}
.loading-fade-enter-from,
.loading-fade-leave-to {
  opacity: 0;
}
</style>
