<script setup lang="ts">
import { LEAD_STATUS_LABELS } from "~/lib/leadStatus";
import { filterClients, sortClients } from "~/lib/clientFilters";
import { readableTextColor } from "~/lib/tagColor";
import type { ClientListItem, ClientSortMode, Tag } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

const api = useApiClient();
const route = useRoute();
const router = useRouter();

const items = ref<ClientListItem[]>([]);
const loading = ref(false);
const error = ref("");

const productTagOptions = ref<Tag[]>([]);
const clientTypeTagOptions = ref<Tag[]>([]);

async function fetchAll() {
  loading.value = true;
  error.value = "";
  try {
    items.value = (await api<ClientListItem[]>("/api/v1/clients")) ?? [];
  } catch {
    error.value = "Не удалось загрузить клиентов";
  } finally {
    loading.value = false;
  }
}

async function fetchTagOptions() {
  [productTagOptions.value, clientTypeTagOptions.value] = await Promise.all([
    api<Tag[]>("/api/v1/tags", { query: { type: "product" } }),
    api<Tag[]>("/api/v1/tags", { query: { type: "client_type" } }),
  ]);
}

await Promise.all([fetchAll(), fetchTagOptions()]);

function parseCsvNumbers(value: unknown): Set<number> {
  if (typeof value !== "string" || value === "") return new Set();
  return new Set(value.split(",").map(Number));
}

const searchQuery = ref((route.query.q as string) ?? "");
const productTagFilter = ref<Set<number>>(parseCsvNumbers(route.query.ptags));
const clientTypeTagFilter = ref<Set<number>>(parseCsvNumbers(route.query.ctags));
const sortMode = ref<ClientSortMode>(route.query.sort === "name" ? "name" : "activity");

watch([searchQuery, productTagFilter, clientTypeTagFilter, sortMode], () => {
  router.replace({
    query: {
      ...route.query,
      q: searchQuery.value || undefined,
      ptags: productTagFilter.value.size > 0 ? [...productTagFilter.value].join(",") : undefined,
      ctags:
        clientTypeTagFilter.value.size > 0 ? [...clientTypeTagFilter.value].join(",") : undefined,
      sort: sortMode.value === "name" ? "name" : undefined,
    },
  });
});

function toggleProductTag(id: number) {
  const next = new Set(productTagFilter.value);
  if (next.has(id)) next.delete(id);
  else next.add(id);
  productTagFilter.value = next;
}

function toggleClientTypeTag(id: number) {
  const next = new Set(clientTypeTagFilter.value);
  if (next.has(id)) next.delete(id);
  else next.add(id);
  clientTypeTagFilter.value = next;
}

const filteredItems = computed(() =>
  sortClients(
    filterClients(items.value, {
      query: searchQuery.value,
      productTagIds: productTagFilter.value,
      clientTypeTagIds: clientTypeTagFilter.value,
    }),
    sortMode.value,
  ),
);
</script>

<template>
  <div>
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-2xl font-semibold">Клиенты</h1>
      <NuxtLink to="/admin/leads" class="text-sm text-[var(--color-primary)] hover:underline">
        К списку заявок →
      </NuxtLink>
    </div>

    <div class="mt-6 flex flex-wrap items-center gap-3">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Поиск по имени, телефону, почте…"
        class="min-w-64 rounded border border-gray-300 px-3 py-2 text-sm"
      />
      <select v-model="sortMode" class="rounded border border-gray-300 px-3 py-2 text-sm">
        <option value="activity">По активности</option>
        <option value="name">По имени</option>
      </select>
    </div>

    <div
      v-if="productTagOptions.length > 0 || clientTypeTagOptions.length > 0"
      class="mt-3 flex flex-wrap gap-4"
    >
      <div v-if="productTagOptions.length > 0" class="flex flex-wrap items-center gap-2">
        <span class="text-xs text-[var(--color-text-muted)]">Продукт:</span>
        <button
          v-for="tag in productTagOptions"
          :key="tag.id"
          type="button"
          class="rounded-full border px-3 py-1 text-xs"
          :class="
            productTagFilter.has(tag.id)
              ? 'border-[var(--color-primary)] bg-[var(--color-primary)] text-white'
              : 'border-gray-300 text-[var(--color-text-muted)]'
          "
          @click="toggleProductTag(tag.id)"
        >
          {{ tag.name }}
        </button>
      </div>
      <div v-if="clientTypeTagOptions.length > 0" class="flex flex-wrap items-center gap-2">
        <span class="text-xs text-[var(--color-text-muted)]">Тип клиента:</span>
        <button
          v-for="tag in clientTypeTagOptions"
          :key="tag.id"
          type="button"
          class="rounded-full border px-3 py-1 text-xs"
          :class="
            clientTypeTagFilter.has(tag.id)
              ? 'border-[var(--color-primary)] bg-[var(--color-primary)] text-white'
              : 'border-gray-300 text-[var(--color-text-muted)]'
          "
          @click="toggleClientTypeTag(tag.id)"
        >
          {{ tag.name }}
        </button>
      </div>
    </div>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 space-y-3">
      <NuxtLink
        v-for="client in filteredItems"
        :key="client.id"
        :to="`/admin/clients/${client.id}`"
        class="block rounded border border-gray-200 bg-white p-4 hover:border-[var(--color-primary)]"
      >
        <div class="flex flex-wrap items-center justify-between gap-2">
          <span class="font-medium text-[var(--color-primary)]">{{ client.name }}</span>
          <span
            v-if="client.latestStatus"
            class="rounded-full bg-gray-100 px-2 py-0.5 text-xs text-[var(--color-text-muted)]"
          >
            {{ LEAD_STATUS_LABELS[client.latestStatus] }}
          </span>
        </div>
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">
          {{ client.phone }}<span v-if="client.email"> · {{ client.email }}</span> ·
          {{ client.requestCount }} {{ client.requestCount === 1 ? "заявка" : "заявок" }}
        </p>
        <div
          v-if="client.productTags.length || client.clientTypeTags.length"
          class="mt-2 flex flex-wrap gap-1"
        >
          <span
            v-for="tag in [...client.productTags, ...client.clientTypeTags]"
            :key="tag.id"
            class="rounded-full px-2 py-0.5 text-xs"
            :style="{ backgroundColor: tag.color, color: readableTextColor(tag.color) }"
          >
            {{ tag.name }}
          </span>
        </div>
      </NuxtLink>

      <p
        v-if="filteredItems.length === 0"
        class="rounded border border-gray-200 bg-white px-4 py-6 text-center text-[var(--color-text-muted)]"
      >
        {{ items.length === 0 ? "Пока пусто" : "Ничего не найдено" }}
      </p>
    </div>
  </div>
</template>
