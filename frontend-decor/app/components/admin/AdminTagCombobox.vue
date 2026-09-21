<script setup lang="ts">
import { onClickOutside } from "@vueuse/core";
import { readableTextColor } from "~/lib/tagColor";
import type { Tag, TagType } from "~/types/api";

/**
 * Inline autocomplete for one of the two independent tag types (product /
 * client-type) on a client — type to search existing tags or create a new
 * one on the spot. No separate tag-management screen: this is the entire
 * tagging UI, per the "create tags on the fly" decision.
 *
 * @example
 * <AdminTagCombobox :client-id="detail.id" tag-type="product" v-model="detail.productTags" />
 */
const props = defineProps<{
  clientId: number;
  tagType: TagType;
}>();

const modelValue = defineModel<Tag[]>({ required: true });

const api = useApiClient();
const endpointPath = computed(() =>
  props.tagType === "product"
    ? `/api/v1/clients/${props.clientId}/tags/product`
    : `/api/v1/clients/${props.clientId}/tags/client-type`,
);

const inputValue = ref("");
const suggestions = ref<Tag[]>([]);
const open = ref(false);
const rootRef = ref<HTMLElement | null>(null);
const saving = ref(false);

onClickOutside(rootRef, () => {
  open.value = false;
});

let searchTimer: ReturnType<typeof setTimeout> | undefined;
watch(inputValue, (query) => {
  clearTimeout(searchTimer);
  if (!query.trim()) {
    suggestions.value = [];
    return;
  }
  searchTimer = setTimeout(async () => {
    suggestions.value = await api<Tag[]>("/api/v1/tags", {
      query: { type: props.tagType, q: query },
    });
  }, 200);
});

async function applyNames(names: string[]) {
  saving.value = true;
  try {
    modelValue.value = await api<Tag[]>(endpointPath.value, {
      method: "PUT",
      body: { tagNames: names },
    });
  } finally {
    saving.value = false;
  }
}

async function addTag(name: string) {
  const trimmed = name.trim();
  if (!trimmed) return;
  if (modelValue.value.some((t) => t.name.toLowerCase() === trimmed.toLowerCase())) {
    inputValue.value = "";
    open.value = false;
    return;
  }
  await applyNames([...modelValue.value.map((t) => t.name), trimmed]);
  inputValue.value = "";
  suggestions.value = [];
  open.value = false;
}

async function removeTag(id: number) {
  await applyNames(modelValue.value.filter((t) => t.id !== id).map((t) => t.name));
}

// Deletes the tag definition itself (not just this client's assignment) —
// it disappears from every client that had it and from the autocomplete
// entirely. The backend cascades the unassignment server-side, so this
// just mirrors that locally rather than issuing a second PUT.
async function deleteTagDefinition(tag: Tag) {
  if (!confirm(`Удалить тег «${tag.name}» полностью — он исчезнет у всех клиентов?`)) return;
  await api(`/api/v1/tags/${tag.id}`, { method: "DELETE", query: { type: props.tagType } });
  suggestions.value = suggestions.value.filter((t) => t.id !== tag.id);
  modelValue.value = modelValue.value.filter((t) => t.id !== tag.id);
}

// A single hidden <input type="color"> is reused for whichever tag's
// swatch was clicked (colorEditingId), rather than one input per tag —
// the color lives on the tag definition, so changing it here updates it
// everywhere the tag appears, not just in this client's chip list.
const colorInputRef = ref<HTMLInputElement | null>(null);
const colorEditingId = ref<number | null>(null);

function openColorPicker(tag: Tag) {
  colorEditingId.value = tag.id;
  const input = colorInputRef.value;
  if (!input) return;
  input.value = tag.color;
  input.click();
}

async function onColorPicked(event: Event) {
  const id = colorEditingId.value;
  if (id === null) return;
  const color = (event.target as HTMLInputElement).value;
  const updated = await api<Tag>(`/api/v1/tags/${id}`, {
    method: "PATCH",
    query: { type: props.tagType },
    body: { color },
  });
  modelValue.value = modelValue.value.map((t) => (t.id === id ? updated : t));
  suggestions.value = suggestions.value.map((t) => (t.id === id ? updated : t));
}
</script>

<template>
  <div ref="rootRef" class="relative">
    <div class="flex flex-wrap items-center gap-2 rounded border border-gray-300 p-2">
      <span
        v-for="tag in modelValue"
        :key="tag.id"
        class="flex items-center gap-1 rounded-full px-2 py-0.5 text-xs"
        :style="{ backgroundColor: tag.color, color: readableTextColor(tag.color) }"
      >
        <button
          type="button"
          title="Изменить цвет тега"
          class="size-3 shrink-0 rounded-full border border-black/10"
          :style="{ backgroundColor: tag.color }"
          @click="openColorPicker(tag)"
        />
        {{ tag.name }}
        <button type="button" class="hover:text-red-600" @click="removeTag(tag.id)">×</button>
      </span>
      <input ref="colorInputRef" type="color" class="hidden" @input="onColorPicked" />
      <input
        v-model="inputValue"
        type="text"
        :disabled="saving"
        placeholder="Добавить тег…"
        class="min-w-32 flex-1 border-none p-1 text-sm outline-none disabled:opacity-50"
        @focus="open = true"
        @keydown.enter.prevent="addTag(inputValue)"
      />
    </div>

    <div
      v-if="open && (suggestions.length > 0 || inputValue.trim())"
      class="absolute inset-x-0 top-full z-10 mt-1 rounded border border-gray-200 bg-white text-sm shadow-md"
    >
      <div
        v-for="tag in suggestions"
        :key="tag.id"
        class="flex items-center justify-between gap-2 hover:bg-gray-50"
      >
        <button
          type="button"
          class="flex flex-1 items-center gap-2 px-3 py-2 text-left"
          @click="addTag(tag.name)"
        >
          <span
            class="size-3 shrink-0 rounded-full border border-black/10"
            :style="{ backgroundColor: tag.color }"
          />
          {{ tag.name }}
        </button>
        <button
          type="button"
          title="Удалить тег полностью"
          class="px-3 py-2 text-gray-400 hover:text-red-600"
          @click="deleteTagDefinition(tag)"
        >
          ×
        </button>
      </div>
      <button
        v-if="
          inputValue.trim() &&
          !suggestions.some((s) => s.name.toLowerCase() === inputValue.trim().toLowerCase())
        "
        type="button"
        class="block w-full px-3 py-2 text-left text-[var(--color-primary)] hover:bg-gray-50"
        @click="addTag(inputValue)"
      >
        Создать «{{ inputValue.trim() }}»
      </button>
    </div>
  </div>
</template>
