<script setup lang="ts">
import type { SiteButton, SiteButtonVariant } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

const api = useApiClient();
const items = ref<SiteButton[]>([]);
const loading = ref(false);
const error = ref("");
const savingKey = ref<string | null>(null);
const savedKey = ref<string | null>(null);

async function fetchAll() {
  loading.value = true;
  error.value = "";
  try {
    items.value = (await api<SiteButton[]>("/api/v1/site-buttons")) ?? [];
  } catch {
    error.value = "Не удалось загрузить данные";
  } finally {
    loading.value = false;
  }
}

await fetchAll();

async function save(item: SiteButton) {
  savingKey.value = item.key;
  savedKey.value = null;
  try {
    const updated = await api<SiteButton>(`/api/v1/site-buttons/${item.key}`, {
      method: "PUT",
      body: { text: item.text, variant: item.variant, url: item.url },
    });
    item.text = updated.text;
    item.variant = updated.variant;
    item.url = updated.url;
    savedKey.value = item.key;
  } catch {
    error.value = `Не удалось сохранить «${item.label}»`;
  } finally {
    savingKey.value = null;
  }
}

const VARIANT_LABELS: Record<SiteButtonVariant, string> = {
  primary: "Основной (синяя заливка)",
  outline: "Контурный (белый, с рамкой)",
};
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Кнопки сайта</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Текст, стиль и ссылка для отдельных кнопок-переходов на сайте (не кнопок внутри карточек
      курсов/мастер-классов и не кнопок отправки форм — у тех своя логика).
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
          <p class="text-sm font-medium">{{ item.label }}</p>
          <span class="font-mono text-xs text-[var(--color-text-muted)]">{{ item.key }}</span>
        </div>
        <div class="grid gap-3 sm:grid-cols-3">
          <label class="flex flex-col gap-1 text-sm">
            Текст
            <input
              v-model="item.text"
              type="text"
              class="rounded border border-gray-300 px-3 py-2"
            />
          </label>
          <label class="flex flex-col gap-1 text-sm">
            Стиль
            <select v-model="item.variant" class="rounded border border-gray-300 px-3 py-2">
              <option v-for="(label, variant) in VARIANT_LABELS" :key="variant" :value="variant">
                {{ label }}
              </option>
            </select>
          </label>
          <label class="flex flex-col gap-1 text-sm">
            Ссылка
            <input
              v-model="item.url"
              type="text"
              class="rounded border border-gray-300 px-3 py-2"
            />
          </label>
        </div>
        <div class="mt-3 flex items-center gap-3">
          <button
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
