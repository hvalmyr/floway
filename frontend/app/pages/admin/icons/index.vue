<script setup lang="ts">
import type { Icon } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Icons have no update endpoint (see icon_handler.go's comment — replacing
// one is delete-and-reupload, not edit-in-place), so this page only needs
// list/create/delete from useAdminResource, unlike every other admin list
// here.
const { items, loading, error, fetchAll, create, remove } = useAdminResource<Icon>("/api/v1/icons");

await fetchAll();

const name = ref("");
const uploading = ref(false);
const uploadError = ref("");
const fileInput = ref<HTMLInputElement | null>(null);

function pickFile() {
  fileInput.value?.click();
}

async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = "";
  if (!file) return;

  uploadError.value = "";
  uploading.value = true;
  try {
    const svg = await file.text();
    await create({ name: name.value.trim() || file.name.replace(/\.svg$/i, ""), svg });
    name.value = "";
  } catch {
    uploadError.value = "Не удалось загрузить SVG — проверьте, что это корректный файл";
  } finally {
    uploading.value = false;
  }
}

async function onDelete(icon: Icon) {
  if (
    !confirm(`Удалить иконку «${icon.name}»? Места, где она уже выбрана, перестанут её показывать.`)
  )
    return;
  await remove(icon.id);
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Библиотека иконок</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      SVG-иконки, доступные для выбора везде, где есть поле «Иконка» (преимущества, бегущая строка
      сертификатов и т.д.), в дополнение к встроенному набору. Для правильной перекраски под разные
      фоны SVG должен использовать fill/stroke="currentColor".
    </p>

    <div class="mt-6 flex flex-wrap items-end gap-3 rounded border border-gray-200 bg-white p-4">
      <div class="flex flex-col gap-1">
        <label for="icon-name" class="text-sm font-medium">Название</label>
        <input
          id="icon-name"
          v-model="name"
          type="text"
          placeholder="Например, «Ромашка»"
          class="rounded border border-gray-300 px-3 py-2"
        />
      </div>
      <button
        type="button"
        :disabled="uploading"
        class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
        @click="pickFile"
      >
        {{ uploading ? "Загрузка…" : "Загрузить SVG" }}
      </button>
      <p v-if="uploadError" class="w-full text-sm text-red-600">{{ uploadError }}</p>
      <input
        ref="fileInput"
        type="file"
        accept=".svg,image/svg+xml"
        class="hidden"
        @change="onFileChange"
      />
    </div>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-4 lg:grid-cols-6">
      <div
        v-for="icon in items"
        :key="icon.id"
        class="flex flex-col items-center gap-2 rounded border border-gray-200 bg-white p-4 text-center"
      >
        <span class="grid size-10 place-items-center text-[var(--color-primary)]">
          <AppIcon :icon="`icon:${icon.id}`" class="size-8" />
        </span>
        <span class="truncate text-sm">{{ icon.name }}</span>
        <button type="button" class="text-xs text-red-600 hover:underline" @click="onDelete(icon)">
          Удалить
        </button>
      </div>
      <p v-if="items.length === 0" class="col-span-full text-center text-[var(--color-text-muted)]">
        Пока пусто
      </p>
    </div>
  </div>
</template>
