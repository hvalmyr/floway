<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

type DisplayStyle =
  | "blue-beige"
  | "brown-beige"
  | "beige-blue"
  | "beige-brown"
  | "blue-brown"
  | "brown-blue";

const displayStyleLabels: Record<DisplayStyle, string> = {
  "blue-beige": "Голубой фон, бежевый текст",
  "brown-beige": "Коричневый фон, бежевый текст",
  "beige-blue": "Бежевый фон, голубой текст",
  "beige-brown": "Бежевый фон, коричневый текст",
  "blue-brown": "Голубой фон, коричневый текст",
  "brown-blue": "Коричневый фон, голубой текст",
};

// Mirrors the 3 hex values in tailwind.config.ts (--color-primary/-surface/-ink)
// — only needed here to render the same preview for the fixed enum styles
// that AdminCourseCardPreview otherwise only gets for custom styles.
const displayStyleColors: Record<DisplayStyle, { bgColor: string; textColor: string }> = {
  "blue-beige": { bgColor: "#82B1CC", textColor: "#F7F5F3" },
  "brown-beige": { bgColor: "#41342A", textColor: "#F7F5F3" },
  "beige-blue": { bgColor: "#F7F5F3", textColor: "#82B1CC" },
  "beige-brown": { bgColor: "#F7F5F3", textColor: "#41342A" },
  "blue-brown": { bgColor: "#82B1CC", textColor: "#41342A" },
  "brown-blue": { bgColor: "#41342A", textColor: "#82B1CC" },
};

interface CourseBlock {
  id: number;
  courseId: number;
  blockName: string;
  description: string;
  blockCover: string;
  lessonCount: string;
  timeLength: string;
  price: string;
  displayStyle: DisplayStyle;
  sortOrder: number;
  visible: boolean;
  customDisplayStyleId: number | null;
}

interface CustomDisplayStyle {
  id: number;
  name: string;
  bgColor: string;
  textColor: string;
}

const route = useRoute();
const courseId = route.params.courseId as string;

const emptyForm = (): Omit<CourseBlock, "id" | "courseId"> => ({
  blockName: "",
  description: "",
  blockCover: "",
  lessonCount: "",
  timeLength: "",
  price: "",
  displayStyle: "blue-beige",
  sortOrder: 0,
  visible: true,
  customDisplayStyleId: null,
});

const { items, loading, error, fetchAll, create, update, remove } = useAdminResource<CourseBlock>(
  `/api/v1/courses/${courseId}/blocks`,
);

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

const customStyles = ref<CustomDisplayStyle[]>([]);
const previewColors = computed(() => {
  const custom = customStyles.value.find((s) => s.id === form.value.customDisplayStyleId);
  return custom ?? displayStyleColors[form.value.displayStyle];
});

// Single <select> combining the fixed enum styles and the admin-defined
// custom ones — form.displayStyle/customDisplayStyleId stay as two separate
// fields underneath (that's what the API expects), this just presents them
// as one control instead of two that could otherwise both apply at once.
const styleSelection = computed<string>({
  get: () =>
    form.value.customDisplayStyleId !== null
      ? `custom:${form.value.customDisplayStyleId}`
      : form.value.displayStyle,
  set: (value) => {
    if (value.startsWith("custom:")) {
      form.value.customDisplayStyleId = Number(value.slice("custom:".length));
    } else {
      form.value.displayStyle = value as DisplayStyle;
      form.value.customDisplayStyleId = null;
    }
  },
});

await fetchAll();
customStyles.value = await useApi().getCustomDisplayStyles();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(block: CourseBlock) {
  editingId.value = block.id;
  form.value = {
    blockName: block.blockName,
    description: block.description,
    blockCover: block.blockCover,
    lessonCount: block.lessonCount,
    timeLength: block.timeLength,
    price: block.price,
    displayStyle: block.displayStyle,
    sortOrder: block.sortOrder,
    visible: block.visible,
    customDisplayStyleId: block.customDisplayStyleId,
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
  if (!confirm("Удалить блок курса?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

async function onToggleVisible(block: CourseBlock) {
  await update(block.id, { ...block, visible: !block.visible });
}
</script>

<template>
  <div>
    <NuxtLink
      to="/admin/course-sections"
      class="text-sm text-[var(--color-primary)] hover:underline"
      >← К списку секций</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Блоки курса #{{ courseId }}</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Блоки — опция для курсов с несколькими программами (каждый блок — своя карточка на сайте).
      Если у курса всего одна программа, блоки не нужны вообще — заполните обложку/цену/занятия
      прямо в карточке курса на предыдущей странице.
    </p>

    <div class="mt-6 grid gap-6 lg:grid-cols-[1fr_auto] lg:items-start">
      <form
        class="grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
        @submit.prevent="onSubmit"
      >
        <input
          v-model="form.blockName"
          type="text"
          placeholder="Название блока (например, «Букеты»)"
          class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
        />
        <AdminMarkdownField
          v-model="form.description"
          placeholder="Вступительный текст над учебным планом блока"
          :rows="3"
          class="sm:col-span-2"
        />
        <AdminImageUpload v-model="form.blockCover" label="Обложка блока" />
        <input
          v-model="form.lessonCount"
          type="text"
          placeholder="Количество занятий (например, «7 занятий»)"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <input
          v-model="form.timeLength"
          type="text"
          placeholder="Продолжительность (например, «30 часов»)"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <input
          v-model="form.price"
          type="text"
          placeholder="Цена (например, «38 500 ₽»)"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <select
          v-model="styleSelection"
          class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
        >
          <optgroup label="Стандартные стили">
            <option v-for="(label, value) in displayStyleLabels" :key="value" :value="value">
              {{ label }}
            </option>
          </optgroup>
          <optgroup v-if="customStyles.length" label="Особые стили">
            <option v-for="style in customStyles" :key="style.id" :value="`custom:${style.id}`">
              {{ style.name }}
            </option>
          </optgroup>
        </select>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="form.visible" type="checkbox" class="size-5" />
          Показывать на сайте
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

      <div class="flex flex-col items-center gap-2 justify-self-center">
        <p class="text-sm text-[var(--color-text-muted)]">Так будет выглядеть карточка блока</p>
        <AdminCourseCardPreview
          :bg-color="previewColors.bgColor"
          :text-color="previewColors.textColor"
          :title="form.blockName || 'Название курса'"
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
          <th class="px-4 py-2">Название</th>
          <th class="px-4 py-2">Занятия</th>
          <th class="px-4 py-2">Часы</th>
          <th class="px-4 py-2">Цена</th>
          <th class="px-4 py-2">Стиль</th>
          <th class="px-4 py-2">Видимость</th>
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
          <td class="px-4 py-2">{{ block.blockName || "—" }}</td>
          <td class="px-4 py-2">{{ block.lessonCount }}</td>
          <td class="px-4 py-2">{{ block.timeLength }}</td>
          <td class="px-4 py-2">{{ block.price }}</td>
          <td class="px-4 py-2">
            {{
              customStyles.find((s) => s.id === block.customDisplayStyleId)?.name ??
              displayStyleLabels[block.displayStyle]
            }}
          </td>
          <td class="px-4 py-2">
            <AdminVisibilityDot :visible="block.visible" @click="onToggleVisible(block)" />
          </td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <NuxtLink
              :to="`/admin/course-blocks/${block.id}/lessons`"
              class="text-[var(--color-primary)] hover:underline"
              >Занятия</NuxtLink
            >
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(block)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(block.id)">
              Удалить
            </button>
          </td>
        </tr>
        <tr v-if="items.length === 0">
          <td colspan="8" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
