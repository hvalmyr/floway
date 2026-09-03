<script setup lang="ts">
import { FEATURE_ICONS } from "~/constants/feature-icons";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface FeatureEntity {
  id: number;
  page: "home" | "masterclasses";
  icon: string;
  title: string;
  description: string;
  sortOrder: number;
}

const pageLabels: Record<FeatureEntity["page"], string> = {
  home: "Главная",
  masterclasses: "Мастер-классы",
};

const iconOptions = Object.entries(FEATURE_ICONS).map(([value, { label }]) => ({ value, label }));

const emptyForm = () => ({
  page: "home" as FeatureEntity["page"],
  icon: iconOptions[0]?.value ?? "",
  title: "",
  description: "",
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<FeatureEntity>("/api/v1/features");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

// Each page's block heading/lead ("Почему стоит учиться"/"Почему стоит
// выбрать мастер-класс...") lives here, right above that page's item list,
// rather than mixed into the Hero admin page — it's the intro text for
// *this* block, not the page's hero.
const {
  items: blockContentItems,
  loading: blockContentLoading,
  savingKey: blockContentSavingKey,
  savedKey: blockContentSavedKey,
  fetchAll: fetchBlockContent,
  save: saveBlockContent,
} = useAdminPageContent([
  "home_features_heading",
  "home_features_lead",
  "masterclasses_features_heading",
  "masterclasses_features_lead",
]);

await fetchBlockContent();

function blockContentFor(page: FeatureEntity["page"]) {
  const prefix = page === "home" ? "home_features_" : "masterclasses_features_";
  return [
    blockContentItems.value.find((item) => item.key === `${prefix}heading`),
    blockContentItems.value.find((item) => item.key === `${prefix}lead`),
  ].filter((item): item is (typeof blockContentItems.value)[number] => item !== undefined);
}

// sort_order is scoped per `page` on the backend (see feature_repository.go's
// `ORDER BY page, sort_order, id`) — each page's features are numbered
// independently, so dragging must stay within one page's group. A writable
// computed reads/writes that page's slice of the shared `items` ref, leaving
// the other page's rows untouched.
function pageGroup(page: FeatureEntity["page"]) {
  return computed<FeatureEntity[]>({
    get: () =>
      items.value.filter((item) => item.page === page).sort((a, b) => a.sortOrder - b.sortOrder),
    set: (reordered) => {
      items.value = [...items.value.filter((item) => item.page !== page), ...reordered];
    },
  });
}

const homeItems = pageGroup("home");
const masterclassItems = pageGroup("masterclasses");
const homeReorder = useAdminDragReorder(homeItems, (item) => update(item.id, item));
const masterclassReorder = useAdminDragReorder(masterclassItems, (item) => update(item.id, item));

const featureGroups = [
  { page: "home" as const, items: homeItems, reorder: homeReorder },
  { page: "masterclasses" as const, items: masterclassItems, reorder: masterclassReorder },
];

function startEdit(item: FeatureEntity) {
  editingId.value = item.id;
  form.value = {
    page: item.page,
    icon: item.icon,
    title: item.title,
    description: item.description,
    sortOrder: item.sortOrder,
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
      const sortOrder = items.value.filter((item) => item.page === form.value.page).length;
      await create({ ...form.value, sortOrder });
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
  if (!confirm("Удалить карточку преимущества?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Карточки преимуществ</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Блок «Почему стоит учиться» на главной и блок преимуществ на странице мастер-классов.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <select v-model="form.page" class="rounded border border-gray-300 px-3 py-2">
        <option value="home">Главная</option>
        <option value="masterclasses">Мастер-классы</option>
      </select>
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

    <p v-if="loading || blockContentLoading" class="mt-6 text-[var(--color-text-muted)]">
      Загрузка…
    </p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <template v-else>
      <div v-for="group in featureGroups" :key="group.page" class="mt-6">
        <h2 class="mb-2 font-medium">{{ pageLabels[group.page] }}</h2>

        <div
          v-for="item in blockContentFor(group.page)"
          :key="item.key"
          class="mb-3 rounded border border-gray-200 bg-white p-4"
        >
          <div class="mb-2 flex items-center justify-between gap-4">
            <label :for="`field-${item.key}`" class="text-sm font-medium">{{ item.label }}</label>
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
              v-for="(item, index) in group.items.value"
              :key="item.id"
              draggable="true"
              class="border-b border-gray-100 last:border-0"
              :class="group.reorder.draggingIndex.value === index ? 'opacity-50' : ''"
              @dragstart="group.reorder.onDragStart(index)"
              @dragover.prevent="group.reorder.onDragOver(index)"
              @drop.prevent="group.reorder.onDrop"
            >
              <td class="px-4 py-2"><AdminDragHandle /></td>
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
            <tr v-if="group.items.value.length === 0">
              <td colspan="4" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
                Пока пусто
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>
