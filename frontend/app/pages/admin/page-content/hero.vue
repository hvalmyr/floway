<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

/**
 * Every page that has a Hero (see components/sections/Hero.vue) stores its
 * title/lead/image as three page_content keys following the same
 * `<page>_hero_<field>` naming. Grouping them here as one block per page —
 * a list on the left, one field-form on the right — instead of a single
 * flat list of 9 fields keeps it clear which title/lead/image go together
 * once there are more than one or two pages with a Hero.
 */
const heroBlocks = [
  {
    id: "home",
    label: "Главная",
    titleKey: "home_hero_title",
    leadKey: "home_hero_lead",
    imageKey: "home_hero_image",
  },
  {
    id: "masterclasses",
    label: "Мастер-классы",
    titleKey: "masterclasses_hero_title",
    leadKey: "masterclasses_hero_lead",
    imageKey: "masterclasses_hero_image",
  },
  {
    id: "gift_certificate",
    label: "Подарочные сертификаты",
    titleKey: "gift_certificate_hero_title",
    leadKey: "gift_certificate_hero_lead",
    imageKey: "gift_certificate_hero_image",
  },
];

const { items, loading, error, savingKey, savedKey, fetchAll, save, onImageUploaded } =
  useAdminPageContent(
    heroBlocks.flatMap((block) => [block.titleKey, block.leadKey, block.imageKey]),
  );

await fetchAll();

const selectedBlockId = ref(heroBlocks[0]!.id);
const selectedBlock = computed(() =>
  heroBlocks.find((block) => block.id === selectedBlockId.value)!,
);

function fieldFor(key: string) {
  return items.value.find((item) => item.key === key);
}

// The block is already named once in the list on the left, so repeating it
// in every field's own label ("Главная — заголовок Hero") is redundant —
// just "Заголовок"/"Описание" here. The underlying page_content `label`
// (still block-specific) stays as-is for the flat "все тексты сайта" list,
// which has no such grouping to lean on.
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
    <h1 class="mt-2 text-2xl font-semibold">Hero</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Заголовок, подзаголовок и фото в верхней части страницы — отдельный набор для каждой страницы
      с Hero-блоком.
    </p>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="error" class="mt-6 text-red-600">{{ error }}</p>
    <div v-else class="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-[200px_1fr]">
      <div class="flex flex-row gap-2 sm:flex-col">
        <button
          v-for="block in heroBlocks"
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
            <AdminMarkdownField
              :id="`field-${field.key}`"
              v-model="fieldFor(field.key)!.value"
              :rows="3"
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

        <div
          v-if="fieldFor(selectedBlock.imageKey)"
          class="rounded border border-gray-200 bg-white p-4"
        >
          <div class="mb-2 flex items-center justify-between gap-4">
            <span class="text-sm font-medium">Фотография</span>
            <span class="font-mono text-xs text-[var(--color-text-muted)]">{{
              selectedBlock.imageKey
            }}</span>
          </div>
          <AdminImageUpload
            :model-value="fieldFor(selectedBlock.imageKey)!.value"
            label="Фотография"
            @update:model-value="(url) => onImageUploaded(fieldFor(selectedBlock.imageKey)!, url)"
          />
        </div>
      </div>
    </div>
  </div>
</template>
