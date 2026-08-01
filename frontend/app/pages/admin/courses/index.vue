<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface Course {
  id: number;
  slug: string;
  title: string;
  shortDescription: string;
  fullDescription: string;
  status: "active" | "archived";
  coverImage: string;
  gallery: string[];
  sortOrder: number;
}

const statusLabels: Record<Course["status"], string> = {
  active: "Активен",
  archived: "В архиве",
};

const emptyForm = () => ({
  slug: "",
  title: "",
  shortDescription: "",
  fullDescription: "",
  status: "active" as Course["status"],
  coverImage: "",
  gallery: [] as string[],
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<Course>("/api/v1/courses");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

function startEdit(course: Course) {
  editingId.value = course.id;
  form.value = {
    slug: course.slug,
    title: course.title,
    shortDescription: course.shortDescription,
    fullDescription: course.fullDescription,
    status: course.status,
    coverImage: course.coverImage,
    gallery: [...course.gallery],
    sortOrder: course.sortOrder,
  };
}

function addGalleryImage(url: string) {
  form.value.gallery = [...form.value.gallery, url];
}

function removeGalleryImage(index: number) {
  form.value.gallery = form.value.gallery.filter((_, i) => i !== index);
}

function cancelEdit() {
  editingId.value = null;
  form.value = emptyForm();
}

async function onSubmit() {
  formError.value = "";
  saving.value = true;
  try {
    const payload = {
      slug: form.value.slug,
      title: form.value.title,
      shortDescription: form.value.shortDescription,
      fullDescription: form.value.fullDescription,
      status: form.value.status,
      coverImage: form.value.coverImage,
      gallery: form.value.gallery.filter(Boolean),
      sortOrder: form.value.sortOrder,
    };
    if (editingId.value === null) {
      await create(payload);
    } else {
      await update(editingId.value, payload);
    }
    cancelEdit();
  } catch {
    formError.value = "Не удалось сохранить. Проверьте поля и попробуйте снова.";
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Удалить курс?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Курсы</h1>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.title"
        type="text"
        placeholder="Заголовок"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model="form.slug"
        type="text"
        placeholder="Slug"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <textarea
        v-model="form.shortDescription"
        placeholder="Краткое описание"
        rows="2"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <textarea
        v-model="form.fullDescription"
        placeholder="Полное описание"
        rows="4"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <AdminImageUpload v-model="form.coverImage" label="Обложка" />

      <div class="flex flex-col gap-2 sm:col-span-2">
        <span class="text-sm text-[var(--color-text-muted)]">Галерея</span>
        <div
          v-for="(url, index) in form.gallery"
          :key="index"
          class="flex items-center gap-3 rounded border border-gray-300 px-3 py-2"
        >
          <img :src="resolveMediaUrl(url)" alt="" class="size-12 shrink-0 rounded object-cover" />
          <span class="min-w-0 flex-1 truncate text-sm text-[var(--color-text-muted)]">{{
            url
          }}</span>
          <button
            type="button"
            class="shrink-0 rounded border border-gray-300 px-2 py-1 text-sm text-red-600"
            @click="removeGalleryImage(index)"
          >
            Убрать
          </button>
        </div>
        <AdminImageUpload
          model-value=""
          label="Добавить фото в галерею"
          @update:model-value="addGalleryImage"
        />
      </div>

      <select v-model="form.status" class="rounded border border-gray-300 px-3 py-2">
        <option value="active">Активен</option>
        <option value="archived">В архиве</option>
      </select>
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
          <th class="px-4 py-2">Slug</th>
          <th class="px-4 py-2">Статус</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr v-for="course in items" :key="course.id" class="border-b border-gray-100 last:border-0">
          <td class="px-4 py-2">{{ course.title }}</td>
          <td class="px-4 py-2">{{ course.slug }}</td>
          <td class="px-4 py-2">{{ statusLabels[course.status] }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <NuxtLink
              :to="`/admin/courses/${course.id}/blocks`"
              class="text-[var(--color-primary)] hover:underline"
              >Блоки курса</NuxtLink
            >
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(course)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(course.id)">
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
