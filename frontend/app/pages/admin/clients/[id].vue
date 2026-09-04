<script setup lang="ts">
import { LEAD_STATUS_LABELS } from "~/lib/leadStatus";
import { contactMethodLabels, requestTypeLabels, sourceLabels } from "~/lib/leadLabels";
import type { ClientDetail } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

const route = useRoute();
const clientId = Number(route.params.id);
const api = useApiClient();

const detail = ref<ClientDetail | null>(null);
const loading = ref(false);
const error = ref("");

async function fetchDetail() {
  loading.value = true;
  error.value = "";
  try {
    detail.value = await api<ClientDetail>(`/api/v1/clients/${clientId}`);
  } catch {
    error.value = "Не удалось загрузить карточку клиента";
  } finally {
    loading.value = false;
  }
}

await fetchDetail();

const newComment = ref("");
const commentSaving = ref(false);

async function addComment() {
  const text = newComment.value.trim();
  if (!text || !detail.value) return;
  commentSaving.value = true;
  try {
    const comment = await api(`/api/v1/clients/${clientId}/comments`, {
      method: "POST",
      body: { text },
    });
    detail.value.comments = [comment, ...detail.value.comments];
    newComment.value = "";
  } finally {
    commentSaving.value = false;
  }
}

const customReminderDays = ref<number | null>(null);

async function addReminder(days: number) {
  if (!detail.value || days < 1) return;
  const reminder = await api(`/api/v1/clients/${clientId}/reminders`, {
    method: "POST",
    body: { days, note: "" },
  });
  detail.value.reminders = [...detail.value.reminders, reminder];
}

async function completeReminder(reminderId: number) {
  if (!detail.value) return;
  await api(`/api/v1/clients/${clientId}/reminders/${reminderId}/complete`, { method: "PATCH" });
  detail.value.reminders = detail.value.reminders.filter((r) => r.id !== reminderId);
}

function isToday(iso: string): boolean {
  return iso.slice(0, 10) === new Date().toISOString().slice(0, 10);
}

const openReminders = computed(() => detail.value?.reminders.filter((r) => !r.completedAt) ?? []);
</script>

<template>
  <div>
    <NuxtLink to="/admin/leads" class="text-sm text-[var(--color-primary)] hover:underline">
      ← К списку заявок
    </NuxtLink>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>

    <template v-else-if="detail">
      <div class="mt-4 flex flex-wrap items-baseline justify-between gap-3">
        <h1 class="text-2xl font-semibold">{{ detail.name }}</h1>
        <div class="text-sm text-[var(--color-text-muted)]">
          {{ detail.phone }} · {{ detail.email || "—" }}
        </div>
      </div>

      <section class="mt-6 grid gap-4 sm:grid-cols-2">
        <div>
          <h2 class="text-lg font-medium">Теги по продукту</h2>
          <AdminTagCombobox
            class="mt-2"
            :client-id="detail.id"
            tag-type="product"
            v-model="detail.productTags"
          />
        </div>
        <div>
          <h2 class="text-lg font-medium">Теги по типу клиента</h2>
          <AdminTagCombobox
            class="mt-2"
            :client-id="detail.id"
            tag-type="client_type"
            v-model="detail.clientTypeTags"
          />
        </div>
      </section>

      <section class="mt-6">
        <h2 class="text-lg font-medium">Напоминания</h2>
        <div class="mt-2 flex flex-wrap items-center gap-2">
          <button
            v-for="days in [1, 3, 7]"
            :key="days"
            type="button"
            class="rounded border border-gray-300 px-3 py-1.5 text-sm hover:bg-gray-50"
            @click="addReminder(days)"
          >
            Через {{ days }} {{ days === 1 ? "день" : "дня" }}
          </button>
          <input
            v-model.number="customReminderDays"
            type="number"
            min="1"
            placeholder="N дней"
            class="w-24 rounded border border-gray-300 px-2 py-1.5 text-sm"
          />
          <button
            type="button"
            class="rounded border border-gray-300 px-3 py-1.5 text-sm hover:bg-gray-50"
            @click="customReminderDays && addReminder(customReminderDays)"
          >
            Добавить
          </button>
        </div>
        <ul v-if="openReminders.length > 0" class="mt-3 space-y-1">
          <li
            v-for="reminder in openReminders"
            :key="reminder.id"
            class="flex items-center gap-2 text-sm"
          >
            <span
              v-if="isToday(reminder.remindAt)"
              class="rounded-full bg-amber-100 px-2 py-0.5 text-xs text-amber-800"
            >
              Сегодня
            </span>
            <span>{{ new Date(reminder.remindAt).toLocaleDateString("ru-RU") }}</span>
            <span v-if="reminder.note" class="text-[var(--color-text-muted)]"
              >— {{ reminder.note }}</span
            >
            <button
              type="button"
              class="text-[var(--color-text-muted)] hover:underline"
              @click="completeReminder(reminder.id)"
            >
              Выполнено
            </button>
          </li>
        </ul>
        <p v-else class="mt-3 text-sm text-[var(--color-text-muted)]">Нет активных напоминаний</p>
      </section>

      <section class="mt-6">
        <h2 class="text-lg font-medium">История заявок</h2>
        <div class="mt-2 space-y-2">
          <div
            v-for="request in detail.requests"
            :key="request.id"
            class="rounded border border-gray-200 bg-white p-3 text-sm"
          >
            <div class="flex flex-wrap items-center justify-between gap-2">
              <span class="font-medium">{{ LEAD_STATUS_LABELS[request.status] }}</span>
              <span class="text-xs text-[var(--color-text-muted)]">
                {{ new Date(request.createdAt).toLocaleString("ru-RU") }}
              </span>
            </div>
            <p class="mt-1 text-[var(--color-text-muted)]">
              {{ requestTypeLabels[request.requestType] ?? request.requestType }}
              <template v-if="request.relatedName || request.relatedSlug">
                · {{ request.relatedName || request.relatedSlug }}</template
              >
              · {{ sourceLabels[request.source] ?? request.source }} ·
              {{ contactMethodLabels[request.contactMethod] ?? request.contactMethod }}
            </p>
          </div>
          <p v-if="detail.requests.length === 0" class="text-sm text-[var(--color-text-muted)]">
            Заявок пока нет
          </p>
        </div>
      </section>

      <section class="mt-6">
        <h2 class="text-lg font-medium">Комментарии</h2>
        <form class="mt-2 flex gap-2" @submit.prevent="addComment">
          <input
            v-model="newComment"
            type="text"
            placeholder="Добавить комментарий…"
            class="flex-1 rounded border border-gray-300 px-3 py-2 text-sm"
          />
          <button
            type="submit"
            :disabled="commentSaving || !newComment.trim()"
            class="rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white disabled:opacity-50"
          >
            Добавить
          </button>
        </form>
        <ul class="mt-3 space-y-2">
          <li
            v-for="comment in detail.comments"
            :key="comment.id"
            class="rounded border border-gray-200 bg-white p-3 text-sm"
          >
            <p>{{ comment.text }}</p>
            <p class="mt-1 text-xs text-[var(--color-text-muted)]">
              {{ new Date(comment.createdAt).toLocaleString("ru-RU") }}
            </p>
          </li>
          <p v-if="detail.comments.length === 0" class="text-sm text-[var(--color-text-muted)]">
            Комментариев пока нет
          </p>
        </ul>
      </section>
    </template>
  </div>
</template>
