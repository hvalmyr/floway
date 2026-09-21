<script setup lang="ts">
import {
  Bold,
  Heading2,
  Image as ImageIcon,
  Italic,
  Link as LinkIcon,
  Quote,
  X,
} from "lucide-vue-next";
import { onClickOutside } from "@vueuse/core";
import { nextTick, onMounted, onUnmounted, ref } from "vue";
import { normalizeLinkUrl, sanitizeRichTextHtml } from "~/lib/richTextSanitize";

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

const wrapperRef = ref<HTMLDivElement | null>(null);
const editorRef = ref<HTMLDivElement | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);
const { upload, uploading, error } = useAdminUpload();

// Clicking an inserted image selects its <figure> and shows a floating
// delete button over it — contenteditable's native "select + Backspace"
// doesn't reliably remove a replaced element like an <img>, so this is the
// only way to remove one once inserted.
const selectedFigure = ref<HTMLElement | null>(null);
const deleteButtonStyle = ref<{ top: string; left: string } | null>(null);

function updateDeleteButtonPosition() {
  if (!selectedFigure.value || !wrapperRef.value) {
    deleteButtonStyle.value = null;
    return;
  }
  const figureRect = selectedFigure.value.getBoundingClientRect();
  const wrapperRect = wrapperRef.value.getBoundingClientRect();
  deleteButtonStyle.value = {
    top: `${figureRect.top - wrapperRect.top + 8}px`,
    left: `${figureRect.right - wrapperRect.left - 32}px`,
  };
}

function selectImage(img: HTMLElement) {
  selectedFigure.value = img.closest("figure") ?? img;
  nextTick(updateDeleteButtonPosition);
}

function clearImageSelection() {
  selectedFigure.value = null;
  deleteButtonStyle.value = null;
}

function deleteSelectedImage() {
  if (!selectedFigure.value) return;
  selectedFigure.value.remove();
  clearImageSelection();
  emitCurrentContent();
}

