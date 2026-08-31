<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface CourseSection {
  id: number;
  heading: string;
  description: string;
  sortOrder: number;
  visible: boolean;
}

const emptyForm = (): Omit<CourseSection, "id"> => ({
  heading: "",
  description: "",
  sortOrder: 0,
  visible: true,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<CourseSection>("/api/v1/course-sections");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onDragStart, onDragOver, onDrop } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(section: CourseSection) {
  editingId.value = section.id;
  form.value = {
    heading: section.heading,
    description: section.description,
    sortOrder: section.sortOrder,
    visible: section.visible,
  };
}

function cancelEdit() {
  editingId.value = null;
  form.value = emptyForm();
}

async function onSubmit() {
  formError.value = "";
  saving.value = true;
  try {
    if (editingId.value === null) {
      await create({ ...form.value, sortOrder: items.value.length });
    } else {
      await update(editingId.value, form.value);
    }
    cancelEdit();
  } catch {
    formError.value = "Не удалось сохранить. Проверьте поля и попробуйте снова.";
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Удалить секцию курсов вместе со всеми курсами внутри неё?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

async function onToggleVisible(section: CourseSection) {
  await update(section.id, { ...section, visible: !section.visible });
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Секции курсов (главная страница)</h1>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.heading"
        type="text"
        placeholder="Заголовок секции"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <AdminMarkdownField
        v-model="form.description"
        placeholder="Описание секции"
        :rows="2"
        class="sm:col-span-2"
      />
      <label class="flex items-center gap-2 text-sm">
        <input v-model="form.visible" type="checkbox" class="size-5" />
        Показывать на сайте
      </label>

      <p v-if="formError" class="text-sm text-red-600 sm:col-span-2">{{ formError }}</p>

      <div class="flex gap-2 sm:col-span-2">
        <button
          type="submit"
          :disabled="saving"
          class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
        >
          {{ editingId === null ? "Добавить" : "Сохранить" }}
        </button>
        <button
          v-if="editingId !== null"
          type="button"
          class="rounded border border-gray-300 px-4 py-2"
          @click="cancelEdit"
        >
          Отмена
        </button>
      </div>
    </form>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <table
      v-else
      class="mt-6 w-full border-collapse overflow-hidden rounded border border-gray-200 bg-white text-sm"
    >
      <thead>
        <tr class="border-b border-gray-200 text-left">
          <th class="px-4 py-2" />
          <th class="px-4 py-2">Заголовок</th>
          <th class="px-4 py-2">Видимость</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(section, index) in items"
          :key="section.id"
          draggable="true"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
          @dragstart="onDragStart(index)"
          @dragover.prevent="onDragOver(index)"
          @drop.prevent="onDrop"
        >
          <td class="px-4 py-2"><AdminDragHandle /></td>
          <td class="px-4 py-2">{{ section.heading }}</td>
          <td class="px-4 py-2">
            <AdminVisibilityDot :visible="section.visible" @click="onToggleVisible(section)" />
          </td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <NuxtLink
              :to="`/admin/course-sections/${section.id}/courses`"
              class="text-[var(--color-primary)] hover:underline"
              >Курсы</NuxtLink
            >
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(section)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(section.id)">
              Удалить
            </button>
          </td>
        </tr>
        <tr v-if="items.length === 0">
          <td colspan="4" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
