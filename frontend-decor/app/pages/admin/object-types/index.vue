<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Shared dictionary of property types — used by landing pages, the lead
// form, and the portfolio filter (see docs/decor-site-audit.md §2). Flat
// CRUD, same shape as school's masterclasses admin screen.
interface ObjectType {
  id: number;
  slug: string;
  name: string;
  visible: boolean;
  sortOrder: number;
}

const emptyForm = (): Omit<ObjectType, "id"> => ({
  slug: "",
  name: "",
  visible: true,
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<ObjectType>("/api/v1/object-types");

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

function startEdit(objectType: ObjectType) {
  slugTouched.value = true;
  editingId.value = objectType.id;
  form.value = {
    slug: objectType.slug,
    name: objectType.name,
    visible: objectType.visible,
    sortOrder: objectType.sortOrder,
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
  if (!confirm("Удалить тип объекта? Посадочные страницы этого типа удалятся вместе с ним."))
    return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

async function onToggleVisible(objectType: ObjectType) {
  await update(objectType.id, { ...objectType, visible: !objectType.visible });
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Типы объектов</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Один справочник для формы заявки, фильтра портфолио и посадочных страниц — изменение здесь
      обновляет все три места сразу.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.name"
        type="text"
        placeholder="Название (например, «Дом»)"
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
          <th class="px-4 py-2">Название</th>
          <th class="px-4 py-2">Slug</th>
          <th class="px-4 py-2">Видимость</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(objectType, index) in items"
          :key="objectType.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">{{ objectType.name }}</td>
          <td class="px-4 py-2">{{ objectType.slug }}</td>
          <td class="px-4 py-2">
            <AdminVisibilityDot
              :visible="objectType.visible"
              @click="onToggleVisible(objectType)"
            />
          </td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button
              class="text-[var(--color-primary)] hover:underline"
              @click="startEdit(objectType)"
            >
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(objectType.id)">
              Удалить
            </button>
          </td>
        </tr>
        <tr v-if="items.length === 0">
          <td colspan="5" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