onClickOutside(wrapperRef, clearImageSelection);
onMounted(() => {
  window.addEventListener("scroll", updateDeleteButtonPosition, true);
  window.addEventListener("resize", updateDeleteButtonPosition);
});
onUnmounted(() => {
  window.removeEventListener("scroll", updateDeleteButtonPosition, true);
  window.removeEventListener("resize", updateDeleteButtonPosition);
});

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
    if (node instanceof HTMLElement && ["P", "H2", "H3", "BLOCKQUOTE"].includes(node.tagName)) {
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

// Cycles the current block through normal text -> H2 -> H3 -> normal text,
// both from the toolbar button and from Ctrl+Shift+H.
function onHeading() {
  focusEditor();
  const current = currentBlockTag();
  const next = current === "H2" ? "H3" : current === "H3" ? "P" : "H2";
  document.execCommand("formatBlock", false, next);
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

// Link tool: a small popover instead of a bare prompt(), so an admin can
// either type any URL or pick another blog post from a live-filtered list
// (inserted as a site-relative "/blog/{slug}" href) instead of having to
// know/copy that post's URL by hand. The selection at the moment the tool
// opens is saved and restored on confirm, since by then focus has moved
// into the popover's own input and window.getSelection() would otherwise
// see nothing inside the editor.
const linkPickerRef = ref<HTMLElement | null>(null);
const linkPickerOpen = ref(false);
const linkUrlInput = ref("");
const linkPosts = ref<{ slug: string; title: string }[]>([]);
const linkPostsLoading = ref(false);
let linkPostsLoaded = false;
let savedRange: Range | null = null;

// Same text box drives both the URL to insert and the post search — typing
// a real URL just happens to match nothing in the list below, which is
// harmless (list stays empty, nothing to click).
const filteredLinkPosts = computed(() => {
  const query = linkUrlInput.value.trim().toLowerCase();
  if (!query) return linkPosts.value;
  return linkPosts.value.filter((post) => post.title.toLowerCase().includes(query));
});

async function ensureLinkPostsLoaded() {
  if (linkPostsLoaded) return;
  linkPostsLoading.value = true;
  try {
    const posts = await useApi().getBlogPosts();
    linkPosts.value = posts.map((post) => ({ slug: post.slug, title: post.title }));
    linkPostsLoaded = true;
  } catch {
    linkPosts.value = [];
  } finally {
    linkPostsLoading.value = false;
  }
}

function openLinkPicker() {
  const selection = window.getSelection();
  savedRange =
    selection && selection.rangeCount > 0 && editorRef.value?.contains(selection.anchorNode)
      ? selection.getRangeAt(0).cloneRange()
      : null;
  linkUrlInput.value = "";
  linkPickerOpen.value = true;
  ensureLinkPostsLoaded();
}

function closeLinkPicker() {
  linkPickerOpen.value = false;
}

onClickOutside(linkPickerRef, closeLinkPicker);

function pickLinkPost(post: { slug: string; title: string }) {
  linkUrlInput.value = `/blog/${post.slug}`;
}

function confirmLink() {
  const url = normalizeLinkUrl(linkUrlInput.value);
  if (!url) {
    window.alert("Не похоже на ссылку — проверьте адрес.");
    return;
  }
  focusEditor();
  const selection = window.getSelection();
  if (selection && savedRange) {
    selection.removeAllRanges();
    selection.addRange(savedRange);
  }
  if (selection && !selection.isCollapsed) {
    document.execCommand("createLink", false, url);
  } else {
    document.execCommand("insertHTML", false, `<a href="${url}">${url}</a>`);
  }
  emitCurrentContent();
  linkPickerOpen.value = false;
}

// Image insertion: instead of an immediate file picker + prompt(), clicking
// the toolbar button drops a "![Alt Text](URL or Image Path)" template at
// the caret so alt text and source are entered inline, in the flow of
// writing, the way a markdown-savvy admin already expects that syntax to
// work — see insertImageTemplate(). It only ever becomes a real <img> once
// finalizeImageTemplate() runs (on Enter), so a half-filled template that
// never gets finished just sits as inert text rather than a broken image.
function selectElementContents(el: Node) {
  const range = document.createRange();
  range.selectNodeContents(el);
  const selection = window.getSelection();
  selection?.removeAllRanges();
  selection?.addRange(range);
}

function closestInEditor(node: Node | null, selector: string): HTMLElement | null {
  const el = node instanceof HTMLElement ? node : (node?.parentElement ?? null);
  if (!el || !editorRef.value?.contains(el)) return null;
  return el.closest<HTMLElement>(selector);
}

function selectionAncestor(selector: string): HTMLElement | null {
  return closestInEditor(window.getSelection()?.anchorNode ?? null, selector);
}

function insertImageTemplate() {
  focusEditor();
  const selection = window.getSelection();
  if (
    !selection ||
    selection.rangeCount === 0 ||
    !editorRef.value?.contains(selection.anchorNode)
  ) {
    return;
  }
  const range = selection.getRangeAt(0);
  range.deleteContents();

  const altSpan = document.createElement("span");
  altSpan.className = "img-tpl-alt";
  altSpan.textContent = "Alt Text";

  const urlWord = document.createElement("span");
  urlWord.className = "img-tpl-opt";
  urlWord.dataset.action = "url";
  urlWord.textContent = "URL";

  const pathWord = document.createElement("span");
  pathWord.className = "img-tpl-opt";
  pathWord.dataset.action = "path";
  pathWord.textContent = "Image Path";

  const urlSpan = document.createElement("span");
  urlSpan.className = "img-tpl-url";
  urlSpan.append(urlWord, document.createTextNode(" или "), pathWord);

  const wrapper = document.createElement("span");
  wrapper.className = "img-tpl";
  wrapper.append(
    document.createTextNode("!["),
    altSpan,
    document.createTextNode("]("),
    urlSpan,
    document.createTextNode(")"),
  );

  range.insertNode(wrapper);
  selectElementContents(altSpan);
}

// Set right before fileInput.value.click() so onImageChange() knows which
// template's URL field to fill once the upload resolves — the file picker
// is only ever opened from a template's "Image Path" word, never directly.
const pendingUrlTarget = ref<HTMLElement | null>(null);

function pickImage() {
  fileInput.value?.click();
}

async function onImageChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = "";
  const target = pendingUrlTarget.value;
  pendingUrlTarget.value = null;
  if (!file || !target) return;
  try {
    const relativeUrl = await upload(file);
    target.textContent = resolveMediaUrl(relativeUrl);
    selectElementContents(target);
  } catch {
    // upload() already captured the failure in `error` below.
  }
}

function onEditorClick(event: MouseEvent) {
  const target = event.target as HTMLElement;
  if (target.tagName === "IMG" && closestInEditor(target, "figure")) {
    event.preventDefault();
    selectImage(target);
    return;
  }
  clearImageSelection();
  onTemplateOptionClick(event);
}

function onTemplateOptionClick(event: MouseEvent) {
  const opt = closestInEditor(event.target as Node, ".img-tpl-opt");
  if (!opt) return;
  event.preventDefault();
  const urlSpan = opt.closest<HTMLElement>(".img-tpl-url");
  if (!urlSpan) return;
  if (opt.dataset.action === "path") {
    pendingUrlTarget.value = urlSpan;
    pickImage();
    return;
  }
  urlSpan.textContent = "";
  focusEditor();
  selectElementContents(urlSpan);
}

// Ready once the two option words have been replaced by a real value,
// either typed by hand or filled in by onImageChange() after an upload.
function finalizeImageTemplate(tpl: HTMLElement): boolean {
  const urlSpan = tpl.querySelector(".img-tpl-url");
  if (!urlSpan || urlSpan.querySelector(".img-tpl-opt")) return false;
  const url = normalizeLinkUrl(urlSpan.textContent?.trim() ?? "");
  if (!url) {
    window.alert("Не похоже на ссылку — проверьте адрес картинки.");
    return false;
  }
  const altRaw = tpl.querySelector(".img-tpl-alt")?.textContent?.trim() ?? "";
  const alt = altRaw && altRaw !== "Alt Text" ? altRaw : "";

  const fragment = document
    .createRange()
    .createContextualFragment(
      `<figure><img src="${escapeHtmlAttr(url)}" alt="${escapeHtmlAttr(alt)}"></figure><p><br></p>`,
    );
  const paragraph = fragment.querySelector("p");
  tpl.replaceWith(fragment);
  if (paragraph) selectElementContents(paragraph);
  emitCurrentContent();
  return true;
}

function escapeHtmlAttr(value: string): string {
  return value.replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;");
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
  if ((event.key === "Delete" || event.key === "Backspace") && selectedFigure.value) {
    event.preventDefault();
    deleteSelectedImage();
    return;
  }
  if (event.key === "Tab" && selectionAncestor(".img-tpl-alt")) {
    event.preventDefault();
    const urlSpan = selectionAncestor(".img-tpl")?.querySelector<HTMLElement>(".img-tpl-url");
    if (urlSpan) selectElementContents(urlSpan);
    return;
  }
  if (event.key === "Enter") {
    const tpl = selectionAncestor(".img-tpl");
    if (tpl) {
      event.preventDefault();
      finalizeImageTemplate(tpl);
      return;
    }
  }
  if (!(event.ctrlKey || event.metaKey)) return;
  const key = event.key.toLowerCase();
  if (event.shiftKey && key === "i") {
    event.preventDefault();
    insertImageTemplate();
    return;
  }
  if (event.shiftKey && key === "h") {
    event.preventDefault();
    onHeading();
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
      openLinkPicker();
      break;
  }
}
</script>

<template>
  <div ref="wrapperRef" class="relative flex flex-col rounded border border-gray-300">
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
        title="Заголовок (Ctrl+Shift+H)"
        @click="onHeading"
      >
        <Heading2 class="size-5" />
      </button>
      <button type="button" class="rounded p-1.5 hover:bg-gray-200" title="Цитата" @click="onQuote">
        <Quote class="size-5" />
      </button>
      <div ref="linkPickerRef" class="relative">
        <button
          type="button"
          class="rounded p-1.5 hover:bg-gray-200"
          title="Ссылка (Ctrl+K)"
          @click="openLinkPicker"
        >
          <LinkIcon class="size-5" />
        </button>
        <div
          v-if="linkPickerOpen"
          class="absolute left-0 top-full z-10 mt-1 w-80 rounded border border-gray-200 bg-white p-3 text-sm shadow-md"
        >
          <input
            v-model="linkUrlInput"
            type="text"
            placeholder="URL или название статьи блога…"
            class="w-full rounded border border-gray-300 px-2 py-1 text-sm"
            @keydown.enter.prevent="confirmLink"
          />
          <p v-if="linkPostsLoading" class="mt-2 text-xs text-[var(--color-text-muted)]">
            Загрузка статей…
          </p>
          <div v-else-if="filteredLinkPosts.length" class="mt-2 max-h-40 overflow-y-auto">
            <button
              v-for="post in filteredLinkPosts"
              :key="post.slug"
              type="button"
              class="block w-full truncate rounded px-2 py-1 text-left hover:bg-gray-50"
              @click="pickLinkPost(post)"
            >
              {{ post.title }}
            </button>
          </div>
          <div class="mt-2 flex justify-end gap-2">
            <button type="button" class="rounded px-2 py-1 text-xs" @click="closeLinkPicker">
              Отмена
            </button>
            <button
              type="button"
              class="rounded bg-[var(--color-primary)] px-2 py-1 text-xs text-white"
              @click="confirmLink"
            >
              Вставить
            </button>
          </div>
        </div>
      </div>
      <button
        type="button"
        class="rounded p-1.5 hover:bg-gray-200 disabled:opacity-50"
        title="Картинка (Ctrl+Shift+I)"
        :disabled="uploading"
        @click="insertImageTemplate"
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
      @click="onEditorClick"
    />
    <button
      v-if="selectedFigure && deleteButtonStyle"
      type="button"
      class="absolute z-10 grid size-6 place-items-center rounded-full bg-red-600 text-white shadow hover:bg-red-700"
      title="Удалить картинку"
      :style="deleteButtonStyle"
      @click="deleteSelectedImage"
    >
      <X class="size-4" />
    </button>
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
.rich-text-editor :deep(h2),
.rich-text-editor :deep(h3) {
  margin-top: 1em;
  margin-bottom: 0.5em;
  font-family: var(--font-display);
  font-weight: 700;
}
.rich-text-editor :deep(h2) {
  font-size: var(--text-h2);
}
.rich-text-editor :deep(h3) {
  font-size: var(--text-h4);
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
  max-height: 400px;
  border-radius: 0.25rem;
  /* display:block + auto margins, not the figure's text-align — Tailwind's
     preflight forces img to display:block, which text-align can't center. */
  display: block;
  margin-left: auto;
  margin-right: auto;
}
.rich-text-editor :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
}
.rich-text-editor :deep(.img-tpl-opt) {
  color: var(--color-primary);
  text-decoration: underline;
  cursor: pointer;
}
</style>
