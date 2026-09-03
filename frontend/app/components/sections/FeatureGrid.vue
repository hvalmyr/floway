<script setup lang="ts">
import type { Component } from "vue";

/**
 * 3→2→1 column grid of icon+title+description cards ("Почему стоит учиться",
 * masterclass listing advantages).
 *
 * @example
 * <FeatureGrid :items="[{ icon: IconFlexStart, title: 'Гибкий старт', description: '...' }]" />
 */
defineProps<{
  items: { icon: Component; title: string; description: string }[];
}>();
</script>

<template>
  <div class="grid grid-cols-1 items-start gap-24 md:grid-cols-2 lg:grid-cols-3 lg:gap-32">
    <div
      v-for="(item, i) in items"
      :key="i"
      class="flex flex-col items-center gap-16 rounded-md bg-white p-32 text-center"
    >
      <!-- Фиксированная высота (а не size-*), т.к. иконки разной ширины должны
      выглядеть одной высоты, а не влезать в одинаковый квадрат. -->
      <component :is="item.icon" class="h-40 w-auto text-primary" aria-hidden="true" />
      <h3 class="font-display text-h4 text-primary">{{ item.title }}</h3>
      <p class="whitespace-pre-line font-body text-body text-ink">{{ item.description }}</p>
    </div>
  </div>
</template>
