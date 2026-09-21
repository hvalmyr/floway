<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface EditingGroup {
  title: string;
  links: { to: string; label: string }[];
}

const editingGroups: EditingGroup[] = [
  {
    title: "Услуга",
    links: [
      { to: "/admin/object-types", label: "Типы объектов" },
      { to: "/admin/landing-pages", label: "Посадочные страницы" },
      { to: "/admin/portfolio", label: "Портфолио" },
    ],
  },
  {
    title: "Компоненты",
    links: [
      { to: "/admin/faq", label: "FAQ главной" },
      { to: "/admin/features", label: "Преимущества" },
      { to: "/admin/thank-you-page", label: "Страница благодарности" },
    ],
  },
  {
    title: "Информация",
    links: [
      { to: "/admin/social-links", label: "Соцсети" },
      { to: "/admin/notification-emails", label: "Email для уведомлений о заявках" },
      { to: "/admin/page-content", label: "Все тексты сайта (полный список)" },
    ],
  },
];

const editingOpen = ref(false);
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
      <button
        type="button"
        class="rounded border bg-white px-4 py-8 text-center text-lg font-medium hover:border-[var(--color-primary)]"
        :class="editingOpen ? 'border-[var(--color-primary)]' : 'border-gray-200'"
        @click="editingOpen = !editingOpen"
      >
        Редактирование
      </button>
    </div>

    <div v-if="editingOpen" class="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
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
