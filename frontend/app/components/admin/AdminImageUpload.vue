<script setup lang="ts">
/**
 * File-picker replacement for the old "paste an image URL" text inputs in
 * the admin forms. Uploads straight to the backend (Garage-backed) and
 * writes the returned relative URL into v-model.
 *
 * @example
 * <AdminImageUpload v-model="form.coverImage" label="Обложка" />
 */
const props = withDefaults(
  defineProps<{
    modelValue: string;
    label?: string;
  }>(),
  { label: "Изображение" },
);

const emit = defineEmits<{ "update:modelValue": [value: string] }>();

const { upload, uploading, error } = useAdminUpload();
const fileInput = ref<HTMLInputElement | null>(null);

const previewSrc = computed(() => (props.modelValue ? resolveMediaUrl(props.modelValue) : null));

function pickFile() {
  fileInput.value?.click();
}

async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = "";
  if (!file) return;

  try {
    const url = await upload(file);
    emit("update:modelValue", url);
  } catch {
    // error already captured in useAdminUpload's error ref
  }
}

function clear() {
  emit("update:modelValue", "");
}
</script>

<template>
  <div class="flex items-center gap-3 rounded border border-gray-300 px-3 py-2">
    <img v-if="previewSrc" :src="previewSrc" alt="" class="size-12 shrink-0 rounded object-cover" />
    <div class="flex min-w-0 flex-1 flex-col gap-1">
      <span class="truncate text-sm text-[var(--color-text-muted)]">
        {{ modelValue ? label : `${label}: не загружено` }}
      </span>
      <p v-if="error" class="text-xs text-red-600">{{ error }}</p>
    </div>
    <div class="flex shrink-0 gap-2">
      <button
        type="button"
        :disabled="uploading"
        class="rounded border border-gray-300 px-2 py-1 text-sm disabled:opacity-50"
        @click="pickFile"
      >
        {{ uploading ? "Загрузка…" : modelValue ? "Заменить" : "Загрузить" }}
      </button>
      <button
        v-if="modelValue"
        type="button"
        class="rounded border border-gray-300 px-2 py-1 text-sm text-red-600"
        @click="clear"
      >
        Удалить
      </button>
    </div>
    <input
      ref="fileInput"
      type="file"
      accept="image/jpeg,image/png,image/webp,image/gif"
      class="hidden"
      @change="onFileChange"
    />
  </div>
</template>
