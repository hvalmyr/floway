<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface Masterclass {
  id: number;
  slug: string;
  title: string;
  shortDescription: string;
  fullDescription: string;
  endingText: string;
  duration: string;
  priceGroup: number;
  priceIndividual: number;
  priceDescription: string;
  coverImage: string;
  status: "active" | "archived";
}

const emptyForm = (): Omit<Masterclass, "id"> => ({
  slug: "",
  title: "",
  shortDescription: "",
  fullDescription: "",
  endingText: "",
  duration: "",
  priceGroup: 0,
  priceIndividual: 0,
  priceDescription: "",
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

function startEdit(masterclass: Masterclass) {
  editingId.value = masterclass.id;
  form.value = {
    slug: masterclass.slug,
    title: masterclass.title,
    shortDescription: masterclass.shortDescription,
    fullDescription: masterclass.fullDescription,
    endingText: masterclass.endingText,
    duration: masterclass.duration,
    priceGroup: masterclass.priceGroup,
    priceIndividual: masterclass.priceIndividual,
    priceDescription: masterclass.priceDescription,
    coverImage: masterclass.coverImage,
    status: masterclass.status,
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
  if (!confirm("Удалить мастер-класс?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
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
      />
      <input
        v-model="form.duration"
        type="text"
        placeholder="Длительность"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <AdminImageUpload v-model="form.coverImage" label="Обложка" class="sm:col-span-2" />

      <textarea
        v-model="form.shortDescription"
        placeholder="Краткое описание"
        rows="2"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <textarea
        v-model="form.fullDescription"
        placeholder="Полное описание"
        rows="3"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />
      <textarea
        v-model="form.endingText"
        placeholder="Завершающий текст"
        rows="2"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
      />

      <input
        v-model.number="form.priceGroup"
        type="number"
        placeholder="Цена (группа)"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model.number="form.priceIndividual"
        type="number"
        placeholder="Цена (индивидуально)"
        class="rounded border border-gray-300 px-3 py-2"
      />
      <textarea
        v-model="form.priceDescription"
        placeholder="Описание цены"
        rows="2"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
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

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <table
      v-else
      class="mt-6 w-full border-collapse overflow-hidden rounded border border-gray-200 bg-white text-sm"
    >
      <thead>
        <tr class="border-b border-gray-200 text-left">
          <th class="px-4 py-2">Заголовок</th>
          <th class="px-4 py-2">Slug</th>
          <th class="px-4 py-2">Цена (группа/индивидуально)</th>
          <th class="px-4 py-2">Статус</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="masterclass in items"
          :key="masterclass.id"
          class="border-b border-gray-100 last:border-0"
        >
          <td class="px-4 py-2">{{ masterclass.title }}</td>
          <td class="px-4 py-2">{{ masterclass.slug }}</td>
          <td class="px-4 py-2">
            {{ masterclass.priceGroup }} / {{ masterclass.priceIndividual }}
          </td>
          <td class="px-4 py-2">{{ masterclass.status === "active" ? "Активен" : "В архиве" }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button
              class="text-[var(--color-primary)] hover:underline"
              @click="startEdit(masterclass)"
            >
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(masterclass.id)">
              Удалить
            </button>
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
