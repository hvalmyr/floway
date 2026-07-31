<script setup lang="ts">
const { adminUser, logout } = useAdminAuth();

async function onLogout() {
  await logout();
  await navigateTo("/admin/login");
}
</script>

<template>
  <div class="min-h-screen bg-[var(--color-background)]">
    <header v-if="adminUser" class="border-b border-gray-200 bg-white px-6 py-4">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4">
        <nav class="flex flex-wrap gap-4 text-sm">
          <NuxtLink to="/admin">Дашборд</NuxtLink>
          <NuxtLink to="/admin/teachers">Преподаватели</NuxtLink>
          <NuxtLink to="/admin/blog-posts">Блог</NuxtLink>
          <NuxtLink to="/admin/masterclasses">Мастер-классы</NuxtLink>
          <NuxtLink to="/admin/courses">Курсы</NuxtLink>
          <NuxtLink to="/admin/faq">FAQ</NuxtLink>
          <NuxtLink to="/admin/leads">Заявки</NuxtLink>
        </nav>
        <div class="flex items-center gap-3 text-sm">
          <span class="text-[var(--color-text-muted)]">{{ adminUser.login }}</span>
          <button class="text-[var(--color-primary)] hover:underline" @click="onLogout">
            Выйти
          </button>
        </div>
      </div>
    </header>
    <main :class="adminUser ? 'mx-auto max-w-6xl px-6 py-8' : ''">
      <slot />
    </main>
  </div>
</template>
