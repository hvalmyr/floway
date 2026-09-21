<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface SocialLinkEntity {
  id: number;
  label: string;
  href: string;
  disclaimer: string;
  sortOrder: number;
}

const emptyForm = () => ({
  label: "",
  href: "",
  disclaimer: "",
  sortOrder: 0,
});

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<SocialLinkEntity>("/api/v1/social-links");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

await fetchAll();

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(item: SocialLinkEntity) {
  editingId.value = item.id;
  form.value = {
    label: item.label,
    href: item.href,
    disclaimer: item.disclaimer,
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
  if (!confirm("Удалить ссылку на соцсеть?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Соцсети</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Иконки соцсетей в подвале сайта и на странице «Контакты». Иконка подбирается по названию —
      сейчас есть готовые иконки только для Telegram, VK и Instagram.
    </p>

    <form
      class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="form.label"
        type="text"
        placeholder="Название, например: Telegram"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <input
        v-model="form.href"
        type="url"
        placeholder="Ссылка, например: https://t.me/floway"
        required
        class="rounded border border-gray-300 px-3 py-2"
      />
      <textarea
        v-model="form.disclaimer"
        placeholder="Дисклеймер (необязательно) — например, обязательная пометка про Meta"
        rows="2"
        class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
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
          <th class="px-4 py-2" />
          <th class="px-4 py-2">Название</th>
          <th class="px-4 py-2">Ссылка</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(item, index) in items"
          :key="item.id"
          :data-row-index="index"
          class="border-b border-gray-100 last:border-0"
          :class="draggingIndex === index ? 'opacity-50' : ''"
        >
          <td class="px-4 py-2">
            <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
          </td>
          <td class="px-4 py-2">{{ item.label }}</td>
          <td class="px-4 py-2 text-[var(--color-text-muted)]">{{ item.href }}</td>
          <td class="flex gap-3 px-4 py-2 text-right">
            <button class="text-[var(--color-primary)] hover:underline" @click="startEdit(item)">
              Редактировать
            </button>
            <button class="text-red-600 hover:underline" @click="onDelete(item.id)">Удалить</button>
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
