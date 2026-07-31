<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

const sections = [
  { to: "/admin/teachers", title: "Преподаватели" },
  { to: "/admin/blog-posts", title: "Блог" },
  { to: "/admin/masterclasses", title: "Мастер-классы" },
  { to: "/admin/courses", title: "Курсы" },
  { to: "/admin/faq", title: "FAQ" },
  { to: "/admin/leads", title: "Заявки" },
];
</script>

<template>
  <div>
    <h1 class="text-3xl font-semibold">Админ-панель</h1>
    <div class="mt-8 grid grid-cols-2 gap-4 sm:grid-cols-3">
      <NuxtLink
        v-for="section in sections"
        :key="section.to"
        :to="section.to"
        class="rounded border border-gray-200 bg-white px-4 py-6 text-center font-medium hover:border-[var(--color-primary)]"
      >
        {{ section.title }}
      </NuxtLink>
    </div>
  </div>
</template>
