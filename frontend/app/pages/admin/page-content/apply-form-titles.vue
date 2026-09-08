<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

/**
 * The shared lead-capture form (ApplyForm.vue) is embedded on three pages,
 * each with its own title/subhead wired in as `<page>_apply_form_<field>`
 * page_content keys — same grouping approach as page-content/hero.vue.
 * The trial-lesson embed on the home page isn't listed here: it renders
 * with an empty title on purpose (its section already has its own heading).
 */
const formBlocks = [
  {
    id: "course",
    label: "Страница курса",
    titleKey: "course_apply_form_title",
    leadKey: "course_apply_form_lead",
  },
  {
    id: "masterclasses",
    label: "Мастер-классы",
    titleKey: "masterclasses_apply_form_title",
    leadKey: "masterclasses_apply_form_lead",
  },
  {
    id: "gift_certificate",
    label: "Подарочные сертификаты",
    titleKey: "gift_certificate_apply_form_title",
    leadKey: "gift_certificate_apply_form_lead",
  },
];

const { items, loading, error, savingKey, savedKey, fetchAll, save } = useAdminPageContent(
  formBlocks.flatMap((block) => [block.titleKey, block.leadKey]),
);

await fetchAll();

const selectedBlockId = ref(formBlocks[0]!.id);
const selectedBlock = computed(() =>
  formBlocks.find((block) => block.id === selectedBlockId.value)!,
);

function fieldFor(key: string) {
  return items.value.find((item) => item.key === key);
}

const titleAndLeadFields = computed(() => [
  { key: selectedBlock.value.titleKey, caption: "Заголовок" },
  { key: selectedBlock.value.leadKey, caption: "Описание" },
]);
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Заголовки формы заявки</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Заголовок и описание над формой заявки — отдельный набор для каждой страницы, где форма
      встречается. Какие поля есть в самой форме — см. «Форма заявки».
    </p>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-[200px_1fr]">
      <div class="flex flex-row gap-2 sm:flex-col">
        <button
          v-for="block in formBlocks"
          :key="block.id"
          type="button"
          class="rounded border px-3 py-2 text-left text-sm"
          :class="
            block.id === selectedBlockId
              ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/10 font-medium'
              : 'border-gray-200 bg-white hover:border-[var(--color-primary)]'
          "
          @click="selectedBlockId = block.id"
        >
          {{ block.label }}
        </button>
      </div>

      <div class="flex flex-col gap-4">
        <template v-for="field in titleAndLeadFields" :key="field.key">
          <div v-if="fieldFor(field.key)" class="rounded border border-gray-200 bg-white p-4">
            <div class="mb-2 flex items-center justify-between gap-4">
              <label :for="`field-${field.key}`" class="text-sm font-medium">{{
                field.caption
              }}</label>
              <span class="font-mono text-xs text-[var(--color-text-muted)]">{{ field.key }}</span>
            </div>
            <input
              :id="`field-${field.key}`"
              v-model="fieldFor(field.key)!.value"
              type="text"
              class="w-full rounded border border-gray-300 px-3 py-2"
            />
            <div class="mt-2 flex items-center gap-3">
              <button
                type="button"
                :disabled="savingKey === field.key"
                class="rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white disabled:opacity-50"
                @click="save(fieldFor(field.key)!)"
              >
                {{ savingKey === field.key ? "Сохранение…" : "Сохранить" }}
              </button>
              <span v-if="savedKey === field.key" class="text-sm text-green-600">Сохранено</span>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
