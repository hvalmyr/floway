<script setup lang="ts">
import { FEATURE_ICONS } from "~/constants/feature-icons";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Decor site only renders FeatureGrid on the home page ("Почему выбирают
// нас") — no masterclasses/gift_certificate blocks, so this is the
// single-page-block case of the school's features/index.vue (no tabs).
interface FeatureEntity {
  id: number;
  page: "home";
  icon: string;
  title: string;
  description: string;
  sortOrder: number;
}

const iconOptions = Object.entries(FEATURE_ICONS).map(([value, { label }]) => ({ value, label }));

const emptyForm = () => ({
  icon: iconOptions[0]?.value ?? "",
  title: "",
  description: "",
});

// No `?page=` query — decor never has masterclasses/gift_certificate rows,
// so the unfiltered list is already just this site's "home" features (and
// unlike a query-string basePath, useAdminResource's update/remove can
// safely append `/${id}` to a plain path).
const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<FeatureEntity>("/api/v1/features");

const editingId = ref<number | null>(null);
const editingSortOrder = ref(0);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

const {
  items: headingContentItems,
  loading: headingContentLoading,
  savingKey: headingContentSavingKey,
  savedKey: headingContentSavedKey,
  fetchAll: fetchHeadingContent,
  save: saveHeadingContent,
} = useAdminPageContent(["home_features_heading", "home_features_lead"]);

await fetchHeadingContent();

function startEdit(item: FeatureEntity) {
  editingId.value = item.id;
  editingSortOrder.value = item.sortOrder;
  form.value = { icon: item.icon, title: item.title, description: item.description };
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
      await create({ ...form.value, page: "home", sortOrder: items.value.length });
    } else {
      await update(editingId.value, {
        ...form.value,
        page: "home",
        sortOrder: editingSortOrder.value,
      });
    }
    cancelEdit();
  } catch {
    formError.value = "Не удалось сохранить. Проверьте поля и попробуйте снова.";
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Удалить карточку преимущества?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Преимущества (главная)</h1>

    <p v-if="loading || headingContentLoading" class="mt-6 text-[var(--color-text-muted)]">
      Загрузка…
    </p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 flex flex-col gap-4">
      <div
        v-for="item in headingContentItems"
        :key="item.key"
        class="rounded border border-gray-200 bg-white p-4"
      >
        <div class="mb-2 flex items-center justify-between gap-4">
          <label :for="`field-${item.key}`" class="text-sm font-medium">{{
            item.key.endsWith("heading") ? "Заголовок" : "Описание"
          }}</label>
          <span class="font-mono text-xs text-[var(--color-text-muted)]">{{ item.key }}</span>
        </div>
        <input
          :id="`field-${item.key}`"
          v-model="item.value"
          type="text"
          class="w-full rounded border border-gray-300 px-3 py-2"
        />
        <div class="mt-2 flex items-center gap-3">
          <button
            type="button"
            :disabled="headingContentSavingKey === item.key"
            class="rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white disabled:opacity-50"
            @click="saveHeadingContent(item)"
          >
            {{ headingContentSavingKey === item.key ? "Сохранение…" : "Сохранить" }}
          </button>
          <span v-if="headingContentSavedKey === item.key" class="text-sm text-green-600"
            >Сохранено</span
          >
        </div>
      </div>

      <form
        class="grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
        @submit.prevent="onSubmit"
      >
        <AdminIconPicker v-model="form.icon" class="sm:col-span-2" />
        <input
          v-model="form.title"
          type="text"
          placeholder="Заголовок"
          required
          class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
        />
        <AdminMarkdownField
          v-model="form.description"
          placeholder="Описание"
          :rows="3"
          required
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

      <table
        class="w-full border-collapse overflow-hidden rounded border border-gray-200 bg-white text-sm"
      >
        <thead>
          <tr class="border-b border-gray-200 text-left">
            <th class="px-4 py-2" />
            <th class="px-4 py-2">Иконка</th>
            <th class="px-4 py-2">Заголовок</th>
            <th class="px-4 py-2" />
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(item, index) in items"
            :key="item.id"
            :data-row-index="index"
            class="border-b border-gray-100 last:border-0"
            :class="draggingIndex === index ? 'opacity-50' : ''"
          >
            <td class="px-4 py-2">
              <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
            </td>
            <td class="px-4 py-2">
              <AppIcon :icon="item.icon" class="size-6 text-[var(--color-primary)]" />
            </td>
            <td class="px-4 py-2">{{ item.title }}</td>
            <td class="flex gap-3 px-4 py-2 text-right">
              <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(item)">
                Редактировать
              </button>
              <button class="text-red-600 hover:underline" @click="onDelete(item.id)">
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
  </div>
</template>
