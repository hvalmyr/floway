<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface Masterclass {
  id: number;
  slug: string;
  title: string;
  description: string;
  description2: string;
  endingText: string;
  duration: string;
  price: string;
  coverImage: string;
  status: "active" | "archived";
}

const emptyForm = (): Omit<Masterclass, "id"> => ({
  slug: "",
  title: "",
  description: "",
  description2: "",
  endingText: "",
  duration: "",
  price: "",
  coverImage: "",
  status: "active",
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<Masterclass>("/api/v1/masterclasses");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

// Auto-fills slug from title for a brand-new masterclass, right up until
// the user types into the slug field themselves — editing an existing
// masterclass's title never touches its (possibly already-live) slug.
const slugTouched = ref(false);
watch(
  () => form.value.title,
  (title) => {
    if (!slugTouched.value) form.value.slug = slugify(title);
  },
);

function startEdit(masterclass: Masterclass) {
  slugTouched.value = true;
  editingId.value = masterclass.id;
  form.value = {
    slug: masterclass.slug,
    title: masterclass.title,
    description: masterclass.description,
    description2: masterclass.description2,
    endingText: masterclass.endingText,
    duration: masterclass.duration,
    price: masterclass.price,
    coverImage: masterclass.coverImage,
    status: masterclass.status,
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
  if (!confirm("Удалить мастер-класс?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

async function onDuplicate(masterclass: Masterclass) {
  await create({
    slug: `${masterclass.slug}-copy-${Date.now()}`,
    title: `${masterclass.title} (копия)`,
    description: masterclass.description,
    description2: masterclass.description2,
    endingText: masterclass.endingText,
    duration: masterclass.duration,
    price: masterclass.price,
    coverImage: masterclass.coverImage,
    status: "archived",
  });
}

const searchQuery = ref("");
const filteredItems = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return items.value;
  return items.value.filter((masterclass) => masterclass.title.toLowerCase().includes(query));
});

const { selectedIds, isSelected, toggle, allSelected, toggleAll, clear } =
  useAdminBulkSelect(filteredItems);

async function onBulkDelete() {
  if (!confirm(`Удалить выбранные мастер-классы (${selectedIds.value.size})?`)) return;
  await Promise.all([...selectedIds.value].map((id) => remove(id)));
  clear();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Мастер-классы</h1>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.title"
        type="text"
        placeholder="Название *"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <input
        v-model="form.slug"
        type="text"
        placeholder="Slug *"
        required
        class="rounded border border-gray-300 px-3 py-2"
        @input="slugTouched = true"
      />
      <input
        v-model="form.duration"
        type="text"
        placeholder="Длительность * (например, «2-3 часа»)"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <AdminImageUpload v-model="form.coverImage" label="Обложка *" class="sm:col-span-2" />

      <AdminMarkdownField
        v-model="form.description"
        placeholder="Описание мастер-класса *"
        :rows="3"
        required
        class="sm:col-span-2"
      />
      <AdminMarkdownField
        v-model="form.description2"
        placeholder="Описание 2 (опционально)"
        :rows="2"
        class="sm:col-span-2"
      />
      <input
        v-model="form.price"
        type="text"
        placeholder="Цена * (например, «3000₽» или «3000₽ или 4500₽»)"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <AdminMarkdownField
        v-model="form.endingText"
        placeholder="Заключительный текст (опционально)"
        :rows="2"
        class="sm:col-span-2"
      />

      <select v-model="form.status" class="rounded border border-gray-300 px-3 py-2">
        <option value="active">Активен</option>
        <option value="archived">В архиве</option>
      </select>

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
          <th class="px-4 py-2">
            <input type="checkbox" class="size-5" :checked="allSelected" @change="toggleAll" />
          </th>
          <th class="px-4 py-2">Заголовок</th>
          <th class="px-4 py-2">Slug</th>
          <th class="px-4 py-2">Цена</th>
          <th class="px-4 py-2">Статус</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="masterclass in filteredItems"
          :key="masterclass.id"
          class="border-b border-gray-100 last:border-0"
        >
          <td class="px-4 py-2">
            <input
              type="checkbox"
              class="size-5"
              :checked="isSelected(masterclass.id)"
              @change="toggle(masterclass.id)"
            />
          </td>
          <td class="px-4 py-2">{{ masterclass.title }}</td>
          <td class="px-4 py-2">{{ masterclass.slug }}</td>
          <td class="px-4 py-2">{{ masterclass.price }}</td>
          <td class="px-4 py-2">{{ masterclass.status === "active" ? "Активен" : "В архиве" }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button
              class="text-[var(--color-primary)] hover:underline"
              @click="startEdit(masterclass)"
            >
              Редактировать
            </button>
            <button
              class="text-[var(--color-primary)] hover:underline"
              @click="onDuplicate(masterclass)"
            >
              Дублировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(masterclass.id)">
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
