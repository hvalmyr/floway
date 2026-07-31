<script setup lang="ts">
/**
 * Minimal text-based stand-in for the VK logo — lucide-vue-next has no VK
 * icon, and hand-drawing the exact trademarked logomark is avoided on
 * purpose. Decorative only; pair with an aria-label on the surrounding link.
 *
 * @example
 * <a href="https://vk.com/floway" aria-label="VK"><IconVk class="size-[20px]" aria-hidden="true" /></a>
 */
</script>

<template>
  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="1" y="1" width="22" height="22" rx="6" stroke="currentColor" stroke-width="1.5" />
    <text x="12" y="16" text-anchor="middle" font-size="10" font-weight="700" fill="currentColor" font-family="sans-serif">
      VK
    </text>
  </svg>
</template>
