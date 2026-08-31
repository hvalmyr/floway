<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface Lead {
  id: number;
  name: string;
  phone: string;
  email: string;
  contactMethod: string;
  source: string;
  requestType: string;
  relatedId: number | null;
  relatedSlug: string;
  status: "new" | "in_progress" | "closed";
  createdAt: string;
}

const contactMethodLabels: Record<string, string> = {
  call: "Звонок",
  telegram: "Telegram",
  whatsapp: "WhatsApp",
  max: "MAX",
};

const sourceLabels: Record<string, string> = {
  referral: "Рекомендация",
  ads: "Реклама",
  internet: "Интернет",
  social: "Соцсети",
  maps: "Карты",
};

const requestTypeLabels: Record<string, string> = {
  course: "Курс",
  masterclass: "Мастер-класс",
  trial_lesson: "Пробный урок",
};

const statusLabels: Record<string, string> = {
  new: "Новая",
  in_progress: "В работе",
  closed: "Закрыта",
};

const api = useApiClient();
const items = ref<Lead[]>([]);
const loading = ref(false);
const error = ref("");

async function fetchAll() {
  loading.value = true;
  error.value = "";
  try {
    items.value = (await api<Lead[]>("/api/v1/leads")) ?? [];
  } catch {
    error.value = "Не удалось загрузить заявки";
  } finally {
    loading.value = false;
  }
}

async function onStatusChange(lead: Lead, status: string) {
  const updated = await api<Lead>(`/api/v1/leads/${lead.id}/status`, {
    method: "PATCH",
    body: { status },
  });
  const idx = items.value.findIndex((i) => i.id === lead.id);
  if (idx !== -1) items.value[idx] = updated;
}

async function onDelete(id: number) {
  if (!confirm("Удалить заявку?")) return;
  await api(`/api/v1/leads/${id}`, { method: "DELETE" });
  items.value = items.value.filter((i) => i.id !== id);
}

await fetchAll();

const searchQuery = ref("");
const statusFilter = ref<"" | Lead["status"]>("");

const filteredItems = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  return items.value.filter((lead) => {
    if (statusFilter.value && lead.status !== statusFilter.value) return false;
    if (!query) return true;
    return (
      lead.name.toLowerCase().includes(query) ||
      lead.phone.toLowerCase().includes(query) ||
      lead.email.toLowerCase().includes(query)
    );
  });
});

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
    ids.map((id) => api<Lead>(`/api/v1/leads/${id}/status`, { method: "PATCH", body: { status } })),
  );
  for (const lead of updated) {
    const idx = items.value.findIndex((i) => i.id === lead.id);
    if (idx !== -1) items.value[idx] = lead;
  }
  clear();
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Заявки</h1>

    <div class="mt-6 flex flex-wrap items-center gap-3">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Поиск по имени, телефону, почте…"
        class="min-w-64 rounded border border-gray-300 px-3 py-2 text-sm"
      />
      <select v-model="statusFilter" class="rounded border border-gray-300 px-3 py-2 text-sm">
        <option value="">Все статусы</option>
        <option value="new">{{ statusLabels.new }}</option>
        <option value="in_progress">{{ statusLabels.in_progress }}</option>
        <option value="closed">{{ statusLabels.closed }}</option>
      </select>
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
        <option value="new">{{ statusLabels.new }}</option>
        <option value="in_progress">{{ statusLabels.in_progress }}</option>
        <option value="closed">{{ statusLabels.closed }}</option>
      </select>
      <button class="text-red-600 hover:underline" @click="onBulkDelete">Удалить выбранные</button>
      <button class="text-[var(--color-text-muted)] hover:underline" @click="clear">
        Отменить выбор
      </button>
    </div>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 overflow-x-auto rounded border border-gray-200 bg-white">
      <table class="w-full border-collapse text-sm">
        <thead>
          <tr class="border-b border-gray-200 text-left">
            <th class="px-4 py-2">
              <input type="checkbox" class="size-5" :checked="allSelected" @change="toggleAll" />
            </th>
            <th class="px-4 py-2">Имя</th>
            <th class="px-4 py-2">Телефон</th>
            <th class="px-4 py-2">Email</th>
            <th class="px-4 py-2">Способ связи</th>
            <th class="px-4 py-2">Источник</th>
            <th class="px-4 py-2">Тип заявки</th>
            <th class="px-4 py-2">Страница</th>
            <th class="px-4 py-2">Статус</th>
            <th class="px-4 py-2">Дата создания</th>
            <th class="px-4 py-2" />
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="lead in filteredItems"
            :key="lead.id"
            class="border-b border-gray-100 last:border-0"
          >
            <td class="px-4 py-2">
              <input
                type="checkbox"
                class="size-5"
                :checked="isSelected(lead.id)"
                @change="toggle(lead.id)"
              />
            </td>
            <td class="px-4 py-2">{{ lead.name }}</td>
            <td class="px-4 py-2">{{ lead.phone }}</td>
            <td class="px-4 py-2">{{ lead.email }}</td>
            <td class="px-4 py-2">
              {{ contactMethodLabels[lead.contactMethod] ?? lead.contactMethod }}
            </td>
            <td class="px-4 py-2">{{ sourceLabels[lead.source] ?? lead.source }}</td>
            <td class="px-4 py-2">{{ requestTypeLabels[lead.requestType] ?? lead.requestType }}</td>
            <td class="px-4 py-2">
              <a
                v-if="lead.requestType === 'course' && lead.relatedSlug"
                :href="`/courses/${lead.relatedSlug}`"
                target="_blank"
                class="text-[var(--color-primary)] hover:underline"
                >{{ lead.relatedSlug }}</a
              >
              <span v-else>{{ lead.relatedSlug || "—" }}</span>
            </td>
            <td class="px-4 py-2">
              <select
                class="rounded border border-gray-300 px-2 py-1"
                :value="lead.status"
                @change="onStatusChange(lead, ($event.target as HTMLSelectElement).value)"
              >
                <option value="new">{{ statusLabels.new }}</option>
                <option value="in_progress">{{ statusLabels.in_progress }}</option>
                <option value="closed">{{ statusLabels.closed }}</option>
              </select>
            </td>
            <td class="px-4 py-2">{{ new Date(lead.createdAt).toLocaleString("ru-RU") }}</td>
            <td class="px-4 py-2 text-right">
              <button class="text-red-600 hover:underline" @click="onDelete(lead.id)">
                Удалить
              </button>
            </td>
          </tr>
          <tr v-if="filteredItems.length === 0">
            <td colspan="11" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
              {{ items.length === 0 ? "Пока пусто" : "Ничего не найдено" }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
