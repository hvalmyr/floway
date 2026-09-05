<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

const { items, loading, error, savingKey, savedKey, fetchAll, save } = useAdminPageContent([
  "apply_form_name_label",
  "apply_form_name_placeholder",
  "apply_form_phone_label",
  "apply_form_email_label",
  "apply_form_email_placeholder",
  "apply_form_contact_method_label",
  "apply_form_contact_method_call",
  "apply_form_contact_method_telegram",
  "apply_form_contact_method_whatsapp",
  "apply_form_contact_method_max",
  "apply_form_source_label",
  "apply_form_source_referral",
  "apply_form_source_ads",
  "apply_form_source_internet",
  "apply_form_source_social",
  "apply_form_source_maps",
  "apply_form_consent_prefix",
  "apply_form_consent_link_text",
  "apply_form_consent_middle",
  "apply_form_consent_link2_text",
  "apply_form_consent_suffix",
  "apply_form_submit_default",
  "apply_form_submit_trial",
  "apply_form_success_title",
  "apply_form_success_message",
]);

await fetchAll();
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Форма заявки</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Подписи полей и текст вариантов в селектах формы заявки (пробное занятие, курс, мастер-класс —
      форма везде одна и та же). Какие поля и варианты есть — фиксировано в коде, здесь можно менять
      только их текст.
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
        <input
          :id="`field-${item.key}`"
          v-model="item.value"
          type="text"
          class="w-full rounded border border-gray-300 px-3 py-2"
        />
        <div class="mt-2 flex items-center gap-3">
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
      </div>
    </div>
  </div>
</template>
