<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface CourseBlock {
  id: number;
  courseId: number;
  title: string;
  lessonsCount: number;
  hours: number;
  price: number;
  sortOrder: number;
}

const route = useRoute();
const courseId = route.params.courseId as string;

const emptyForm = (): Omit<CourseBlock, "id" | "courseId"> => ({
  title: "",
  lessonsCount: 0,
  hours: 0,
  price: 0,
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } = useAdminResource<CourseBlock>(
  `/api/v1/courses/${courseId}/blocks`,
);

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

function startEdit(block: CourseBlock) {
  editingId.value = block.id;
  form.value = {
    title: block.title,
    lessonsCount: block.lessonsCount,
    hours: block.hours,
    price: block.price,
    sortOrder: block.sortOrder,
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
  if (!confirm("Удалить блок курса?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <NuxtLink to="/admin/courses" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К списку курсов</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Блоки курса #{{ courseId }}</h1>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.title"
        type="text"
        placeholder="Заголовок"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <input
        v-model.number="form.lessonsCount"
        type="number"
        placeholder="Количество уроков"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model.number="form.hours"
        type="number"
        placeholder="Часы"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model.number="form.price"
        type="number"
        placeholder="Цена"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model.number="form.sortOrder"
        type="number"
        placeholder="Порядок сортировки"
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
          <th class="px-4 py-2">Заголовок</th>
          <th class="px-4 py-2">Часы</th>
          <th class="px-4 py-2">Цена</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr v-for="block in items" :key="block.id" class="border-b border-gray-100 last:border-0">
          <td class="px-4 py-2">{{ block.title }}</td>
          <td class="px-4 py-2">{{ block.hours }}</td>
          <td class="px-4 py-2">{{ block.price }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <NuxtLink
              :to="`/admin/course-blocks/${block.id}/lessons`"
              class="text-[var(--color-primary)] hover:underline"
              >Уроки</NuxtLink
            >
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(block)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(block.id)">
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
