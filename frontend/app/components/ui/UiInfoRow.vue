<script setup lang="ts">
/**
 * Price/duration info panel — one row per course block (or the single
 * synthetic block) on the course page. A single row of plain-text columns,
 * all one uniform color — that stacks vertically on mobile. Rows alternate
 * two full styles via `highlighted`: beige fill + brown text (default) vs.
 * white fill + blue text + blue border (`highlighted` — the border is
 * needed here specifically because a plain white fill would otherwise have
 * no visible edge against this section's own white background, unlike the
 * beige default which already contrasts on its own) — the caller cycles it
 * per row (see courses/[slug].vue's `i % 2 === 1`). Items sit in equal-width
 * columns (`auto-cols-fr`), left-aligned — varying text lengths previously
 * left items scattered unevenly across the row under `justify-between`.
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
    class="grid rounded-md md:grid-flow-col md:auto-cols-fr md:items-center"
    :class="highlighted ? 'border-2 border-primary bg-white text-primary' : 'bg-surface text-ink'"
  >
    <div
      v-for="(item, i) in items"
      :key="i"
      class="px-24 py-32 text-left font-body text-body font-bold"
    >
      {{ item }}
    </div>
  </div>
</template>
