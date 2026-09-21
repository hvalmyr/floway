<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Посадочные страницы под тип объекта (п. 7.5 ТЗ): create/copy/rename/
// hide/delete из админки, без деплоя. Одна страница = один тип объекта.
interface LandingPage {
  id: number;
  objectTypeId: number;
  slug: string;
  h1: string;
  metaTitle: string;
  metaDescription: string;
  faqTitle: string;
  faqDescription: string;
  faqVisible: boolean;
  visible: boolean;
  sortOrder: number;
}

interface ObjectType {
  id: number;
  slug: string;
  name: string;
  visible: boolean;
}

const emptyForm = (): Omit<LandingPage, "id"> => ({
  objectTypeId: 0,
  slug: "",
  h1: "",
  metaTitle: "",
  metaDescription: "",
  faqTitle: "",
  faqDescription: "",
  faqVisible: false,
  visible: true,
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<LandingPage>("/api/v1/landing-pages");
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

function objectTypeName(id: number) {
  return objectTypes.value.find((t) => t.id === id)?.name ?? `#${id}`;
}

const slugTouched = ref(false);
watch(
  () => form.value.h1,
  (h1) => {
    if (!slugTouched.value) form.value.slug = slugify(h1);
  },
);

function startEdit(page: LandingPage) {
  slugTouched.value = true;
  editingId.value = page.id;
  form.value = {
    objectTypeId: page.objectTypeId,
    slug: page.slug,
    h1: page.h1,
    metaTitle: page.metaTitle,
    metaDescription: page.metaDescription,
    faqTitle: page.faqTitle,
    faqDescription: page.faqDescription,
    faqVisible: page.faqVisible,
    visible: page.visible,
    sortOrder: page.sortOrder,
  };
}

function cancelEdit() {
  editingId.value = null;
  form.value = emptyForm();
  slugTouched.value = false;
}

async function onSubmit() {
  formError.value = "";
  if (!form.value.objectTypeId) {
    formError.value = "Выберите тип объекта.";
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
  if (!confirm("Удалить посадочную страницу вместе со всеми блоками и FAQ?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

async function onToggleVisible(page: LandingPage) {
  await update(page.id, { ...page, visible: !page.visible });
}

async function onDuplicate(page: LandingPage) {
  await create({
    objectTypeId: page.objectTypeId,
    slug: `${page.slug}-copy-${Date.now()}`,
    h1: `${page.h1} (копия)`,
    metaTitle: page.metaTitle,
    metaDescription: page.metaDescription,
    faqTitle: page.faqTitle,
    faqDescription: page.faqDescription,
    faqVisible: page.faqVisible,
    sortOrder: items.value.length,
    visible: false,
  });
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Посадочные страницы</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Блоки страницы (фото + текст) и FAQ редактируются на отдельных страницах — ссылки «Блоки» и
      «FAQ» в списке ниже.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <select v-model="form.objectTypeId" required class="rounded border border-gray-300 px-3 py-2">
        <option :value="0" disabled>Тип объекта *</option>
        <option v-for="type in objectTypes" :key="type.id" :value="type.id">{{ type.name }}</option>
      </select>
      <input
        v-model="form.slug"
        type="text"
        placeholder="Slug (адрес страницы) *"
        required
        class="rounded border border-gray-300 px-3 py-2"
        @input="slugTouched = true"
      />
      <input
        v-model="form.h1"
        type="text"
        placeholder="Заголовок H1 *"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <input
        v-model="form.metaTitle"
        type="text"
        placeholder="Meta title (для поисковиков)"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model="form.metaDescription"
        type="text"
        placeholder="Meta description"
        class="rounded border border-gray-300 px-3 py-2"
      />

      <div class="sm:col-span-2">
        <p class="text-sm font-medium">FAQ страницы</p>
        <p class="text-sm text-[var(--color-text-muted)]">
          Вопросы и ответы редактируются на отдельной странице (ссылка «FAQ» в списке ниже).
        </p>
      </div>
      <input
        v-model="form.faqTitle"
        type="text"
        placeholder="Заголовок FAQ"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <AdminMarkdownField
        v-model="form.faqDescription"
        placeholder="Текст перед вопросами"
        :rows="2"
        class="sm:col-span-2"
      />
      <label class="flex items-center gap-2 text-sm">
        <input v-model="form.faqVisible" type="checkbox" class="size-5" />
        Показывать FAQ на странице
      </label>
      <label class="flex items-center gap-2 text-sm">
        <input v-model="form.visible" type="checkbox" class="size-5" />
        Показывать страницу на сайте
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
          <th class="px-4 py-2">Тип объекта</th>
          <th class="px-4 py-2">Slug</th>
          <th class="px-4 py-2">Видимость</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(page, index) in items"
          :key="page.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">{{ page.h1 }}</td>
          <td class="px-4 py-2">{{ objectTypeName(page.objectTypeId) }}</td>
          <td class="px-4 py-2">{{ page.slug }}</td>
          <td class="px-4 py-2">
            <AdminVisibilityDot :visible="page.visible" @click="onToggleVisible(page)" />
          </td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <NuxtLink
              :to="`/admin/landing-pages/${page.id}/blocks`"
              class="text-[var(--color-primary)] hover:underline"
              >Блоки</NuxtLink
            >
            <NuxtLink
              :to="`/admin/landing-pages/${page.id}/faq`"
              class="text-[var(--color-primary)] hover:underline"
              >FAQ</NuxtLink
            >
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(page)">
              Редактировать
            </button>
            <button class="text-[var(--color-primary)] hover:underline" @click="onDuplicate(page)">
              Дублировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(page.id)">Удалить</button>
          </td>
        </tr>
        <tr v-if="items.length === 0">
          <td colspan="6" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
