<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface Lesson {
  id: number;
  courseBlockId: number;
  number: number;
  title: string;
  topics: string;
  outcomes: string;
  durationHours: number;
}

const route = useRoute();
const blockId = route.params.blockId as string;

const emptyForm = (): Omit<Lesson, "id" | "courseBlockId"> => ({
  number: 0,
  title: "",
  topics: "",
  outcomes: "",
  durationHours: 0,
});

const { items, loading, error, fetchAll, create, update, remove } = useAdminResource<Lesson>(
  `/api/v1/course-blocks/${blockId}/lessons`,
);

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

function startEdit(lesson: Lesson) {
  editingId.value = lesson.id;
  form.value = {
    number: lesson.number,
    title: lesson.title,
    topics: lesson.topics,
    outcomes: lesson.outcomes,
    durationHours: lesson.durationHours,
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
      await create(form.value);
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
  if (!confirm("Удалить урок?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Уроки блока #{{ blockId }}</h1>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model.number="form.number"
        type="number"
        placeholder="Номер урока"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model="form.title"
        type="text"
        placeholder="Заголовок"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <textarea
        v-model="form.topics"
        placeholder="Темы"
        rows="3"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <textarea
        v-model="form.outcomes"
        placeholder="Результаты"
        rows="3"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <input
        v-model.number="form.durationHours"
        type="number"
        placeholder="Продолжительность (часы)"
        class="rounded border border-gray-300 px-3 py-2"
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
          <th class="px-4 py-2">№</th>
          <th class="px-4 py-2">Заголовок</th>
          <th class="px-4 py-2">Часы</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr v-for="lesson in items" :key="lesson.id" class="border-b border-gray-100 last:border-0">
          <td class="px-4 py-2">{{ lesson.number }}</td>
          <td class="px-4 py-2">{{ lesson.title }}</td>
          <td class="px-4 py-2">{{ lesson.durationHours }}</td>
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
          <td colspan="4" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
