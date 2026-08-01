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
      await create(form.value);
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
      <select v-model="form.icon" class="rounded border border-gray-300 px-3 py-2">
        <option v-for="opt in iconOptions" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </option>
      </select>
      <input
        v-model="form.title"
        type="text"
        placeholder="Заголовок"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <textarea
        v-model="form.description"
        placeholder="Описание"
        rows="3"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
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
          <th class="px-4 py-2">Страница</th>
          <th class="px-4 py-2">Иконка</th>
          <th class="px-4 py-2">Заголовок</th>
          <th class="px-4 py-2">Порядок</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" class="border-b border-gray-100 last:border-0">
          <td class="px-4 py-2">{{ pageLabels[item.page] }}</td>
          <td class="px-4 py-2">{{ FEATURE_ICONS[item.icon]?.label ?? item.icon }}</td>
          <td class="px-4 py-2">{{ item.title }}</td>
          <td class="px-4 py-2">{{ item.sortOrder }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(item)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(item.id)">Удалить</button>
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
