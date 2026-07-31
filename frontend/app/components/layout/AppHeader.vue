<script setup lang="ts">
// TODO(design): вёрстка ждёт токенов/структуры из Figma (phone: 0:477, pad: 0:483,
// desktop-main: 0:15 — там проработан header). Пока — минимальная заглушка с навигацией.
const links = [
  { to: '/', label: 'Главная' },
  { to: '/masterclasses', label: 'Мастер-классы' },
  { to: '/blog', label: 'Блог' },
  { to: '/contacts', label: 'Контакты' }
]
</script>

<template>
  <header class="border-b border-black/5">
    <div class="mx-auto flex max-w-6xl items-center justify-between px-4 py-4">
      <NuxtLink to="/" class="text-xl font-semibold" style="font-family: var(--font-family-heading)">
        Floway
      </NuxtLink>
      <nav class="flex gap-6">
        <NuxtLink
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          class="text-sm text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>
    </div>
  </header>
</template>
