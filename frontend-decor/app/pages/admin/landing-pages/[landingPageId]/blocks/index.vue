<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface LandingPageBlock {
  id: number;
  landingPageId: number;
  image: string;
  text: string;
  sortOrder: number;
}

const route = useRoute();
const landingPageId = route.params.landingPageId as string;

const emptyForm = (): Omit<LandingPageBlock, "id" | "landingPageId"> => ({
  image: "",
  text: "",
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<LandingPageBlock>(`/api/v1/landing-pages/${landingPageId}/blocks`);

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(block: LandingPageBlock) {
  editingId.value = block.id;
  form.value = { image: block.image, text: block.text, sortOrder: block.sortOrder };
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
  if (!confirm("Удалить блок?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <NuxtLink to="/admin/landing-pages" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К списку страниц</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Блоки страницы #{{ landingPageId }}</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Фото + текст — примеры работ и форматы услуги на странице, в порядке снизу.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <AdminImageUpload v-model="form.image" label="Фото" class="sm:col-span-2" />
      <AdminMarkdownField
        v-model="form.text"
        placeholder="Текст блока"
        :rows="3"
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
          <th class="px-4 py-2">Фото</th>
          <th class="px-4 py-2">Текст</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(block, index) in items"
          :key="block.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">
            <img
              v-if="block.image"
              :src="block.image"
              class="h-12 w-12 rounded object-cover"
              alt=""
            />
            <span v-else class="text-[var(--color-text-muted)]">—</span>
          </td>
          <td class="max-w-md truncate px-4 py-2">{{ block.text }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
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
