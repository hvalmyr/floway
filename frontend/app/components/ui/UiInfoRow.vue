<script setup lang="ts">
/**
 * Price/duration info panel (course price table rows, masterclass duration+
 * price bar). A single row of plain-text columns — first column in blue
 * (label-like), the rest in bold ink — that stacks vertically on mobile.
 * `highlighted` switches from the plain surface fill to the white+blue-
 * border treatment used for the second row in the course price table.
 *
 * @example
 * <UiInfoRow :items="['Блок “Букеты”', '7 занятий', '30 часов', '38 500 ₽']" />
 * <UiInfoRow :items="['Блок “Композиции”', '4 занятия', '17 часов', '22 000 ₽']" highlighted />
 */
withDefaults(
  defineProps<{
    items: string[];
    highlighted?: boolean;
  }>(),
  { highlighted: false },
);
</script>

<template>
  <div
    class="flex flex-col divide-y divide-primary/20 rounded-md md:flex-row md:items-center md:divide-x md:divide-y-0"
    :class="highlighted ? 'border-2 border-primary bg-white' : 'bg-surface'"
  >
    <div
      v-for="(item, i) in items"
      :key="i"
      class="flex-1 px-24 py-16 font-body text-body"
      :class="i === 0 ? 'font-bold text-primary' : 'font-bold text-ink'"
    >
      {{ item }}
    </div>
  </div>
</template>
