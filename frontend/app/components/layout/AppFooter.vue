<script setup lang="ts">
// TODO(design): содержимое и вёрстка footer уточняются через Figma при вёрстке страниц.
const year = new Date().getFullYear()
</script>

<template>
  <footer class="mt-auto border-t border-black/5">
    <div class="mx-auto max-w-6xl px-4 py-8 text-sm text-[var(--color-text-muted)]">
      <p>&copy; {{ year }} Floway — школа флористики</p>
    </div>
  </footer>
</template>
