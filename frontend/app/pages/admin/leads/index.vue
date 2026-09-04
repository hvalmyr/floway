<script setup lang="ts">
import { LEAD_STATUSES, LEAD_STATUS_LABELS } from "~/lib/leadStatus";
import { contactMethodLabels } from "~/lib/leadLabels";
import { filterLeads, formatLeadExcerpt, lastContactAt, sortLeads } from "~/lib/leadFilters";
import type { LeadListItem, LeadSortMode, LeadStatus, Tag } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

const api = useApiClient();
const route = useRoute();
const router = useRouter();

const items = ref<LeadListItem[]>([]);
const loading = ref(false);
const error = ref("");

const productTagOptions = ref<Tag[]>([]);
const clientTypeTagOptions = ref<Tag[]>([]);

const conversion = ref<{
  closedWon: number;
  closedLost: number;
  conversionRate: number | null;
} | null>(null);

async function fetchAll() {
  loading.value = true;
  error.value = "";
  try {
    items.value = (await api<LeadListItem[]>("/api/v1/leads")) ?? [];
  } catch {
    error.value = "Не удалось загрузить заявки";
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

async function fetchConversion() {
  conversion.value = await api("/api/v1/leads/stats/conversion");
}

await Promise.all([fetchAll(), fetchTagOptions(), fetchConversion()]);

// Filter/sort state lives in the URL query string so it survives a trip to
// a client's detail page and back, instead of resetting on every navigation.
function parseCsvNumbers(value: unknown): Set<number> {
  if (typeof value !== "string" || value === "") return new Set();
  return new Set(value.split(",").map(Number));
}
function parseCsvStatuses(value: unknown): Set<LeadStatus> {
  if (typeof value !== "string" || value === "") return new Set();
  return new Set(value.split(",") as LeadStatus[]);
}

const searchQuery = ref((route.query.q as string) ?? "");
const statusFilter = ref<Set<LeadStatus>>(parseCsvStatuses(route.query.status));
const productTagFilter = ref<Set<number>>(parseCsvNumbers(route.query.ptags));
const clientTypeTagFilter = ref<Set<number>>(parseCsvNumbers(route.query.ctags));
const staleOnly = ref(route.query.stale === "1");
const sortMode = ref<LeadSortMode>(route.query.sort === "nextAction" ? "nextAction" : "createdAt");

watch(
  [searchQuery, statusFilter, productTagFilter, clientTypeTagFilter, staleOnly, sortMode],
  () => {
    router.replace({
      query: {
        ...route.query,
        q: searchQuery.value || undefined,
        status: statusFilter.value.size > 0 ? [...statusFilter.value].join(",") : undefined,
        ptags: productTagFilter.value.size > 0 ? [...productTagFilter.value].join(",") : undefined,
        ctags:
          clientTypeTagFilter.value.size > 0 ? [...clientTypeTagFilter.value].join(",") : undefined,
        stale: staleOnly.value ? "1" : undefined,
        sort: sortMode.value === "nextAction" ? "nextAction" : undefined,
      },
    });
  },
);

function toggleStatus(status: LeadStatus) {
  const next = new Set(statusFilter.value);
  if (next.has(status)) next.delete(status);
  else next.add(status);
  statusFilter.value = next;
}

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
  sortLeads(
    filterLeads(items.value, {
      query: searchQuery.value,
      statuses: statusFilter.value,
      productTagIds: productTagFilter.value,
      clientTypeTagIds: clientTypeTagFilter.value,
      staleOnly: staleOnly.value,
    }),
    sortMode.value,
  ),
);

function isToday(iso: string): boolean {
  return iso.slice(0, 10) === new Date().toISOString().slice(0, 10);
}
const dueTodayItems = computed(() =>
  items.value.filter((i) => i.nextReminderAt && isToday(i.nextReminderAt)),
);

async function onStatusChange(lead: LeadListItem, status: string) {
  const updated = await api<LeadListItem>(`/api/v1/leads/${lead.id}/status`, {
    method: "PATCH",
    body: { status },
  });
  const idx = items.value.findIndex((i) => i.id === lead.id);
  if (idx !== -1) items.value[idx] = { ...items.value[idx], ...updated };
}

async function onDismissReview(lead: LeadListItem) {
  const updated = await api<LeadListItem>(`/api/v1/leads/${lead.id}/review`, { method: "PATCH" });
  const idx = items.value.findIndex((i) => i.id === lead.id);
  if (idx !== -1) items.value[idx] = { ...items.value[idx], ...updated };
}

async function onDelete(id: number) {
  if (!confirm("Удалить заявку?")) return;
  await api(`/api/v1/leads/${id}`, { method: "DELETE" });
  items.value = items.value.filter((i) => i.id !== id);
}

const { selectedIds, isSelected, toggle, allSelected, toggleAll, clear } =
  useAdminBulkSelect(filteredItems);

async function onBulkDelete() {
  if (!confirm(`Удалить выбранные заявки (${selectedIds.value.size})?`)) return;
  await Promise.all(
    [...selectedIds.value].map((id) => api(`/api/v1/leads/${id}`, { method: "DELETE" })),
  );
  items.value = items.value.filter((i) => !selectedIds.value.has(i.id));
  clear();
}

async function onBulkStatusChange(status: string) {
  if (!status) return;
  const ids = [...selectedIds.value];
  const updated = await Promise.all(
    ids.map((id) =>
      api<LeadListItem>(`/api/v1/leads/${id}/status`, { method: "PATCH", body: { status } }),
    ),
  );
  for (const lead of updated) {
    const idx = items.value.findIndex((i) => i.id === lead.id);
    if (idx !== -1) items.value[idx] = { ...items.value[idx], ...lead };
  }
  clear();
}
</script>

<template>
  <div>
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-2xl font-semibold">Заявки</h1>
      <p v-if="conversion" class="text-sm text-[var(--color-text-muted)]">
        Конверсия:
        {{
          conversion.conversionRate === null
            ? "—"
            : `${Math.round(conversion.conversionRate * 100)}%`
        }}
        ({{ conversion.closedWon }} успешно / {{ conversion.closedLost }} отказ)
      </p>
    </div>

    <div
      v-if="dueTodayItems.length > 0"
      class="mt-4 rounded border border-amber-300 bg-amber-50 px-4 py-2 text-sm"
    >
      Сегодня напомнить связаться:
      <NuxtLink
        v-for="(item, idx) in dueTodayItems"
        :key="item.clientId"
        :to="`/admin/clients/${item.clientId}`"
        class="text-[var(--color-primary)] hover:underline"
      >
        {{ item.client.name }}<span v-if="idx < dueTodayItems.length - 1">, </span>
      </NuxtLink>
    </div>

    <div class="mt-6 flex flex-wrap items-center gap-3">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Поиск по имени, телефону, почте…"
        class="min-w-64 rounded border border-gray-300 px-3 py-2 text-sm"
      />
      <button
        type="button"
        class="rounded border px-3 py-2 text-xs"
        :class="
          staleOnly
            ? 'border-[var(--color-primary)] bg-[var(--color-primary)] text-white'
            : 'border-gray-300'
        "
        @click="staleOnly = !staleOnly"
      >
        Не отвечали 3+ дня
      </button>
      <select v-model="sortMode" class="rounded border border-gray-300 px-3 py-2 text-sm">
        <option value="createdAt">Сначала новые</option>
        <option value="nextAction">По дате напоминания</option>
      </select>
    </div>

    <div class="mt-3 flex flex-wrap gap-2">
      <button
        v-for="status in LEAD_STATUSES"
        :key="status"
        type="button"
        class="rounded-full border px-3 py-1 text-xs"
        :class="
          statusFilter.has(status)
            ? 'border-[var(--color-primary)] bg-[var(--color-primary)] text-white'
            : 'border-gray-300 text-[var(--color-text-muted)]'
        "
        @click="toggleStatus(status)"
      >
        {{ LEAD_STATUS_LABELS[status] }}
      </button>
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

    <div
      v-if="selectedIds.size > 0"
      class="mt-3 flex flex-wrap items-center gap-3 rounded border border-[var(--color-primary)] bg-white p-3 text-sm"
    >
      <span>Выбрано: {{ selectedIds.size }}</span>
      <select
        class="rounded border border-gray-300 px-2 py-1"
        value=""
        @change="onBulkStatusChange(($event.target as HTMLSelectElement).value)"
      >
        <option value="" disabled>Сменить статус…</option>
        <option v-for="status in LEAD_STATUSES" :key="status" :value="status">
          {{ LEAD_STATUS_LABELS[status] }}
        </option>
      </select>
      <button class="text-red-600 hover:underline" @click="onBulkDelete">Удалить выбранные</button>
      <button class="text-[var(--color-text-muted)] hover:underline" @click="clear">
        Отменить выбор
      </button>
    </div>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 space-y-3">
      <div
        v-for="lead in filteredItems"
        :key="lead.id"
        class="flex flex-wrap items-start justify-between gap-4 rounded border border-gray-200 bg-white p-4"
        :class="lead.status === 'postponed' ? 'opacity-60' : ''"
      >
        <div class="min-w-0 flex-1">
          <div class="flex flex-wrap items-center gap-2">
            <input
              type="checkbox"
              class="size-4"
              :checked="isSelected(lead.id)"
              @change="toggle(lead.id)"
            />
            <NuxtLink
              :to="`/admin/clients/${lead.clientId}`"
              class="font-medium text-[var(--color-primary)] hover:underline"
            >
              {{ lead.client.name }}
            </NuxtLink>
            <span class="text-xs text-[var(--color-text-muted)]">{{ lead.client.phone }}</span>
            <button
              v-if="lead.needsStatusReview"
              type="button"
              title="Автоматически перенесено из старого статуса «закрыта» — проверьте и подтвердите"
              class="rounded-full border border-amber-400 bg-amber-100 px-2 py-0.5 text-xs text-amber-800"
              @click="onDismissReview(lead)"
            >
              Проверить статус
            </button>
          </div>
          <p class="mt-1 text-sm text-[var(--color-text-muted)]">
            {{ new Date(lastContactAt(lead)).toLocaleDateString("ru-RU") }} ·
            {{ formatLeadExcerpt(lead) }} ·
            {{ contactMethodLabels[lead.contactMethod] ?? lead.contactMethod }}
          </p>
          <div
            v-if="lead.productTags.length || lead.clientTypeTags.length"
            class="mt-2 flex flex-wrap gap-1"
          >
            <span
              v-for="tag in [...lead.productTags, ...lead.clientTypeTags]"
              :key="tag.id"
              class="rounded-full bg-gray-100 px-2 py-0.5 text-xs text-[var(--color-text-muted)]"
            >
              {{ tag.name }}
            </span>
          </div>
        </div>

        <div class="flex shrink-0 items-center gap-3">
          <select
            class="rounded border border-gray-300 px-2 py-1 text-sm"
            :value="lead.status"
            @change="onStatusChange(lead, ($event.target as HTMLSelectElement).value)"
          >
            <option v-for="status in LEAD_STATUSES" :key="status" :value="status">
              {{ LEAD_STATUS_LABELS[status] }}
            </option>
          </select>
          <button class="text-sm text-red-600 hover:underline" @click="onDelete(lead.id)">
            Удалить
          </button>
        </div>
      </div>

      <p
        v-if="filteredItems.length === 0"
        class="rounded border border-gray-200 bg-white px-4 py-6 text-center text-[var(--color-text-muted)]"
      >
        {{ items.length === 0 ? "Пока пусто" : "Ничего не найдено" }}
      </p>
    </div>
  </div>
</template>
