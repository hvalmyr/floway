<script setup lang="ts">
import type { HomeSection, HomeSectionKey } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

const HOME_SECTION_LABELS: Record<HomeSectionKey, string> = {
  features: "Почему стоит учиться (преимущества)",
  gallery: "Фотокарусель",
  courses: "Курсы (секции курсов)",
  gift_certificates: "Бегущая строка сертификатов",
  trial: "Пробное занятие",
  about: "О школе",
  teachers: "Педагоги",
  reviews: "Отзывы (Яндекс Карты)",
  faq: "Вопросы и ответы (FAQ)",
};

const { items, loading, error, fetchAll, update } =
  useAdminResource<HomeSection>("/api/v1/home-sections");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

async function onToggleVisible(section: HomeSection) {
  await update(section.id, { ...section, visible: !section.visible });
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Порядок секций главной страницы</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Перетаскивайте строки, чтобы изменить порядок блоков на главной странице, и переключайте
      видимость, чтобы скрыть блок без удаления. Hero-блок (самый верх страницы) всегда показывается
      первым и здесь не настраивается — его текст редактируется на
      <NuxtLink to="/admin/page-content/hero" class="text-[var(--color-primary)] hover:underline"
        >отдельной странице</NuxtLink
      >.
    </p>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <table
      v-else
      class="mt-6 w-full border-collapse overflow-hidden rounded border border-gray-200 bg-white text-sm"
    >
      <thead>
        <tr class="border-b border-gray-200 text-left">
          <th class="px-4 py-2" />
          <th class="px-4 py-2">Блок</th>
          <th class="px-4 py-2">Видимость</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(section, index) in items"
          :key="section.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">{{ HOME_SECTION_LABELS[section.key] ?? section.key }}</td>
          <td class="px-4 py-2">
            <AdminVisibilityDot :visible="section.visible" @click="onToggleVisible(section)" />
          </td>
        </tr>
        <tr v-if="items.length === 0">
          <td colspan="3" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
