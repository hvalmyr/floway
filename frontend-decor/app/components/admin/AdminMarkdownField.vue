<script setup lang="ts">
import { Bold, Image as ImageIcon, Italic, Link as LinkIcon } from "lucide-vue-next";
import { nextTick, ref } from "vue";

/**
 * Drop-in replacement for a plain <textarea v-model="..."> admin field —
 * same placeholder-as-label look, but adds a formatting toolbar (bold,
 * italic, link) and an image button that uploads a file (via
 * useAdminUpload, same endpoint AdminImageUpload.vue uses) and inserts a
 * markdown image reference. The inserted URL is resolveMediaUrl()'d (not
 * the bare relative path useAdminUpload returns) because MarkdownContent.vue
 * renders the raw markdown source's URLs as-is, with no resolution step —
 * only the browser knows the API origin, and only at render time here.
 *
 * @example
 * <AdminMarkdownField v-model="form.description" placeholder="Описание *" :rows="3" required />
 */
const props = withDefaults(
  defineProps<{
    modelValue: string;
    id?: string;
    placeholder?: string;
    rows?: number;
    required?: boolean;
  }>(),
  { id: undefined, placeholder: "", rows: 4, required: false },
);
const emit = defineEmits<{ "update:modelValue": [value: string] }>();

const textareaRef = ref<HTMLTextAreaElement | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);
const { upload, uploading, error } = useAdminUpload();

function placeCursor(pos: number) {
  const el = textareaRef.value;
  if (!el) return;
  nextTick(() => {
    el.focus();
    el.setSelectionRange(pos, pos);
  });
}

// Wraps the current selection (or a placeholder word, if nothing's
// selected) in `before`/`after` — used for bold/italic.
function surroundSelection(before: string, after: string, placeholderText: string) {
  const el = textareaRef.value;
  if (!el) return;
  const start = el.selectionStart;
  const end = el.selectionEnd;
  const selected = props.modelValue.slice(start, end) || placeholderText;
  const next =
    props.modelValue.slice(0, start) + before + selected + after + props.modelValue.slice(end);
  emit("update:modelValue", next);
  placeCursor(start + before.length + selected.length + after.length);
}

function insertAtCursor(text: string) {
  const el = textareaRef.value;
  if (!el) return;
  const start = el.selectionStart;
  const end = el.selectionEnd;
  const next = props.modelValue.slice(0, start) + text + props.modelValue.slice(end);
  emit("update:modelValue", next);
  placeCursor(start + text.length);
}

function onBold() {
  surroundSelection("**", "**", "текст");
}

function onItalic() {
  surroundSelection("*", "*", "текст");
}

function onLink() {
  const url = window.prompt("Ссылка (URL):");
  if (!url) return;
  const el = textareaRef.value;
  const selected = el ? props.modelValue.slice(el.selectionStart, el.selectionEnd) : "";
  insertAtCursor(`[${selected || "текст ссылки"}](${url})`);
}

function pickImage() {
  fileInput.value?.click();
}

async function onImageChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = "";
  if (!file) return;
  try {
    const relativeUrl = await upload(file);
    insertAtCursor(`![](${resolveMediaUrl(relativeUrl)})`);
  } catch {
    // upload() already captured the failure in `error` below.
  }
}

// Ctrl/Cmd+B/I/K mirror the toolbar buttons; image is Ctrl/Cmd+Shift+I
// (checked before the plain-I italic case) — Ctrl+P was the original pick
// but collided with the browser's print shortcut, and Ctrl+K often focuses
// the address bar, so preventDefault only fires once we've matched one of
// these, leaving every other shortcut (copy/paste, undo, tab switching,
// ...) reaching the browser exactly as it would from a plain textarea.
function onKeydown(event: KeyboardEvent) {
  if (!(event.ctrlKey || event.metaKey)) return;
  const key = event.key.toLowerCase();
  if (key === "i" && event.shiftKey) {
    event.preventDefault();
    pickImage();
    return;
  }
  switch (key) {
    case "b":
      event.preventDefault();
      onBold();
      break;
    case "i":
      event.preventDefault();
      onItalic();
      break;
    case "k":
      event.preventDefault();
      onLink();
      break;
  }
}
</script>

<template>
  <div class="flex flex-col rounded border border-gray-300">
    <div class="flex items-center gap-1 border-b border-gray-200 bg-gray-50 px-2 py-1">
      <button
        type="button"
        class="rounded p-1.5 hover:bg-gray-200"
        title="Жирный (Ctrl+B)"
        @click="onBold"
      >
        <Bold class="size-5" />
      </button>
      <button
        type="button"
        class="rounded p-1.5 hover:bg-gray-200"
        title="Курсив (Ctrl+I)"
        @click="onItalic"
      >
        <Italic class="size-5" />
      </button>
      <button
        type="button"
        class="rounded p-1.5 hover:bg-gray-200"
        title="Ссылка (Ctrl+K)"
        @click="onLink"
      >
        <LinkIcon class="size-5" />
      </button>
      <button
        type="button"
        class="rounded p-1.5 hover:bg-gray-200 disabled:opacity-50"
        title="Картинка (Ctrl+Shift+I)"
        :disabled="uploading"
        @click="pickImage"
      >
        <ImageIcon class="size-5" />
      </button>
      <span class="ml-auto text-xs text-[var(--color-text-muted)]">Markdown</span>
    </div>
    <textarea
      :id="id"
      ref="textareaRef"
      :value="modelValue"
      :placeholder="placeholder"
      :rows="rows"
      :required="required"
      class="w-full resize-y rounded-b px-3 py-2 outline-none"
      @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
      @keydown="onKeydown"
    />
    <p v-if="error" class="border-t border-gray-200 px-3 py-1 text-xs text-red-600">{{ error }}</p>
    <input
      ref="fileInput"
      type="file"
      accept="image/jpeg,image/png,image/webp,image/gif"
      class="hidden"
      @change="onImageChange"
    />
  </div>
</template>
