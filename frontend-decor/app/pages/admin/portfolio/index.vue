<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Portfolio — decor's repurposing of the school's gallery_photos table
// (docs/decor-site-audit.md §2: "репрофилируется"), tagged with object type
// and format so the public /portfolio page can filter (п. 6 ТЗ). Not built
// on AdminPhotoCarouselManager (the school's generic carousel-manager
// component) since that has no room for these two extra fields.
interface GalleryPhoto {
  id: number;
  image: string;
  objectTypeId: number | null;
  format: "" | "season" | "event";
  sortOrder: number;
}

interface ObjectType {
  id: number;
  name: string;
}

const emptyForm = (): Omit<GalleryPhoto, "id"> => ({
  image: "",
  objectTypeId: null,
  format: "",
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<GalleryPhoto>("/api/v1/gallery-photos");
const { items: objectTypes, fetchAll: fetchObjectTypes } =
  useAdminResource<ObjectType>("/api/v1/object-types");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await Promise.all([fetchAll(), fetchObjectTypes()]);

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function objectTypeName(id: number | null) {
  if (!id) return "—";
  return objectTypes.value.find((t) => t.id === id)?.name ?? `#${id}`;
}

const formatLabels: Record<string, string> = { "": "—", season: "На сезон", event: "На праздник" };

function startEdit(photo: GalleryPhoto) {
  editingId.value = photo.id;
  form.value = {
    image: photo.image,
    objectTypeId: photo.objectTypeId,
    format: photo.format,
    sortOrder: photo.sortOrder,
  };
}

function cancelEdit() {
  editingId.value = null;
  form.value = emptyForm();
}

async function onSubmit() {
  formError.value = "";
  if (!form.value.image) {
    formError.value = "Загрузите фото.";
    return;
  }
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
  if (!confirm("Удалить фото из портфолио?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Портфолио</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Фото примеров работ — тип объекта и формат используются фильтром на публичной странице
      «Портфолио».
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <AdminImageUpload v-model="form.image" label="Фото *" class="sm:col-span-2" />
      <select v-model="form.objectTypeId" class="rounded border border-gray-300 px-3 py-2">
        <option :value="null">Тип объекта — не указан</option>
        <option v-for="type in objectTypes" :key="type.id" :value="type.id">{{ type.name }}</option>
      </select>
      <select v-model="form.format" class="rounded border border-gray-300 px-3 py-2">
        <option value="">Формат — не указан</option>
        <option value="season">На сезон</option>
        <option value="event">На праздник</option>
      </select>

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
          <th class="px-4 py-2">Фото</th>
          <th class="px-4 py-2">Тип объекта</th>
          <th class="px-4 py-2">Формат</th>
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
              v-if="photo.image"
              :src="resolveMediaUrl(photo.image)"
              class="h-12 w-12 rounded object-cover"
              alt=""
            />
          </td>
          <td class="px-4 py-2">{{ objectTypeName(photo.objectTypeId) }}</td>
          <td class="px-4 py-2">{{ formatLabels[photo.format] }}</td>
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
          <td colspan="5" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
