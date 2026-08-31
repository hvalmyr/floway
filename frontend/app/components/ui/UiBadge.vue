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
    class="inline-flex items-center rounded-pill border-2 border-primary bg-transparent px-[20px] py-12 font-body text-h4 font-medium text-primary"
  >
    <slot />
  </span>
</template>
