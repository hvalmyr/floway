<script setup lang="ts">
/**
 * Price/duration info panel — one row per course block (or the single
 * synthetic block) on the course page. A single row of plain-text columns,
 * all one uniform color — a 2-column grid below `lg`, one row above it (see
 * the responsive-grid note further down). Rows alternate
 * two full styles via `highlighted`: beige fill + brown text (default) vs.
 * white fill + blue text + blue border (`highlighted` — the border is
 * needed here specifically because a plain white fill would otherwise have
 * no visible edge against this section's own white background, unlike the
 * beige default which already contrasts on its own) — the caller cycles it
 * per row (see courses/[slug].vue's `i % 2 === 1`). Items sit in equal-width
 * columns (`auto-cols-fr`), left-aligned — varying text lengths previously
 * left items scattered unevenly across the row under `justify-between`.
 * Below `lg` there isn't room for 4 columns without either wrapping
 * mid-word (tablet) or stretching into a very tall single-file stack
 * (mobile), so it's a 2-column grid instead — an odd-count row (e.g. no
 * `blockLabel`) spans its lone last item across both columns rather than
 * leaving an empty cell dangling next to it.
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
    class="grid grid-cols-2 rounded-md lg:grid-flow-col lg:grid-cols-none lg:auto-cols-fr lg:items-center"
    :class="highlighted ? 'border-2 border-primary bg-white text-primary' : 'bg-surface text-ink'"
  >
    <div
      v-for="(item, i) in items"
      :key="i"
      class="px-16 py-24 text-left font-body text-body font-bold lg:px-24 lg:py-32"
      :class="i === items.length - 1 && items.length % 2 === 1 ? 'col-span-2 lg:col-span-1' : ''"
    >
      {{ item }}
    </div>
  </div>
</template>
