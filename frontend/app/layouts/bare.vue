<script setup lang="ts">
import { computed } from "vue";

const { treeEnabled } = await useTreeMode();

useHead({
  htmlAttrs: { class: computed(() => (treeEnabled.value ? "tree-enabled" : "")) },
});
</script>

<template>
  <div class="min-h-screen">
    <LazyAmbientTreeBackground v-if="treeEnabled" />
    <main>
      <slot />
    </main>
  </div>
</template>
