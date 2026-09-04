<script setup lang="ts">
import {
  Bold,
  Heading2,
  Image as ImageIcon,
  Italic,
  Link as LinkIcon,
  Quote,
} from "lucide-vue-next";
import { onMounted, ref } from "vue";
import { sanitizeRichTextHtml } from "~/lib/richTextSanitize";

/**
 * Minimalist WYSIWYG editor for long-form content (blog posts) — a
 * telegra.ph-style single contenteditable region instead of a raw markdown
 * textarea. Supports paragraphs, one heading level, quotes, bold/italic,
 * links, and inline images. Output is HTML, sanitized against a fixed
 * whitelist (see richTextSanitize.ts) on every input and especially on
 * paste, since pasted content can carry arbitrary markup from Word/web
 * pages that execCommand-driven typing never produces on its own.
 *
 * @example
 * <AdminRichTextEditor v-model="form.content" placeholder="Текст статьи" />
 */
const props = withDefaults(
  defineProps<{
    modelValue: string;
    placeholder?: string;
  }>(),
  { placeholder: "Текст статьи…" },
);
const emit = defineEmits<{ "update:modelValue": [value: string] }>();

const editorRef = ref<HTMLDivElement | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);
const { upload, uploading, error } = useAdminUpload();

// Tracks our own last emission so the modelValue watcher can tell "the
// parent echoed back what we just sent" (skip re-sync, keep the caret)
// apart from "the parent swapped to a different post" (sync the DOM).
let lastEmitted = props.modelValue;

function syncFromProp() {
  if (editorRef.value && editorRef.value.innerHTML !== props.modelValue) {
    editorRef.value.innerHTML = props.modelValue;
  }
}

onMounted(() => {
  syncFromProp();
  document.execCommand("defaultParagraphSeparator", false, "p");
});

watch(
  () => props.modelValue,
  (value) => {
    if (value === lastEmitted) return;
    syncFromProp();
  },
);

function emitCurrentContent() {
  if (!editorRef.value) return;
  const html = sanitizeRichTextHtml(editorRef.value.innerHTML);
  lastEmitted = html;
  emit("update:modelValue", html);
}

function focusEditor() {
  editorRef.value?.focus();
}

function currentBlockTag(): string {
  const selection = window.getSelection();
  let node = selection?.anchorNode ?? null;
  while (node && node !== editorRef.value) {
    if (node instanceof HTMLElement && ["P", "H3", "BLOCKQUOTE"].includes(node.tagName)) {
      return node.tagName;
    }
    node = node.parentNode;
  }
  return "P";
}

function onBold() {
  focusEditor();
  document.execCommand("bold");
  emitCurrentContent();
}

function onItalic() {
  focusEditor();
  document.execCommand("italic");
  emitCurrentContent();
}

function onHeading() {
  focusEditor();
  document.execCommand("formatBlock", false, currentBlockTag() === "H3" ? "P" : "H3");
  emitCurrentContent();
}

function onQuote() {
  focusEditor();
  document.execCommand(
    "formatBlock",
    false,
    currentBlockTag() === "BLOCKQUOTE" ? "P" : "BLOCKQUOTE",
  );
  emitCurrentContent();
}

function onLink() {
  const url = window.prompt("Ссылка (URL):");
  if (!url) return;
  focusEditor();
  const selection = window.getSelection();
  if (selection && !selection.isCollapsed) {
    document.execCommand("createLink", false, url);
  } else {
    document.execCommand("insertHTML", false, `<a href="${url}">${url}</a>`);
  }
  emitCurrentContent();
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
    focusEditor();
    document.execCommand(
      "insertHTML",
      false,
      `<figure><img src="${resolveMediaUrl(relativeUrl)}" alt=""></figure><p><br></p>`,
    );
    emitCurrentContent();
  } catch {
    // upload() already captured the failure in `error` below.
  }
}

function onPaste(event: ClipboardEvent) {
  event.preventDefault();
  const html = event.clipboardData?.getData("text/html");
  const text = event.clipboardData?.getData("text/plain") ?? "";
  const cleaned = html
    ? sanitizeRichTextHtml(html)
    : text.replace(/&/g, "&amp;").replace(/</g, "&lt;");
  document.execCommand("insertHTML", false, cleaned || text);
  emitCurrentContent();
}

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
        title="Заголовок"
        @click="onHeading"
      >
        <Heading2 class="size-5" />
      </button>
      <button type="button" class="rounded p-1.5 hover:bg-gray-200" title="Цитата" @click="onQuote">
        <Quote class="size-5" />
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
    </div>
    <div
      ref="editorRef"
      contenteditable="true"
      :data-placeholder="placeholder"
      class="rich-text-editor min-h-32 w-full resize-y rounded-b px-3 py-2 outline-none"
      @input="emitCurrentContent"
      @keydown="onKeydown"
      @paste="onPaste"
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

<style scoped>
.rich-text-editor:empty::before {
  content: attr(data-placeholder);
  color: #9ca3af;
}
.rich-text-editor :deep(p) {
  margin-bottom: 0.75em;
}
.rich-text-editor :deep(h3) {
  margin-top: 1em;
  margin-bottom: 0.5em;
  font-weight: 700;
  font-size: 1.25em;
}
.rich-text-editor :deep(blockquote) {
  margin: 0.75em 0;
  padding-left: 0.75em;
  border-left: 3px solid #d1d5db;
  color: #4b5563;
}
.rich-text-editor :deep(figure) {
  margin: 0.75em 0;
}
.rich-text-editor :deep(img) {
  max-width: 100%;
  border-radius: 0.25rem;
}
.rich-text-editor :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
}
</style>
