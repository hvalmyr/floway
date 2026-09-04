<script setup lang="ts">
import { onClickOutside } from "@vueuse/core";
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
</script>

<template>
  <div ref="rootRef" class="relative">
    <div class="flex flex-wrap items-center gap-2 rounded border border-gray-300 p-2">
      <span
        v-for="tag in modelValue"
        :key="tag.id"
        class="flex items-center gap-1 rounded-full bg-gray-100 px-2 py-0.5 text-xs text-[var(--color-text-muted)]"
      >
        {{ tag.name }}
        <button type="button" class="text-gray-400 hover:text-red-600" @click="removeTag(tag.id)">
          ×
        </button>
      </span>
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
      <button
        v-for="tag in suggestions"
        :key="tag.id"
        type="button"
        class="block w-full px-3 py-2 text-left hover:bg-gray-50"
        @click="addTag(tag.name)"
      >
        {{ tag.name }}
      </button>
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
