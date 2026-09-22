<script setup lang="ts">
/**
 * 3→2→1 column grid of icon+title+description cards ("Почему стоит учиться",
 * masterclass listing advantages). `icon` is a raw icon value (built-in key
 * or "icon:<id>") — see AppIcon.vue for how it resolves to a rendered icon.
 *
 * Uses flex-wrap (not CSS grid) so an incomplete last row centers itself
 * automatically. Below lg, items wrap freely (1 or 2 per row) via a flat
 * flex-wrap list. At lg, items are grouped into explicit row blocks of 3 —
 * except when that would leave a single dangling card (items.length % 3 ===
 * 1), in which case the last two rows are 2+2 instead of 3+1. Cards always
 * keep the 3-column width, never stretched; each row's `contents` display
 * below lg removes the grouping so the flat wrap is unaffected.
 *
 * @example
 * <FeatureGrid :items="[{ icon: 'flex-start', title: 'Гибкий старт', description: '...' }]" />
 */
const props = defineProps<{
  items: { icon: string; title: string; description: string }[];
}>();

const rows = computed(() => {
  const { items } = props;
  const n = items.length;
  const result: (typeof items)[] = [];
  let i = 0;
  if (n > 1 && n % 3 === 1) {
    while (n - i > 4) {
      result.push(items.slice(i, i + 3));
      i += 3;
    }
    result.push(items.slice(i, i + 2));
    result.push(items.slice(i + 2, i + 4));
  } else {
    while (i < n) {
      result.push(items.slice(i, i + 3));
      i += 3;
    }
  }
  return result;
});
</script>

<template>
  <div class="flex flex-wrap items-start justify-center gap-24 lg:flex-col lg:items-center lg:gap-32">
    <div
      v-for="(row, r) in rows"
      :key="r"
      class="contents lg:flex lg:w-full lg:flex-wrap lg:justify-center lg:gap-32"
    >
      <div
        v-for="(item, i) in row"
        :key="i"
        class="flex w-full flex-col items-center gap-16 rounded-md bg-white p-32 text-center md:w-[calc((100%_-_24px)/2)] lg:w-[calc((100%_-_64px)/3)]"
      >
        <!-- Фиксированная высота (а не size-*), т.к. иконки разной ширины должны
        выглядеть одной высоты, а не влезать в одинаковый квадрат. -->
        <AppIcon :icon="item.icon" class="h-40 w-auto text-primary" />
        <h3 class="font-display text-h4 text-primary">{{ item.title }}</h3>
        <p class="whitespace-pre-line font-body text-body text-ink">{{ item.description }}</p>
      </div>
    </div>
  </div>
</template>
