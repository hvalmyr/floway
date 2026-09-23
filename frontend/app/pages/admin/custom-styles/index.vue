<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface CustomDisplayStyle {
  id: number;
  name: string;
  bgColor: string;
  textColor: string;
  sortOrder: number;
}

const emptyForm = (): Omit<CustomDisplayStyle, "id"> => ({
  name: "",
  bgColor: "#82b1cc",
  textColor: "#41342a",
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<CustomDisplayStyle>("/api/v1/custom-display-styles");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(style: CustomDisplayStyle) {
  editingId.value = style.id;
  form.value = {
    name: style.name,
    bgColor: style.bgColor,
    textColor: style.textColor,
    sortOrder: style.sortOrder,
  };
}

function cancelEdit() {
  editingId.value = null;
  form.value = emptyForm();
}

function swapColors() {
  const { bgColor, textColor } = form.value;
  form.value.bgColor = textColor;
  form.value.textColor = bgColor;
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
    formError.value = "Не удалось сохранить. Проверьте, что оба цвета выбраны.";
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Удалить стиль? Курсы и блоки, использующие его, вернутся к обычному стилю.")) {
    return;
  }
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Особые стили карточек</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Свои сочетания цвета фона и текста для карточек курсов — например, новогодний или осенний
      стиль. После создания стиль можно выбрать в настройках курса или блока курса вместо
      стандартного набора цветов.
    </p>

    <div class="mt-6 grid gap-6 lg:grid-cols-[1fr_auto] lg:items-start">
      <form
        class="grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
        @submit.prevent="onSubmit"
      >
        <input
          v-model="form.name"
          type="text"
          placeholder="Название (например, «Новый год»)"
          required
          class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
        />

        <label class="flex items-center gap-3 rounded border border-gray-300 px-3 py-2">
          <span class="text-sm">Фон</span>
          <input
            v-model="form.bgColor"
            type="color"
            required
            class="h-8 w-8 cursor-pointer rounded border-0"
          />
          <span class="text-sm text-[var(--color-text-muted)]">{{ form.bgColor }}</span>
        </label>
        <label class="flex items-center gap-3 rounded border border-gray-300 px-3 py-2">
          <span class="text-sm">Текст</span>
          <input
            v-model="form.textColor"
            type="color"
            required
            class="h-8 w-8 cursor-pointer rounded border-0"
          />
          <span class="text-sm text-[var(--color-text-muted)]">{{ form.textColor }}</span>
        </label>

        <button
          type="button"
          title="Поменять фон и текст местами"
          class="flex items-center justify-center gap-2 rounded border border-gray-300 px-3 py-2 text-sm hover:bg-gray-50 sm:col-span-2"
          @click="swapColors"
        >
          ⇄ Поменять цвета местами
        </button>

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

      <div class="flex flex-col items-center gap-2 justify-self-center">
        <p class="text-sm text-[var(--color-text-muted)]">Так будет выглядеть карточка курса</p>
        <AdminCourseCardPreview
          :bg-color="form.bgColor"
          :text-color="form.textColor"
          :title="form.name || 'Название курса'"
        />
      </div>
    </div>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <table
      v-else
      class="mt-6 w-full border-collapse overflow-hidden rounded border border-gray-200 bg-white text-sm"
    >
      <thead>
        <tr class="border-b border-gray-200 text-left">
          <th class="px-4 py-2" />
          <th class="px-4 py-2">Превью</th>
          <th class="px-4 py-2">Название</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(style, index) in items"
          :key="style.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">
            <span
              class="inline-flex size-8 items-center justify-center rounded-sm text-xs font-bold"
              :style="{ backgroundColor: style.bgColor, color: style.textColor }"
              >Aa</span
            >
          </td>
          <td class="px-4 py-2">{{ style.name }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(style)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(style.id)">
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
