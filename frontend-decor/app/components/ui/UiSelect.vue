<script setup lang="ts">
import { onClickOutside } from "@vueuse/core";
import { Check, ChevronDown } from "lucide-vue-next";
import { useField } from "vee-validate";
import { ref, useId } from "vue";

/**
 * Custom single-select dropdown wired to vee-validate via useField(name).
 * Replaces a native <select> so it can match the cream-field visual style
 * shared with UiInput/UiPhoneInput (per the redesigned form mockup) — a
 * native select can't be restyled to that degree across browsers.
 *
 * @example
 * <UiSelect
 *   name="source"
 *   label="Как вы о нас узнали?"
 *   :options="[{ value: 'ads', label: 'Реклама' }, ...]"
 * />
 */
const props = withDefaults(
  defineProps<{
    name: string;
    label: string;
    options: { value: string; label: string }[];
    placeholder?: string;
    required?: boolean;
  }>(),
  { required: false },
);

const { value, errorMessage } = useField<string | undefined>(() => props.name);
const inputId = useId();
const errorId = useId();
const open = ref(false);
const rootRef = ref<HTMLElement | null>(null);

onClickOutside(rootRef, () => {
  open.value = false;
});

function select(optionValue: string) {
  value.value = optionValue;
  open.value = false;
}

const selectedLabel = () => props.options.find((o) => o.value === value.value)?.label;
</script>

<template>
  <div ref="rootRef" class="relative flex flex-col gap-8">
    <label :for="inputId" class="font-body text-body font-light text-ink">
      {{ label }}<span v-if="required" aria-hidden="true"> *</span>
    </label>
    <button
      :id="inputId"
      type="button"
      class="flex h-[56px] w-full items-center justify-between gap-12 rounded-[15px] border-2 bg-surface px-[20px] text-left font-body text-body outline-none"
      :class="[
        errorMessage ? 'border-primary' : open ? 'border-primary' : 'border-transparent',
        selectedLabel() ? 'text-ink' : 'text-ink/45',
      ]"
      aria-haspopup="listbox"
      :aria-expanded="open"
      :aria-describedby="errorId"
      @click="open = !open"
      @keydown.esc="open = false"
    >
      <span>{{ selectedLabel() ?? placeholder ?? "Выберите вариант" }}</span>
      <ChevronDown
        class="size-[14px] shrink-0 text-ink transition-transform duration-150"
        :class="open ? 'rotate-180' : ''"
        aria-hidden="true"
      />
    </button>

    <div
      v-if="open"
      role="listbox"
      class="absolute inset-x-0 top-full z-20 mt-8 flex flex-col gap-4 rounded-[15px] bg-white p-8 shadow-[0_12px_32px_rgba(65,52,42,0.16)]"
    >
      <button
        v-for="option in options"
        :key="option.value"
        type="button"
        role="option"
        :aria-selected="value === option.value"
        class="flex items-center justify-between gap-12 rounded-[10px] px-12 py-8 text-left font-body text-body text-ink"
        :class="value === option.value ? 'bg-surface' : 'hover:bg-surface/60'"
        @click="select(option.value)"
      >
        {{ option.label }}
        <Check
          v-if="value === option.value"
          class="size-[16px] shrink-0 text-primary"
          aria-hidden="true"
        />
      </button>
    </div>

    <!-- Рендерится только при ошибке — поле "растёт" вниз вместо того, чтобы
    всегда резервировать место под текст ошибки. -->
    <p v-if="errorMessage" :id="errorId" class="font-body text-body text-primary">
      {{ errorMessage }}
    </p>
  </div>
</template>
