<script setup lang="ts">
/**
 * Price/duration info panel (docs/floway-design.md §5.2 "инфо-плашка цены").
 * A horizontal row of 2–4 values on desktop that stacks vertically on
 * mobile, per §7 of the design doc.
 *
 * @example
 * <UiInfoRow
 *   :items="[
 *     { label: 'Длительность', value: '2–3 часа' },
 *     { label: 'Цена (групповое)', value: '3 000 ₽' },
 *     { label: 'Цена (индивидуальное)', value: '4 500 ₽' },
 *   ]"
 * />
 */
defineProps<{
  items: { label: string; value: string; oldValue?: string }[];
}>();
</script>

<template>
  <div
    class="flex flex-col divide-y divide-line rounded-md border border-line bg-surface md:flex-row md:divide-x md:divide-y-0"
  >
    <div v-for="item in items" :key="item.label" class="flex flex-1 flex-col gap-4 px-24 py-16">
      <span class="text-small text-ink-700">{{ item.label }}</span>
      <span class="font-display text-h4 text-ink-900">
        <span v-if="item.oldValue" class="mr-8 text-body font-body text-ink-400 line-through">
          <span class="sr-only">Старая цена: </span>{{ item.oldValue }}
        </span>
        {{ item.value }}
      </span>
    </div>
  </div>
</template>
