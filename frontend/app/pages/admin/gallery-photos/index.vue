<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface GalleryPhoto {
  id: number;
  image: string;
  sortOrder: number;
}

const emptyForm = (): Omit<GalleryPhoto, "id"> => ({ image: "", sortOrder: 0 });

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<GalleryPhoto>("/api/v1/gallery-photos");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(photo: GalleryPhoto) {
  editingId.value = photo.id;
  form.value = { image: photo.image, sortOrder: photo.sortOrder };
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
  if (!confirm("Удалить фото из карусели?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Фотогалерея (карусель на главной)</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Вертикальные фото, показанные в слайд-шоу между блоками «Преимущества» и «О школе». Порядок
      сортировки определяет порядок слайдов.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <AdminImageUpload v-model="form.image" label="Фото" />

      <p v-if="formError" class="text-sm text-red-600 sm:col-span-2">{{ formError }}</p>

      <div class="flex gap-2 sm:col-span-2">
        <button
          type="submit"
          :disabled="saving || !form.image"
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
          <th class="px-4 py-2">Фото</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(photo, index) in items"
          :key="photo.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">
            <img
              :src="resolveMediaUrl(photo.image)"
              alt=""
              class="h-16 w-12 rounded object-cover"
            />
          </td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(photo)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(photo.id)">
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
