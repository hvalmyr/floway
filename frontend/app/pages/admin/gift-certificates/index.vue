<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

type GiftCertificateKind = "amount" | "course" | "masterclass" | "any_masterclass";

interface GiftCertificate {
  id: number;
  number: string;
  issuedAt: string;
  kind: GiftCertificateKind;
  value: string;
  recipient: string;
}

const kindOptions: { value: GiftCertificateKind; label: string; valuePlaceholder?: string }[] = [
  { value: "amount", label: "Сумма", valuePlaceholder: "5000 рублей" },
  { value: "course", label: "Курс", valuePlaceholder: "основы флористики" },
  { value: "masterclass", label: "Мастер-класс", valuePlaceholder: "раскидистый букет" },
  { value: "any_masterclass", label: "Любой мастер-класс" },
];
const kindLabels = Object.fromEntries(kindOptions.map((o) => [o.value, o.label])) as Record<
  GiftCertificateKind,
  string
>;

const { items, loading, error, fetchAll, create, remove } = useAdminResource<GiftCertificate>(
  "/api/v1/gift-certificates",
);

const config = useRuntimeConfig();
function pdfUrl(id: number) {
  return `${config.public.apiBase}/api/v1/gift-certificates/${id}/pdf`;
}

const kind = ref<GiftCertificateKind>("amount");
const value = ref("");
const recipient = ref("");
const saving = ref(false);
const formError = ref("");

const selectedKind = computed(() => kindOptions.find((o) => o.value === kind.value));

await fetchAll();

async function onSubmit() {
  formError.value = "";
  saving.value = true;
  try {
    await create({
      kind: kind.value,
      value: kind.value === "any_masterclass" ? "" : value.value,
      recipient: recipient.value,
    });
    value.value = "";
    recipient.value = "";
  } catch {
    formError.value = "Не удалось создать сертификат. Проверьте поля и попробуйте снова.";
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Удалить этот сертификат из списка?")) return;
  await remove(id);
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("ru-RU");
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Подарочные сертификаты</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Выпустите сертификат — номер и дата проставятся автоматически, PDF в дизайне школы можно сразу
      скачать и отправить получателю.
    </p>

    <form
      class="mt-6 flex flex-col gap-3 rounded border border-gray-200 bg-white p-4"
      @submit.prevent="onSubmit"
    >
      <div class="flex flex-col gap-3 sm:flex-row">
        <select v-model="kind" class="rounded border border-gray-300 px-3 py-2 sm:w-56">
          <option v-for="opt in kindOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <input
          v-if="kind !== 'any_masterclass'"
          v-model="value"
          type="text"
          :placeholder="selectedKind?.valuePlaceholder"
          required
          class="flex-1 rounded border border-gray-300 px-3 py-2"
        />
        <input
          v-model="recipient"
          type="text"
          placeholder="Васильевой Василисе Васильевне"
          required
          class="flex-1 rounded border border-gray-300 px-3 py-2"
        />
      </div>
      <div>
        <button
          type="submit"
          :disabled="saving"
          class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
        >
          Выпустить сертификат
        </button>
      </div>
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
          <th class="px-4 py-2">Номер</th>
          <th class="px-4 py-2">Дата</th>
          <th class="px-4 py-2">Тип</th>
          <th class="px-4 py-2">На что</th>
          <th class="px-4 py-2">Получатель</th>
          <th class="px-4 py-2" />
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" class="border-b border-gray-100 last:border-0">
          <td class="px-4 py-2">{{ item.number }}</td>
          <td class="px-4 py-2">{{ formatDate(item.issuedAt) }}</td>
          <td class="px-4 py-2">{{ kindLabels[item.kind] }}</td>
          <td class="px-4 py-2">{{ item.value }}</td>
          <td class="px-4 py-2">{{ item.recipient }}</td>
          <td class="px-4 py-2 text-right whitespace-nowrap">
            <a
              :href="pdfUrl(item.id)"
              target="_blank"
              class="mr-4 text-[var(--color-primary)] hover:underline"
              >Скачать PDF</a
            >
            <button class="text-red-600 hover:underline" @click="onDelete(item.id)">Удалить</button>
          </td>
        </tr>
        <tr v-if="items.length === 0">
          <td colspan="6" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
            Пока ни одного сертификата не выпущено
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
