<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Content-editing screens (FAQ, преимущества, тексты страниц, портфолио,
// соцсети, уведомления, экспорт) land in Phase 4 alongside the public pages
// that actually consume them — no point wiring up a page-content editor
// before the page reading it exists.
</script>

<template>
  <div>
    <h1 class="text-3xl font-semibold">Админ-панель flo-way</h1>

    <div class="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <NuxtLink
        to="/admin/leads"
        class="rounded border border-gray-200 bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
      >
        Заявки
      </NuxtLink>
      <NuxtLink
        to="/admin/clients"
        class="rounded border border-gray-200 bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
      >
        Клиенты
      </NuxtLink>
      <NuxtLink
        to="/admin/object-types"
        class="rounded border border-gray-200 bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
      >
        Типы объектов
      </NuxtLink>
      <NuxtLink
        to="/admin/landing-pages"
        class="rounded border border-gray-200 bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
      >
        Посадочные страницы
      </NuxtLink>
    </div>
  </div>
</template>
