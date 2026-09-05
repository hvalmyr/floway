<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

type DisplayStyle = "blue-beige" | "brown-beige" | "beige-blue" | "beige-brown";

const displayStyleLabels: Record<DisplayStyle, string> = {
  "blue-beige": "Голубой фон, бежевый текст",
  "brown-beige": "Коричневый фон, бежевый текст",
  "beige-blue": "Бежевый фон, голубой текст",
  "beige-brown": "Бежевый фон, коричневый текст",
};

interface Course {
  id: number;
  sectionId: number;
  slug: string;
  name: string;
  description: string;
  coverImage: string;
  lessonCount: string;
  timeLength: string;
  price: string;
  displayStyle: DisplayStyle;
  sortOrder: number;
  visible: boolean;
  singleCard: boolean;
  faqTitle: string;
  faqDescription: string;
  faqVisible: boolean;
}

const route = useRoute();
const sectionId = route.params.sectionId as string;

const emptyForm = (): Omit<Course, "id" | "sectionId"> => ({
  slug: "",
  name: "",
  description: "",
  coverImage: "",
  lessonCount: "",
  timeLength: "",
  price: "",
  displayStyle: "blue-beige",
  sortOrder: 0,
  visible: true,
  singleCard: false,
  faqTitle: "",
  faqDescription: "",
  faqVisible: false,
});

const { items, loading, error, fetchAll, create, update, remove } = useAdminResource<Course>(
  `/api/v1/course-sections/${sectionId}/courses`,
);

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

// Auto-fills slug from name for a brand-new course, right up until the
// user types into the slug field themselves — editing an existing
// course's name never touches its (possibly already-live) slug.
const slugTouched = ref(false);
watch(
  () => form.value.name,
  (name) => {
    if (!slugTouched.value) form.value.slug = slugify(name);
  },
);

function startEdit(course: Course) {
  slugTouched.value = true;
  editingId.value = course.id;
  form.value = {
    slug: course.slug,
    name: course.name,
    description: course.description,
    coverImage: course.coverImage,
    lessonCount: course.lessonCount,
    timeLength: course.timeLength,
    price: course.price,
    displayStyle: course.displayStyle,
    sortOrder: course.sortOrder,
    visible: course.visible,
    singleCard: course.singleCard,
    faqTitle: course.faqTitle,
    faqDescription: course.faqDescription,
    faqVisible: course.faqVisible,
  };
}

