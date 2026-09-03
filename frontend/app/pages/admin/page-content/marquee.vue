<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

const { items, loading, error, savingKey, savedKey, fetchAll, save, onIconSelected } =
  useAdminPageContent(["gift_certificate_marquee_text", "gift_certificate_marquee_icon"]);

await fetchAll();
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Бегущая строка сертификатов</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Синяя полоса с текстом и иконкой на главной странице, ведущая на страницу сертификатов.
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
        <AdminIconPicker
          v-if="item.type === 'icon'"
          :model-value="item.value"
          @update:model-value="(value) => onIconSelected(item, value)"
        />
        <input
          v-else
          :id="`field-${item.key}`"
          v-model="item.value"
          type="text"
          class="w-full rounded border border-gray-300 px-3 py-2"
        />
        <div v-if="item.type !== 'icon'" class="mt-2 flex items-center gap-3">
          <button
            type="button"
            :disabled="savingKey === item.key"
            class="rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white disabled:opacity-50"
            @click="save(item)"
          >
            {{ savingKey === item.key ? "Сохранение…" : "Сохранить" }}
          </button>
          <span v-if="savedKey === item.key" class="text-sm text-green-600">Сохранено</span>
        </div>
        <span v-else-if="savedKey === item.key" class="mt-2 block text-sm text-green-600"
          >Сохранено</span
        >
      </div>
    </div>
  </div>
</template>
