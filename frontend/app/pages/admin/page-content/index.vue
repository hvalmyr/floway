<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface PageContentItem {
  key: string;
  label: string;
  value: string;
  type: "text" | "image" | "icon";
  updatedAt: string;
}

const api = useApiClient();
const items = ref<PageContentItem[]>([]);
const loading = ref(false);
const error = ref("");
const savingKey = ref<string | null>(null);
const savedKey = ref<string | null>(null);

async function fetchAll() {
  loading.value = true;
  error.value = "";
  try {
    items.value = (await api<PageContentItem[]>("/api/v1/page-content")) ?? [];
  } catch {
    error.value = "Не удалось загрузить данные";
  } finally {
    loading.value = false;
  }
}

await fetchAll();

async function save(item: PageContentItem) {
  savingKey.value = item.key;
  savedKey.value = null;
  try {
    const updated = await api<PageContentItem>(`/api/v1/page-content/${item.key}`, {
      method: "PUT",
      body: { value: item.value },
    });
    item.value = updated.value;
    savedKey.value = item.key;
  } catch {
    error.value = `Не удалось сохранить «${item.label}»`;
  } finally {
    savingKey.value = null;
  }
}

function onImageUploaded(item: PageContentItem, url: string) {
  item.value = url;
  save(item);
}

function onIconSelected(item: PageContentItem, value: string) {
  item.value = value;
  save(item);
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Тексты сайта</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Свободные тексты и блоки сайта (заголовки, вводные абзацы, юридические страницы). Поддерживают
      markdown там, где это отмечено в подписи поля.
    </p>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 flex flex-col gap-4">
      <div
        v-for="item in items"
        :key="item.key"
        class="rounded border border-gray-200 bg-white p-4"
      >
        <div class="mb-2 flex items-center justify-between gap-4">
          <label :for="`field-${item.key}`" class="text-sm font-medium">{{ item.label }}</label>
          <span class="font-mono text-xs text-[var(--color-text-muted)]">{{ item.key }}</span>
        </div>
        <AdminImageUpload
          v-if="item.type === 'image'"
          :model-value="item.value"
          :label="item.label"
          @update:model-value="(url) => onImageUploaded(item, url)"
        />
        <AdminIconPicker
          v-else-if="item.type === 'icon'"
          :model-value="item.value"
          @update:model-value="(value) => onIconSelected(item, value)"
        />
        <AdminMarkdownField v-else :id="`field-${item.key}`" v-model="item.value" :rows="3" />
        <div class="mt-2 flex items-center gap-3">
          <button
            v-if="item.type !== 'image' && item.type !== 'icon'"
            type="button"
            :disabled="savingKey === item.key"
            class="rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white disabled:opacity-50"
            @click="save(item)"
          >
            {{ savingKey === item.key ? "Сохранение…" : "Сохранить" }}
          </button>
          <span v-if="savedKey === item.key" class="text-sm text-green-600">Сохранено</span>
        </div>
      </div>
      <p v-if="items.length === 0" class="text-center text-[var(--color-text-muted)]">Пока пусто</p>
    </div>
  </div>
</template>
