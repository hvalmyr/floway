<script setup lang="ts">
import type { Component } from "vue";

/**
 * 3→2→1 column grid of icon+title+description cards ("Почему стоит учиться",
 * masterclass listing advantages). See docs/floway-design.md §5.2.
 *
 * @example
 * <FeatureGrid :items="[{ icon: Handshake, title: 'Гибкий старт', description: '...' }]" />
 */
defineProps<{
  items: { icon: Component; title: string; description: string }[];
}>();
</script>

<template>
  <div class="grid grid-cols-1 gap-24 md:grid-cols-2 lg:grid-cols-3 lg:gap-32">
    <div
      v-for="(item, i) in items"
      :key="i"
      class="flex flex-col items-center gap-16 rounded-md bg-white p-32 text-center"
    >
      <component :is="item.icon" class="size-40 text-primary-500" aria-hidden="true" />
      <h3 class="font-display text-h4 text-ink-900">{{ item.title }}</h3>
      <p class="text-body text-ink-700">{{ item.description }}</p>
    </div>
  </div>
</template>
