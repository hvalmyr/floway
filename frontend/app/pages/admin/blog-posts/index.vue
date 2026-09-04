<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface BlogPost {
  id: number;
  slug: string;
  title: string;
  coverImage: string;
  category: string;
  tags: string[];
  author: string;
  publishedAt: string | null;
  content: string;
  status: "draft" | "published";
}

const emptyForm = () => ({
  slug: "",
  title: "",
  coverImage: "",
  category: "",
  tagsInput: "",
  author: "",
  publishedAtInput: "",
  content: "",
  status: "draft" as "draft" | "published",
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<BlogPost>("/api/v1/blog-posts");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

// Auto-fills slug from title for a brand-new post, right up until the user
// types into the slug field themselves — editing an existing post's title
// never touches its (possibly already-live) slug.
const slugTouched = ref(false);
watch(
  () => form.value.title,
  (title) => {
    if (!slugTouched.value) form.value.slug = slugify(title);
  },
);

function startEdit(post: BlogPost) {
  slugTouched.value = true;
  editingId.value = post.id;
  form.value = {
    slug: post.slug,
    title: post.title,
    coverImage: post.coverImage,
    category: post.category,
    tagsInput: post.tags.join(", "),
    author: post.author,
    publishedAtInput: post.publishedAt ?? "",
    content: post.content,
    status: post.status,
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
    const payload = {
      slug: form.value.slug,
      title: form.value.title,
      coverImage: form.value.coverImage,
      category: form.value.category,
      tags: form.value.tagsInput
        .split(",")
        .map((tag) => tag.trim())
        .filter((tag) => tag.length > 0),
      author: form.value.author,
      publishedAt:
        form.value.publishedAtInput.trim() === "" ? null : form.value.publishedAtInput.trim(),
      content: form.value.content,
      status: form.value.status,
    };
    if (editingId.value === null) {
      await create(payload);
    } else {
      await update(editingId.value, payload);
    }
    cancelEdit();
  } catch {
    formError.value = "Не удалось сохранить. Проверьте поля и попробуйте снова.";
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Удалить запись блога?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

async function onDuplicate(post: BlogPost) {
  await create({
    slug: `${post.slug}-copy-${Date.now()}`,
    title: `${post.title} (копия)`,
    coverImage: post.coverImage,
    category: post.category,
    tags: post.tags,
    author: post.author,
    publishedAt: null,
    content: post.content,
    status: "draft",
  });
}

const searchQuery = ref("");
const filteredItems = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return items.value;
  return items.value.filter((post) => post.title.toLowerCase().includes(query));
});

const { selectedIds, isSelected, toggle, allSelected, toggleAll, clear } =
  useAdminBulkSelect(filteredItems);

async function onBulkDelete() {
  if (!confirm(`Удалить выбранные записи блога (${selectedIds.value.size})?`)) return;
  await Promise.all([...selectedIds.value].map((id) => remove(id)));
  clear();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Блог</h1>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.title"
        type="text"
        placeholder="Заголовок"
        required
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <input
        v-model="form.slug"
        type="text"
        placeholder="Slug"
        required
        class="rounded border border-gray-300 px-3 py-2"
        @input="slugTouched = true"
      />
      <input
        v-model="form.author"
        type="text"
        placeholder="Автор"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <AdminImageUpload v-model="form.coverImage" label="Обложка" />
      <input
        v-model="form.category"
        type="text"
        placeholder="Категория"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model="form.tagsInput"
        type="text"
        placeholder="Теги через запятую, например: цветы, новости"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <input
        v-model="form.publishedAtInput"
        type="text"
        placeholder="Дата публикации, например: 2026-01-15T10:00:00Z"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <select v-model="form.status" class="rounded border border-gray-300 px-3 py-2">
        <option value="draft">Черновик</option>
        <option value="published">Опубликовано</option>
      </select>
      <AdminRichTextEditor v-model="form.content" placeholder="Содержимое" class="sm:col-span-2" />

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
        placeholder="Поиск по заголовку…"
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
          <th class="px-4 py-2">Статус</th>
          <th class="px-4 py-2">Дата публикации</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="post in filteredItems"
          :key="post.id"
          class="border-b border-gray-100 last:border-0"
        >
          <td class="px-4 py-2">
            <input
              type="checkbox"
              class="size-5"
              :checked="isSelected(post.id)"
              @change="toggle(post.id)"
            />
          </td>
          <td class="px-4 py-2">{{ post.title }}</td>
          <td class="px-4 py-2">{{ post.slug }}</td>
          <td class="px-4 py-2">{{ post.status === "published" ? "Опубликовано" : "Черновик" }}</td>
          <td class="px-4 py-2">{{ post.publishedAt ?? "—" }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(post)">
              Редактировать
            </button>
            <button class="text-[var(--color-primary)] hover:underline" @click="onDuplicate(post)">
              Дублировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(post.id)">Удалить</button>
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
