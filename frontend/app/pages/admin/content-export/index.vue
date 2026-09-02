<script setup lang="ts">
import type { ImportMode, ImportResult } from "~/types/api";

definePageMeta({ layout: "admin", middleware: "admin-auth" });

const config = useRuntimeConfig();
// Browser-side only (this page always renders client-side after admin
// auth) — same public API base the rest of the admin panel's fetch calls
// use, see useApiClient.ts.
const exportUrl = `${config.public.apiBase}/api/v1/admin/content/export`;

const api = useApiClient();

const fileInput = ref<HTMLInputElement | null>(null);
const fileName = ref("");
// Parsed once at file-select time (not re-read at submit time) so a bad/
// non-JSON file surfaces its error immediately, next to the file picker,
// rather than after the admin has already committed to a mode and clicked
// submit.
const parsedData = ref<unknown>(null);
const parseError = ref("");

const mode = ref<ImportMode>("merge");
const confirmReplace = ref(false);

const importing = ref(false);
const importResult = ref<ImportResult | null>(null);
const importError = ref("");

async function onFileChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0];
  importResult.value = null;
  importError.value = "";
  parsedData.value = null;
  parseError.value = "";
  fileName.value = file?.name ?? "";
  if (!file) return;

  try {
    const text = await file.text();
    parsedData.value = JSON.parse(text);
  } catch {
    parseError.value = "Не удалось прочитать файл — это не валидный JSON";
  }
}

async function onImport() {
  if (!parsedData.value || (mode.value === "replace" && !confirmReplace.value)) return;

  importing.value = true;
  importResult.value = null;
  importError.value = "";
  try {
    importResult.value = await api<ImportResult>("/api/v1/admin/content/import", {
      method: "POST",
      body: { mode: mode.value, data: parsedData.value },
    });
    // Force re-selecting a file (and re-confirming replace) for the next
    // import rather than letting a second click resubmit the same payload.
    fileName.value = "";
    parsedData.value = null;
    confirmReplace.value = false;
    if (fileInput.value) fileInput.value.value = "";
  } catch (err) {
    importError.value =
      (err as { data?: { error?: string } } | undefined)?.data?.error ??
      "Не удалось выполнить импорт";
  } finally {
    importing.value = false;
  }
}

const countEntries = computed(() => Object.entries(importResult.value?.counts ?? {}));
</script>

<template>
  <div>
    <h1 class="text-2xl font-semibold">Экспорт и импорт контента</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Весь контент сайта целиком: курсы, мастер-классы, преподаватели, фотогалерея, FAQ, тексты
      страниц и все загруженные картинки — в одном файле. Аккаунты админов и заявки клиентов сюда не
      входят.
    </p>

    <section class="mt-8 rounded border border-gray-200 bg-white p-4">
      <h2 class="font-medium">Экспорт</h2>
      <p class="mt-1 text-sm text-[var(--color-text-muted)]">
        Скачивает весь текущий контент одним JSON-файлом.
      </p>
      <a
        :href="exportUrl"
        class="mt-3 inline-block rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white"
      >
        Скачать экспорт
      </a>
    </section>

    <section class="mt-6 rounded border border-gray-200 bg-white p-4">
      <h2 class="font-medium">Импорт</h2>

      <div class="mt-3 flex flex-col gap-3">
        <input ref="fileInput" type="file" accept="application/json" @change="onFileChange" />
        <p v-if="parseError" class="text-sm text-red-600">{{ parseError }}</p>

        <fieldset class="flex flex-col gap-2 text-sm">
          <legend class="mb-1 font-medium">Режим импорта</legend>
          <label class="flex items-start gap-2">
            <input v-model="mode" type="radio" value="merge" class="mt-1" />
            <span>
              <strong>Безопасное объединение</strong> — существующие записи обновляются по id, новые
              добавляются, ничего не удаляется.
            </span>
          </label>
          <label class="flex items-start gap-2">
            <input v-model="mode" type="radio" value="replace" class="mt-1" />
            <span>
              <strong>Полная замена</strong> — весь текущий контент удаляется, вместо него
              загружается файл целиком. Необратимо.
            </span>
          </label>
        </fieldset>

        <label v-if="mode === 'replace'" class="flex items-start gap-2 text-sm text-red-600">
          <input v-model="confirmReplace" type="checkbox" class="mt-1" />
          <span>Понимаю, что весь текущий контент сайта будет удалён и заменён файлом.</span>
        </label>

        <button
          type="button"
          class="self-start rounded bg-[var(--color-primary)] px-4 py-2 text-sm text-white disabled:opacity-50"
          :disabled="importing || !parsedData || (mode === 'replace' && !confirmReplace)"
          @click="onImport"
        >
          {{ importing ? "Импортирую…" : "Импортировать" }}
        </button>

        <p v-if="importError" class="text-sm text-red-600">{{ importError }}</p>

        <div v-if="importResult" class="rounded bg-green-50 p-3 text-sm">
          <p class="font-medium text-green-800">
            Готово ({{ importResult.mode === "merge" ? "объединение" : "полная замена" }})
          </p>
          <ul class="mt-1 text-[var(--color-text-muted)]">
            <li v-for="[key, count] in countEntries" :key="key">{{ key }}: {{ count }}</li>
          </ul>
          <p v-if="importResult.pageContentSkipped?.length" class="mt-2 text-amber-700">
            Пропущено (нет на этом сервере): {{ importResult.pageContentSkipped.join(", ") }}
          </p>
        </div>
      </div>
    </section>
  </div>
</template>
