<script setup lang="ts">
/**
 * Pill label used for callouts like "работаем с 2013 года" or
 * "маленькие группы до 7 человек". Purely presentational text, no props.
 *
 * @example
 * <UiBadge>работаем с 2013 года</UiBadge>
 */
</script>

<template>
  <span
    class="inline-flex items-center rounded-pill border-[1.5px] border-primary bg-transparent px-[20px] py-8 font-body text-small font-bold text-primary"
  >
    <slot />
  </span>
</template>
