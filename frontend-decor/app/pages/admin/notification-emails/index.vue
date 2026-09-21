<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface NotificationEmail {
  id: number;
  email: string;
}

const { items, loading, error, fetchAll, create, remove } = useAdminResource<NotificationEmail>(
  "/api/v1/notification-emails",
);

const newEmail = ref("");
const saving = ref(false);
const formError = ref("");

await fetchAll();

async function onSubmit() {
  formError.value = "";
  saving.value = true;
  try {
    await create({ email: newEmail.value });
    newEmail.value = "";
  } catch {
    formError.value = "Не удалось добавить. Проверьте адрес и попробуйте снова.";
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Убрать этот адрес из списка уведомлений?")) return;
  await remove(id);
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Email для уведомлений о заявках</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Письмо о новой заявке уходит на все адреса из этого списка. Добавьте или уберите адрес — без
      деплоя.
    </p>

    <form
      class="mt-6 flex flex-col gap-3 rounded border border-gray-200 bg-white p-4 sm:flex-row"
      @submit.prevent="onSubmit"
    >
      <input
        v-model="newEmail"
        type="email"
        placeholder="you@example.com"
        required
        class="flex-1 rounded border border-gray-300 px-3 py-2"
      />
      <button
        type="submit"
        :disabled="saving"
        class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
      >
        Добавить
      </button>
    </form>
    <p v-if="formError" class="mt-2 text-sm text-red-600">{{ formError }}</p>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <table
      v-else
      class="mt-6 w-full border-collapse overflow-hidden rounded border border-gray-200 bg-white text-sm"
    >
      <thead>
        <tr class="border-b border-gray-200 text-left">
          <th class="px-4 py-2">Email</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" class="border-b border-gray-100 last:border-0">
          <td class="px-4 py-2">{{ item.email }}</td>
          <td class="px-4 py-2 text-right">
            <button class="text-red-600 hover:underline" @click="onDelete(item.id)">Убрать</button>
          </td>
        </tr>
        <tr v-if="items.length === 0">
          <td colspan="2" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока пусто — уведомления никуда не отправляются
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
