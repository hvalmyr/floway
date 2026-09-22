<script setup lang="ts">
/**
 * 3→2→1 column grid of icon+title+description cards ("Почему стоит учиться",
 * masterclass listing advantages). `icon` is a raw icon value (built-in key
 * or "icon:<id>") — see AppIcon.vue for how it resolves to a rendered icon.
 *
 * Uses flex-wrap (not CSS grid) so an incomplete last row centers itself
 * automatically. At the lg breakpoint cards always keep the 3-column width,
 * even when the last row has only 1 or 2 cards — they stay centered at
 * their normal size instead of stretching to fill the row.
 *
 * @example
 * <FeatureGrid :items="[{ icon: 'flex-start', title: 'Гибкий старт', description: '...' }]" />
 */
defineProps<{
  items: { icon: string; title: string; description: string }[];
}>();
</script>

<template>
  <div class="flex flex-wrap items-start justify-center gap-24 lg:gap-32">
    <div
      v-for="(item, i) in items"
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
</template>
