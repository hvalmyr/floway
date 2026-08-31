<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

const { items, loading, error, savingKey, savedKey, fetchAll, save, onImageUploaded } =
  useAdminPageContent([
    "trial_section_heading",
    "trial_section_description",
    "trial_heading",
    "trial_description",
    "home_trial_image",
  ]);

await fetchAll();
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Пробное занятие</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Блок пробного занятия на главной странице.
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
          <label :for="`field-${item.key}`" class="text-sm font-medium">{{ item.label }}</label>
          <span class="font-mono text-xs text-[var(--color-text-muted)]">{{ item.key }}</span>
        </div>
        <AdminImageUpload
          v-if="item.type === 'image'"
          :model-value="item.value"
          :label="item.label"
          @update:model-value="(url) => onImageUploaded(item, url)"
        />
        <AdminMarkdownField v-else :id="`field-${item.key}`" v-model="item.value" :rows="3" />
        <div class="mt-2 flex items-center gap-3">
          <button
            v-if="item.type !== 'image'"
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
    </div>

    <div class="mt-8 flex flex-col gap-2 rounded border border-gray-200 bg-white p-4">
      <p class="text-sm font-medium">Другие блоки главной страницы</p>
      <NuxtLink to="/admin/teachers" class="text-sm text-[var(--color-primary)] hover:underline"
        >Преподаватели</NuxtLink
      >
      <NuxtLink to="/admin/about-items" class="text-sm text-[var(--color-primary)] hover:underline"
        >О школе</NuxtLink
      >
    </div>
  </div>
</template>
