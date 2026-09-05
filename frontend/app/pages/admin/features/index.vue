<script setup lang="ts">
import { FEATURE_ICONS } from "~/constants/feature-icons";
import type { FeaturePage } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface FeatureEntity {
  id: number;
  page: FeaturePage;
  icon: string;
  title: string;
  description: string;
  sortOrder: number;
}

/**
 * Every page that renders a FeatureGrid stores its heading/lead as
 * `<prefix>heading`/`<prefix>lead` page_content keys, plus its own slice of
 * `features` rows (scoped by `page`). One block per page — a list on the
 * left, one add/edit form + heading/lead fields + items table on the right
 * for whichever block is selected — same pattern as the Hero admin page,
 * instead of stacking every page's section one after another.
 */
const pageBlocks: { id: FeaturePage; label: string; prefix: string }[] = [
  { id: "home", label: "Главная", prefix: "home_features_" },
  { id: "masterclasses", label: "Мастер-классы", prefix: "masterclasses_features_" },
  { id: "gift_certificate", label: "Подарочные сертификаты", prefix: "gift_certificate_features_" },
];

const iconOptions = Object.entries(FEATURE_ICONS).map(([value, { label }]) => ({ value, label }));

const emptyForm = () => ({
  icon: iconOptions[0]?.value ?? "",
  title: "",
  description: "",
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<FeatureEntity>("/api/v1/features");

const editingId = ref<number | null>(null);
const editingSortOrder = ref(0);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const {
  items: blockContentItems,
  loading: blockContentLoading,
  savingKey: blockContentSavingKey,
  savedKey: blockContentSavedKey,
  fetchAll: fetchBlockContent,
  save: saveBlockContent,
} = useAdminPageContent(
  pageBlocks.flatMap((block) => [`${block.prefix}heading`, `${block.prefix}lead`]),
);

await fetchBlockContent();

function blockContentFor(prefix: string) {
  return [
    blockContentItems.value.find((item) => item.key === `${prefix}heading`),
    blockContentItems.value.find((item) => item.key === `${prefix}lead`),
  ].filter((item): item is (typeof blockContentItems.value)[number] => item !== undefined);
}

// sort_order is scoped per `page` on the backend (see feature_repository.go's
// `ORDER BY page, sort_order, id`) — each page's features are numbered
// independently, so dragging must stay within one page's group. A writable
// computed reads/writes that page's slice of the shared `items` ref, leaving
// the other pages' rows untouched.
function pageGroup(page: FeaturePage) {
  return computed<FeatureEntity[]>({
    get: () =>
      items.value.filter((item) => item.page === page).sort((a, b) => a.sortOrder - b.sortOrder),
    set: (reordered) => {
      items.value = [...items.value.filter((item) => item.page !== page), ...reordered];
    },
  });
}

const featureGroups = pageBlocks.map((block) => {
  const groupItems = pageGroup(block.id);
  return {
    ...block,
    items: groupItems,
    reorder: useAdminDragReorder(groupItems, (item) => update(item.id, item)),
  };
});

const selectedPageId = ref(pageBlocks[0]!.id);
const selectedGroup = computed(() =>
  featureGroups.find((group) => group.id === selectedPageId.value)!,
);

function selectPage(pageId: FeaturePage) {
  selectedPageId.value = pageId;
  cancelEdit();
}

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
      const sortOrder = selectedGroup.value.items.value.length;
      await create({ ...form.value, page: selectedPageId.value, sortOrder });
    } else {
      await update(editingId.value, {
        ...form.value,
        page: selectedPageId.value,
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
    <h1 class="text-2xl font-semibold">Карточки преимуществ</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Блок «Почему стоит учиться» / «Почему это отличный подарок» и т.п. — свой набор карточек для
      каждой страницы.
    </p>

    <p v-if="loading || blockContentLoading" class="mt-6 text-[var(--color-text-muted)]">
      Загрузка…
    </p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-[200px_1fr]">
      <div class="flex flex-row gap-2 sm:flex-col">
        <button
          v-for="block in pageBlocks"
          :key="block.id"
          type="button"
          class="rounded border px-3 py-2 text-left text-sm"
          :class="
            block.id === selectedPageId
              ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/10 font-medium'
              : 'border-gray-200 bg-white hover:border-[var(--color-primary)]'
          "
          @click="selectPage(block.id)"
        >
          {{ block.label }}
        </button>
      </div>

      <div class="flex flex-col gap-4">
        <div
          v-for="item in blockContentFor(selectedGroup.prefix)"
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
              :disabled="blockContentSavingKey === item.key"
              class="rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white disabled:opacity-50"
              @click="saveBlockContent(item)"
            >
              {{ blockContentSavingKey === item.key ? "Сохранение…" : "Сохранить" }}
            </button>
            <span v-if="blockContentSavedKey === item.key" class="text-sm text-green-600"
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
              v-for="(item, index) in selectedGroup.items.value"
              :key="item.id"
              :data-row-index="index"
              class="border-b border-gray-100 last:border-0"
              :class="selectedGroup.reorder.draggingIndex.value === index ? 'opacity-50' : ''"
            >
              <td class="px-4 py-2">
                <AdminDragHandle
                  @pointerdown="selectedGroup.reorder.onPointerDown(index, $event)"
                />
              </td>
              <td class="px-4 py-2">
                <AppIcon :icon="item.icon" class="size-6 text-[var(--color-primary)]" />
              </td>
              <td class="px-4 py-2">{{ item.title }}</td>
              <td class="flex gap-3 px-4 py-2 text-right">
                <button
                  class="text-[var(--color-primary)] hover:underline"
                  @click="startEdit(item)"
                >
                  Редактировать
                </button>
                <button class="text-red-600 hover:underline" @click="onDelete(item.id)">
                  Удалить
                </button>
              </td>
            </tr>
            <tr v-if="selectedGroup.items.value.length === 0">
              <td colspan="4" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
                Пока пусто
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
