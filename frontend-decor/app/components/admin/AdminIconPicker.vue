<script setup lang="ts">
import { FEATURE_ICONS } from "~/constants/feature-icons";

/**
 * Icon field for admin forms — pick a built-in icon, pick one already
 * uploaded, or upload a new SVG on the spot. v-model is the same string
 * convention AppIcon.vue reads: a bare FEATURE_ICONS key or "icon:<id>".
 *
 * @example
 * <AdminIconPicker v-model="form.icon" />
 */
const modelValue = defineModel<string>({ required: true });

const { icons, findIcon, refresh } = useIconLibrary();
const api = useApiClient();

const uploading = ref(false);
const uploadError = ref("");
const fileInput = ref<HTMLInputElement | null>(null);

const selectedUploaded = computed(() =>
  modelValue.value.startsWith("icon:") ? findIcon(Number(modelValue.value.slice(5))) : undefined,
);

function pickFile() {
  fileInput.value?.click();
}

async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = "";
  if (!file) return;

  uploadError.value = "";
  uploading.value = true;
  try {
    const svg = await file.text();
    const name = file.name.replace(/\.svg$/i, "");
    const created = await api<{ id: number }>("/api/v1/icons", {
      method: "POST",
      body: { name, svg },
    });
    await refresh();
    modelValue.value = `icon:${created.id}`;
  } catch {
    uploadError.value = "Не удалось загрузить SVG — проверьте, что это корректный файл";
  } finally {
    uploading.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="flex flex-wrap items-center gap-2">
      <select v-model="modelValue" class="rounded border border-gray-300 px-3 py-2">
        <optgroup label="Встроенные">
          <option v-for="(def, key) in FEATURE_ICONS" :key="key" :value="key">
            {{ def.label }}
          </option>
        </optgroup>
        <optgroup v-if="icons.length" label="Загруженные">
          <option v-for="uploaded in icons" :key="uploaded.id" :value="`icon:${uploaded.id}`">
            {{ uploaded.name }}
          </option>
        </optgroup>
      </select>

      <span
        class="grid size-9 shrink-0 place-items-center rounded border border-gray-300 text-[var(--color-primary)]"
      >
        <AppIcon :icon="modelValue" class="size-6" />
      </span>

      <button
        type="button"
        :disabled="uploading"
        class="rounded border border-gray-300 px-3 py-2 text-sm disabled:opacity-50"
        @click="pickFile"
      >
        {{ uploading ? "Загрузка…" : "Загрузить SVG" }}
      </button>
    </div>

    <p v-if="selectedUploaded" class="text-xs text-[var(--color-text-muted)]">
      Загруженная иконка «{{ selectedUploaded.name }}». Для правильной перекраски под разные фоны
      SVG должен использовать fill/stroke="currentColor".
    </p>
    <p v-if="uploadError" class="text-xs text-red-600">{{ uploadError }}</p>

    <input
      ref="fileInput"
      type="file"
      accept=".svg,image/svg+xml"
      class="hidden"
      @change="onFileChange"
    />
  </div>
</template>
