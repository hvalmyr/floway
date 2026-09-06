<script setup lang="ts">
import type { FaqPage } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

interface FAQItem {
  id: number;
  question: string;
  answer: string;
  sortOrder: number;
}

const emptyForm = (): Omit<FAQItem, "id"> => ({ question: "", answer: "", sortOrder: 0 });

const { items, loading, error, fetchAll, create, update, remove } =
  useAdminResource<FAQItem>("/api/v1/faq");

const editingId = ref<number | null>(null);
const form = ref(emptyForm());
const saving = ref(false);
const formError = ref("");

const { draggingIndex, onPointerDown } = useAdminDragReorder(items, (item) =>
  update(item.id, item),
);

function startEdit(faqItem: FAQItem) {
  editingId.value = faqItem.id;
  form.value = {
    question: faqItem.question,
    answer: faqItem.answer,
    sortOrder: faqItem.sortOrder,
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
  if (!confirm("Удалить вопрос?")) return;
  await remove(id);
  if (editingId.value === id) cancelEdit();
}

/**
 * Masterclasses/gift-certificate each get their own FAQ block — title,
 * intro text, visible toggle and their own Q&A list — separate from the
 * flat, unscoped homepage list above (see useAdminPageFaq's doc comment).
 * One useAdminPageFaq instance per page (not a single dynamic one) since
 * there are only ever these two known pages.
 */
const pageBlocks: { id: FaqPage; label: string }[] = [
  { id: "masterclasses", label: "Мастер-классы" },
  { id: "gift_certificate", label: "Подарочные сертификаты" },
];
const pageFaqs = {
  masterclasses: useAdminPageFaq("masterclasses"),
  gift_certificate: useAdminPageFaq("gift_certificate"),
};

type Tab = "home" | FaqPage;
const activeTab = ref<Tab>("home");
const activePageFaq = computed(() =>
  activeTab.value === "home" ? null : pageFaqs[activeTab.value],
);

interface PageFaqItemForm {
  id: number | null;
  question: string;
  answer: string;
  sortOrder: number;
}
const pageItemForm = ref<PageFaqItemForm>({ id: null, question: "", answer: "", sortOrder: 0 });
const pageItemSaving = ref(false);
const pageItemFormError = ref("");

function selectTab(tab: Tab) {
  activeTab.value = tab;
  cancelEdit();
  cancelPageItemEdit();
}

function startPageItemEdit(item: {
  id: number;
  question: string;
  answer: string;
  sortOrder: number;
}) {
  pageItemForm.value = { ...item };
}

function cancelPageItemEdit() {
  pageItemForm.value = { id: null, question: "", answer: "", sortOrder: 0 };
}

async function onPageItemSubmit() {
  const faq = activePageFaq.value;
  if (!faq) return;
  pageItemFormError.value = "";
  pageItemSaving.value = true;
  try {
    if (pageItemForm.value.id === null) {
      await faq.create({
        question: pageItemForm.value.question,
        answer: pageItemForm.value.answer,
        sortOrder: faq.items.value.length,
      });
    } else {
      await faq.update(pageItemForm.value.id, {
        question: pageItemForm.value.question,
        answer: pageItemForm.value.answer,
        sortOrder: pageItemForm.value.sortOrder,
      });
    }
    cancelPageItemEdit();
  } catch {
    pageItemFormError.value = "Не удалось сохранить. Проверьте поля и попробуйте снова.";
  } finally {
    pageItemSaving.value = false;
  }
}

async function onPageItemDelete(id: number) {
  const faq = activePageFaq.value;
  if (!faq) return;
  if (!confirm("Удалить вопрос?")) return;
  await faq.remove(id);
  if (pageItemForm.value.id === id) cancelPageItemEdit();
}

const pageFaqReorders = {
  masterclasses: useAdminDragReorder(pageFaqs.masterclasses.items, (item) =>
    pageFaqs.masterclasses.update(item.id, item),
  ),
  gift_certificate: useAdminDragReorder(pageFaqs.gift_certificate.items, (item) =>
    pageFaqs.gift_certificate.update(item.id, item),
  ),
};
const activeReorder = computed(() =>
  activeTab.value === "home" ? null : pageFaqReorders[activeTab.value],
);

await Promise.all([
  fetchAll(),
  pageFaqs.masterclasses.fetchAll(),
  pageFaqs.gift_certificate.fetchAll(),
]);
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">FAQ</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      «Главная» — общий блок вопросов-ответов на главной странице. У мастер-классов и подарочных
      сертификатов — свой заголовок, текст и список вопросов, с отдельным переключателем видимости.
    </p>

    <div class="mt-6 flex flex-wrap gap-2">
      <button
        type="button"
        class="rounded border px-3 py-2 text-sm"
        :class="
          activeTab === 'home'
            ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/10 font-medium'
            : 'border-gray-200 bg-white hover:border-[var(--color-primary)]'
        "
        @click="selectTab('home')"
      >
        Главная
      </button>
      <button
        v-for="block in pageBlocks"
        :key="block.id"
        type="button"
        class="rounded border px-3 py-2 text-sm"
        :class="
          activeTab === block.id
            ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/10 font-medium'
            : 'border-gray-200 bg-white hover:border-[var(--color-primary)]'
        "
        @click="selectTab(block.id)"
      >
        {{ block.label }}
      </button>
    </div>

    <template v-if="activeTab === 'home'">
      <form
        class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
        @submit.prevent="onSubmit"
      >
        <AdminMarkdownField
          v-model="form.question"
          placeholder="Вопрос"
          :rows="2"
          required
          class="sm:col-span-2"
        />
        <AdminMarkdownField
          v-model="form.answer"
          placeholder="Ответ"
          :rows="3"
          required
          class="sm:col-span-2"
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
            <th class="px-4 py-2">Вопрос</th>
            <th class="px-4 py-2" />
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(faqItem, index) in items"
            :key="faqItem.id"
            :data-row-index="index"
            class="border-b border-gray-100 last:border-0"
            :class="draggingIndex === index ? 'opacity-50' : ''"
          >
            <td class="px-4 py-2">
              <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
            </td>
            <td class="px-4 py-2">{{ faqItem.question }}</td>
            <td class="flex gap-3 px-4 py-2 text-right">
              <button
                class="text-[var(--color-primary)] hover:underline"
                @click="startEdit(faqItem)"
              >
                Редактировать
              </button>
              <button class="text-red-600 hover:underline" @click="onDelete(faqItem.id)">
                Удалить
              </button>
            </td>
          </tr>
          <tr v-if="items.length === 0">
            <td colspan="3" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
              Пока пусто
            </td>
          </tr>
        </tbody>
      </table>
    </template>

    <template v-else-if="activePageFaq">
      <p v-if="activePageFaq.loading.value" class="mt-6 text-[var(--color-text-muted)]">
        Загрузка…
      </p>
      <p v-else-if="activePageFaq.error.value" class="mt-6 text-red-600">
        {{ activePageFaq.error.value }}
      </p>
      <div v-else class="mt-6 flex flex-col gap-6">
        <form
          class="grid gap-3 rounded border border-gray-200 bg-white p-4"
          @submit.prevent="activePageFaq.saveSettings"
        >
          <input
            v-model="activePageFaq.settings.value.title"
            type="text"
            placeholder="Заголовок FAQ (например, «Вопросы и ответы»)"
            class="rounded border border-gray-300 px-3 py-2"
          />
          <AdminMarkdownField
            v-model="activePageFaq.settings.value.description"
            placeholder="Текст перед вопросами"
            :rows="2"
          />
          <label class="flex items-center gap-2 text-sm">
            <input v-model="activePageFaq.settings.value.visible" type="checkbox" class="size-5" />
            Показывать FAQ на странице
          </label>
          <div>
            <button
              type="submit"
              :disabled="activePageFaq.saving.value"
              class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
            >
              Сохранить настройки
            </button>
          </div>
        </form>

        <form
          class="grid gap-3 rounded border border-gray-200 bg-white p-4 sm:grid-cols-2"
          @submit.prevent="onPageItemSubmit"
        >
          <AdminMarkdownField
            v-model="pageItemForm.question"
            placeholder="Вопрос"
            :rows="2"
            required
            class="sm:col-span-2"
          />
          <AdminMarkdownField
            v-model="pageItemForm.answer"
            placeholder="Ответ"
            :rows="3"
            required
            class="sm:col-span-2"
          />
          <p v-if="pageItemFormError" class="text-sm text-red-600 sm:col-span-2">
            {{ pageItemFormError }}
          </p>

          <div class="flex gap-2 sm:col-span-2">
            <button
              type="submit"
              :disabled="pageItemSaving"
              class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
            >
              {{ pageItemForm.id === null ? "Добавить" : "Сохранить" }}
            </button>
            <button
              v-if="pageItemForm.id !== null"
              type="button"
              class="rounded border border-gray-300 px-4 py-2"
              @click="cancelPageItemEdit"
            >
              Отмена
            </button>
          </div>
        </form>

        <table
          class="w-full border-collapse overflow-hidden rounded border border-gray-200 bg-white text-sm"
        >
          <thead>
            <tr class="border-b border-gray-200 text-left">
              <th class="px-4 py-2" />
              <th class="px-4 py-2">Вопрос</th>
              <th class="px-4 py-2" />
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, index) in activePageFaq.items.value"
              :key="item.id"
              :data-row-index="index"
              class="border-b border-gray-100 last:border-0"
              :class="activeReorder?.draggingIndex.value === index ? 'opacity-50' : ''"
            >
              <td class="px-4 py-2">
                <AdminDragHandle @pointerdown="activeReorder?.onPointerDown(index, $event)" />
              </td>
              <td class="px-4 py-2">{{ item.question }}</td>
              <td class="flex gap-3 px-4 py-2 text-right">
                <button
                  class="text-[var(--color-primary)] hover:underline"
                  @click="startPageItemEdit(item)"
                >
                  Редактировать
                </button>
                <button class="text-red-600 hover:underline" @click="onPageItemDelete(item.id)">
                  Удалить
                </button>
              </td>
            </tr>
            <tr v-if="activePageFaq.items.value.length === 0">
              <td colspan="3" class="px-4 py-6 text-center text-[var(--color-text-muted)]">
                Пока пусто
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>
