<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface Lesson {
  id: number;
  courseId: number;
  name: string;
  description: string;
  sortOrder: number;
}

const route = useRoute();
const courseId = route.params.courseId as string;

const emptyForm = (): Omit<Lesson, "id" | "courseId"> => ({
  name: "",
  description: "",
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } = useAdminResource<Lesson>(
  `/api/v1/courses/${courseId}/lessons`,
);

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onDragStart, onDragOver, onDrop } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(lesson: Lesson) {
  editingId.value = lesson.id;
  form.value = {
    name: lesson.name,
    description: lesson.description,
    sortOrder: lesson.sortOrder,
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
  if (!confirm("Удалить занятие?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <NuxtLink
      to="/admin/course-sections"
      class="text-sm text-[var(--color-primary)] hover:underline"
      >← К списку секций</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Занятия курса #{{ courseId }}</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Для курса без блоков — занятия привязаны прямо к курсу. Если у курса есть блоки, занятия
      редактируются у каждого блока отдельно (вкладка «Блоки» → «Занятия»), а этот список не
      используется.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.name"
        type="text"
        placeholder="Название занятия (например, «Занятие 1. Спиральная техника»)"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <AdminMarkdownField
        v-model="form.description"
        placeholder="Описание занятия"
        :rows="4"
        class="sm:col-span-2"
      />
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
          <th class="px-4 py-2">Название</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(lesson, index) in items"
          :key="lesson.id"
          draggable="true"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
          @dragstart="onDragStart(index)"
          @dragover.prevent="onDragOver(index)"
          @drop.prevent="onDrop"
        >
          <td class="px-4 py-2"><AdminDragHandle /></td>
          <td class="px-4 py-2">{{ lesson.name }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(lesson)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(lesson.id)">
              Удалить
            </button>
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
