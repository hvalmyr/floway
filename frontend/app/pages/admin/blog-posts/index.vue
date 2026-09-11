<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

type DisplayStyle = "blue-beige" | "brown-beige" | "beige-blue" | "beige-brown";

const displayStyleLabels: Record<DisplayStyle, string> = {
  "blue-beige": "Голубой фон, бежевый текст",
  "brown-beige": "Коричневый фон, бежевый текст",
  "beige-blue": "Бежевый фон, голубой текст",
  "beige-brown": "Бежевый фон, коричневый текст",
};

interface BlogPost {
  id: number;
  slug: string;
  title: string;
  metaTitle: string;
  metaDescription: string;
  coverImage: string;
  displayStyle: DisplayStyle;
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
  metaTitle: "",
  metaDescription: "",
  coverImage: "",
  displayStyle: "blue-beige" as DisplayStyle,
  category: "",
  tags: [] as string[],
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
const slugTouched = ref(false);

// Restores an in-progress edit/create after a reload and keeps saving it as
// the admin types, so a refresh or an accidental tab close doesn't lose
// unsaved work. Cleared on successful submit and on explicit cancel.
const { clearDraft } = useAdminFormDraft("blog-posts", { editingId, form, slugTouched });

await fetchAll();

// Auto-fills slug from title for a brand-new post, right up until the user
// types into the slug field themselves — editing an existing post's title
// never touches its (possibly already-live) slug.
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
    metaTitle: post.metaTitle,
    metaDescription: post.metaDescription,
    coverImage: post.coverImage,
    displayStyle: post.displayStyle,
    category: post.category,
    tags: [...post.tags],
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
  clearDraft();
}

async function onSubmit() {
  formError.value = "";
  saving.value = true;
  try {
    const payload = {
      slug: form.value.slug,
      title: form.value.title,
      metaTitle: form.value.metaTitle,
      metaDescription: form.value.metaDescription,
      coverImage: form.value.coverImage,
      displayStyle: form.value.displayStyle,
      category: form.value.category,
      tags: form.value.tags,
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
    metaTitle: post.metaTitle,
    metaDescription: post.metaDescription,
    coverImage: post.coverImage,
    displayStyle: post.displayStyle,
    category: post.category,
    tags: post.tags,
    author: post.author,
    publishedAt: null,
    content: post.content,
    status: "draft",
  });
}

// Suggestions for the category/tag inputs below — drawn from posts already
// in the list rather than a separate backend endpoint, since blog
// category/tags are free-text columns (not the client Tag table's
// product/client-type system), so there's nothing structured to query yet.
const categorySuggestions = computed(() =>
  [...new Set(items.value.map((p) => p.category).filter((c) => c.length > 0))].sort(),
);
const tagSuggestions = computed(() => [...new Set(items.value.flatMap((p) => p.tags))].sort());

const tagInput = ref("");
function addTagFromInput() {
  const trimmed = tagInput.value.trim();
  tagInput.value = "";
  if (!trimmed) return;
  if (form.value.tags.some((t) => t.toLowerCase() === trimmed.toLowerCase())) return;
  form.value.tags.push(trimmed);
}
function removeTag(tag: string) {
  form.value.tags = form.value.tags.filter((t) => t !== tag);
}
function onTagInputKeydown(event: KeyboardEvent) {
  if (event.key !== "Enter" && event.key !== ",") return;
  event.preventDefault();
  addTagFromInput();
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
      <select v-model="form.displayStyle" class="rounded border border-gray-300 px-3 py-2">
        <option v-for="(label, value) in displayStyleLabels" :key="value" :value="value">
          {{ label }}
        </option>
      </select>

      <div class="flex flex-col gap-1 sm:col-span-2">
        <input
          v-model="form.metaTitle"
          type="text"
          placeholder="Meta title (заголовок в поиске и вкладке браузера)"
          maxlength="60"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <p class="text-right text-xs text-[var(--color-text-muted)]">
          {{ form.metaTitle.length }}/60 символов
          <span v-if="form.metaTitle.length === 0">— иначе используется заголовок статьи</span>
        </p>
      </div>
      <div class="flex flex-col gap-1 sm:col-span-2">
        <textarea
          v-model="form.metaDescription"
          placeholder="Meta description (описание в выдаче поисковика)"
          maxlength="160"
          rows="2"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <p class="text-right text-xs text-[var(--color-text-muted)]">
          {{ form.metaDescription.length }}/160 символов
        </p>
      </div>

      <div class="flex flex-col gap-1">
        <input
          v-model="form.category"
          type="text"
          list="blog-category-suggestions"
          placeholder="Категория"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <datalist id="blog-category-suggestions">
          <option v-for="c in categorySuggestions" :key="c" :value="c" />
        </datalist>
      </div>

      <div class="flex flex-col gap-1">
        <div class="flex flex-wrap items-center gap-2 rounded border border-gray-300 p-2">
          <span
            v-for="tag in form.tags"
            :key="tag"
            class="flex items-center gap-1 rounded-full bg-gray-100 px-2 py-0.5 text-xs"
          >
            {{ tag }}
            <button type="button" class="hover:text-red-600" @click="removeTag(tag)">×</button>
          </span>
          <input
            v-model="tagInput"
            type="text"
            list="blog-tag-suggestions"
            placeholder="Добавить тег…"
            class="min-w-32 flex-1 border-none p-1 text-sm outline-none"
            @keydown="onTagInputKeydown"
            @blur="addTagFromInput"
          />
        </div>
        <datalist id="blog-tag-suggestions">
          <option v-for="t in tagSuggestions" :key="t" :value="t" />
        </datalist>
      </div>

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
