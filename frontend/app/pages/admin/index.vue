<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface EditingGroup {
  title: string;
  links: { to: string; label: string }[];
}

const editingGroups: EditingGroup[] = [
  {
    title: "Учебная программа",
    links: [
      { to: "/admin/course-sections", label: "Курсы" },
      { to: "/admin/masterclasses", label: "Мастер-классы" },
    ],
  },
  {
    title: "Компоненты",
    links: [
      { to: "/admin/faq", label: "FAQ" },
      { to: "/admin/features", label: "Преимущества" },
      { to: "/admin/page-content/hero", label: "Hero-блок" },
      { to: "/admin/gallery-photos", label: "Фотогалерея" },
    ],
  },
  {
    title: "Главная",
    links: [
      { to: "/admin/page-content/home", label: "Пробное занятие" },
      { to: "/admin/teachers", label: "Преподаватели" },
      { to: "/admin/about-items", label: "О школе" },
      { to: "/admin/page-content/apply-form", label: "Форма заявки" },
    ],
  },
  {
    title: "Информация",
    links: [
      { to: "/admin/page-content/info", label: "Контакты и реквизиты" },
      { to: "/admin/page-content", label: "Все тексты сайта (полный список)" },
      { to: "/admin/content-export", label: "Экспорт / импорт" },
    ],
  },
];

const editingOpen = ref(false);
</script>

<template>
  <div>
    <h1 class="text-3xl font-semibold">Админ-панель</h1>

    <div class="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-3">
      <NuxtLink
        to="/admin/leads"
        class="rounded border border-gray-200 bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
      >
        Заявки
      </NuxtLink>
      <button
        type="button"
        class="rounded border bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
        :class="editingOpen ? 'border-[var(--color-primary)]' : 'border-gray-200'"
        @click="editingOpen = !editingOpen"
      >
        Редактирование
      </button>
      <NuxtLink
        to="/admin/blog-posts"
        class="rounded border border-gray-200 bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
      >
        Блог
      </NuxtLink>
    </div>

    <div v-if="editingOpen" class="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2">
      <div
        v-for="group in editingGroups"
        :key="group.title"
        class="rounded border border-gray-200 bg-white p-4"
      >
        <p class="mb-3 font-semibold">{{ group.title }}</p>
        <div class="flex flex-col gap-2">
          <NuxtLink
            v-for="link in group.links"
            :key="link.to"
            :to="link.to"
            class="text-sm text-[var(--color-primary)] hover:underline"
          >
            {{ link.label }}
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