function cancelEdit() {
  editingId.value = null;
  form.value = emptyForm();
  slugTouched.value = false;
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
  if (!confirm("Удалить курс вместе со всеми его блоками?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

async function onToggleVisible(course: Course) {
  await update(course.id, { ...course, visible: !course.visible });
}

async function onDuplicate(course: Course) {
  await create({
    slug: `${course.slug}-copy-${Date.now()}`,
    name: `${course.name} (копия)`,
    description: course.description,
    coverImage: course.coverImage,
    lessonCount: course.lessonCount,
    timeLength: course.timeLength,
    price: course.price,
    displayStyle: course.displayStyle,
    sortOrder: items.value.length,
    visible: false,
    singleCard: course.singleCard,
  });
}

// Drag reorder needs the row's index in the unfiltered `items` list, so
// it's only enabled when nothing is filtered out (filteredItems === items,
// same order) — searching and dragging at once would silently reorder the
// wrong rows.
const searchQuery = ref("");
const filteredItems = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return items.value;
  return items.value.filter((course) => course.name.toLowerCase().includes(query));
});
const dragEnabled = computed(() => searchQuery.value.trim() === "");

const { selectedIds, isSelected, toggle, allSelected, toggleAll, clear } =
  useAdminBulkSelect(filteredItems);

async function onBulkDelete() {
  if (!confirm(`Удалить выбранные курсы вместе со всеми их блоками (${selectedIds.value.size})?`))
    return;
  await Promise.all([...selectedIds.value].map((id) => remove(id)));
  clear();
}
</script>

<template>
  <div>
    <NuxtLink
      to="/admin/course-sections"
      class="text-sm text-[var(--color-primary)] hover:underline"
      >← К списку секций</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Курсы секции #{{ sectionId }}</h1>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.name"
        type="text"
        placeholder="Название курса"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model="form.slug"
        type="text"
        placeholder="Slug"
        required
        class="rounded border border-gray-300 px-3 py-2"
        @input="slugTouched = true"
      />
      <AdminMarkdownField
        v-model="form.description"
        placeholder="Описание"
        :rows="3"
        class="sm:col-span-2"
      />

      <div class="sm:col-span-2">
        <p class="text-sm font-medium">Карточка курса на сайте</p>
        <p class="text-sm text-[var(--color-text-muted)]">
          Используется, если у курса нет блоков, или если включена опция «Одна карточка» ниже.
          Иначе, если у курса есть блоки (вкладка «Блоки»), вместо неё показывается по одной
          карточке на каждый блок.
        </p>
      </div>
      <AdminImageUpload v-model="form.coverImage" label="Обложка курса" />
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
      <select v-model="form.displayStyle" class="rounded border border-gray-300 px-3 py-2">
        <option v-for="(label, value) in displayStyleLabels" :key="value" :value="value">
          {{ label }}
        </option>
      </select>
      <label class="flex items-center gap-2 text-sm">
        <input v-model="form.visible" type="checkbox" class="size-5" />
        Показывать на сайте
      </label>
      <label class="flex items-center gap-2 text-sm">
        <input v-model="form.singleCard" type="checkbox" class="size-5" />
        Одна карточка (даже если есть несколько блоков)
      </label>

      <div class="sm:col-span-2">
        <p class="text-sm font-medium">FAQ курса</p>
        <p class="text-sm text-[var(--color-text-muted)]">
          Показывается на странице курса после формы заявки. Сами вопросы и ответы редактируются на
          отдельной странице (ссылка «FAQ» в списке курсов ниже).
        </p>
      </div>
      <input
        v-model="form.faqTitle"
        type="text"
        placeholder="Заголовок FAQ (например, «Вопросы и ответы»)"
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
        Показывать FAQ на странице курса
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

    <div class="mt-6 flex flex-wrap items-center gap-3">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Поиск по названию…"
        class="min-w-64 rounded border border-gray-300 px-3 py-2 text-sm"
      />
      <span v-if="!dragEnabled" class="text-xs text-[var(--color-text-muted)]"
        >Перетаскивание для сортировки отключено во время поиска</span
      >
    </div>

    <div
      v-if="selectedIds.size > 0"
      class="mt-3 flex flex-wrap items-center gap-3 rounded border border-[var(--color-primary)] bg-white p-3 text-sm"
    >
      <span>Выбрано: {{ selectedIds.size }}</span>
      <button class="text-red-600 hover:underline" @click="onBulkDelete">Удалить выбранные</button>
      <button class="text-[var(--color-text-muted)] hover:underline" @click="clear">
        Отменить выбор
      </button>
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
          <th class="px-4 py-2">
            <input type="checkbox" class="size-5" :checked="allSelected" @change="toggleAll" />
          </th>
          <th class="px-4 py-2">Название</th>
          <th class="px-4 py-2">Slug</th>
          <th class="px-4 py-2">Видимость</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(course, index) in filteredItems"
          :key="course.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle v-if="dragEnabled" @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">
            <input
              type="checkbox"
              class="size-5"
              :checked="isSelected(course.id)"
              @change="toggle(course.id)"
            />
          </td>
          <td class="px-4 py-2">{{ course.name }}</td>
          <td class="px-4 py-2">{{ course.slug }}</td>
          <td class="px-4 py-2">
            <AdminVisibilityDot :visible="course.visible" @click="onToggleVisible(course)" />
          </td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <NuxtLink
              :to="`/admin/courses/${course.id}/blocks`"
              class="text-[var(--color-primary)] hover:underline"
              >Блоки</NuxtLink
            >
            <NuxtLink
              :to="`/admin/courses/${course.id}/lessons`"
              class="text-[var(--color-primary)] hover:underline"
              >Занятия</NuxtLink
            >
            <NuxtLink
              :to="`/admin/courses/${course.id}/faq`"
              class="text-[var(--color-primary)] hover:underline"
              >FAQ</NuxtLink
            >
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(course)">
              Редактировать
            </button>
            <button
              class="text-[var(--color-primary)] hover:underline"
              @click="onDuplicate(course)"
            >
              Дублировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(course.id)">
              Удалить
            </button>
          </td>
        </tr>
        <tr v-if="filteredItems.length === 0">
          <td colspan="6" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            {{ items.length === 0 ? "Пока пусто" : "Ничего не найдено" }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
