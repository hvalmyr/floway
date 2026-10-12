<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Tabs shown above the masterclasses list on the public page (see
// masterclasses/index.vue) — flat CRUD, same shape as the decor site's
// object-types admin screen, minus "visible" (every category is its own
// tab; there's no hidden/draft state here).
interface MasterclassCategory {
  id: number;
  slug: string;
  name: string;
  sortOrder: number;
}

const emptyForm = (): Omit<MasterclassCategory, "id"> => ({
  slug: "",
  name: "",
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<MasterclassCategory>("/api/v1/masterclass-categories");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

const slugTouched = ref(false);
watch(
  () => form.value.name,
  (name) => {
    if (!slugTouched.value) form.value.slug = slugify(name);
  },
);

function startEdit(category: MasterclassCategory) {
  slugTouched.value = true;
  editingId.value = category.id;
  form.value = {
    slug: category.slug,
    name: category.name,
    sortOrder: category.sortOrder,
  };
}

function cancelEdit() {
  editingId.value = null;
  form.value = emptyForm();
  slugTouched.value = false;
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
  if (!confirm("Удалить категорию? Мастер-классы этой категории останутся без категории.")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Категории мастер-классов</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Категории отображаются вкладками над списком мастер-классов на сайте.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.name"
        type="text"
        placeholder="Название (например, «Новогодние»)"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model="form.slug"
        type="text"
        placeholder="Slug"
        required
        class="rounded border border-gray-300 px-3 py-2"
        @input="slugTouched = true"
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
          <th class="px-4 py-2">Slug</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(category, index) in items"
          :key="category.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">{{ category.name }}</td>
          <td class="px-4 py-2">{{ category.slug }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button
              class="text-[var(--color-primary)] hover:underline"
              @click="startEdit(category)"
            >
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(category.id)">
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
