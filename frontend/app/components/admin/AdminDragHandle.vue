<script setup lang="ts">
import { GripVertical } from "lucide-vue-next";
</script>

<template>
  <span
    class="inline-flex cursor-grab rounded p-1 text-[var(--color-text-muted)] hover:bg-gray-100 active:cursor-grabbing"
    title="Перетащите, чтобы изменить порядок"
  >
    <GripVertical class="size-6" />
  </span>
</template>
