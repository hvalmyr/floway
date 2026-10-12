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
 * An option can carry `chips` — a second-level choice (e.g. "which social
 * network") that expands inline under that option inside the still-open
 * popup, instead of a second field appearing below the select once it
 * closes. Picking a chip sets this select's own `value` to the parent
 * option *and* writes the chip's value to a second vee-validate field named
 * by `chipFieldName`. Options without `chips` behave exactly as before
 * (pick → close).
 *
 * @example
 * <UiSelect
 *   name="source"
 *   label="Как вы о нас узнали?"
 *   :options="[
 *     { value: 'ads', label: 'Реклама' },
 *     { value: 'social', label: 'В социальных сетях', chips: [{ value: 'vk', label: 'ВКонтакте' }] },
 *   ]"
 *   chip-field-name="socialNetwork"
 * />
 */
interface SelectOption {
  value: string;
  label: string;
  chips?: { value: string; label: string }[];
}

const props = withDefaults(
  defineProps<{
    name: string;
    label: string;
    options: SelectOption[];
    placeholder?: string;
    required?: boolean;
    /** vee-validate field name the chip row writes to — required by any
     * option that sets `chips`. */
    chipFieldName?: string;
  }>(),
  { required: false, chipFieldName: undefined },
);

const { value, errorMessage } = useField<string | undefined>(() => props.name);
// Conditional composable call is fine here — Vue has no React-style "rules
// of hooks": chipFieldName is a per-instance prop, not something that
// changes across this component's own lifetime.
const chipField = props.chipFieldName
  ? useField<string | undefined>(() => props.chipFieldName as string)
  : null;
const inputId = useId();
const errorId = useId();
const open = ref(false);
// Which option's chip row is currently expanded inside the open popup —
// null when no chip-bearing option is being picked.
const expandedChipsFor = ref<string | null>(null);
const rootRef = ref<HTMLElement | null>(null);

onClickOutside(rootRef, () => {
  open.value = false;
});

function toggleOpen() {
  if (!open.value) {
    // Re-expand the current value's chip row on reopen, so editing an
    // already-made "social" pick doesn't require reselecting it first.
    const current = props.options.find((o) => o.value === value.value);
    expandedChipsFor.value = current?.chips ? current.value : null;
  }
  open.value = !open.value;
}

function select(option: SelectOption) {
  value.value = option.value;
  if (option.chips && option.chips.length > 0) {
    // Tentatively record the category; keep the popup open until a chip
    // is actually picked.
    expandedChipsFor.value = option.value;
    return;
  }
  expandedChipsFor.value = null;
  open.value = false;
}

function selectChip(option: SelectOption, chip: { value: string; label: string }) {
  value.value = option.value;
  if (chipField) chipField.value.value = chip.value;
  expandedChipsFor.value = null;
  open.value = false;
}

const selectedLabel = () => {
  const option = props.options.find((o) => o.value === value.value);
  if (!option) return undefined;
  const chip = option.chips?.find((c) => c.value === chipField?.value.value);
  return chip ? `${option.label} — ${chip.label}` : option.label;
};

const combinedErrorMessage = computed(() => errorMessage.value || chipField?.errorMessage.value);
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
        combinedErrorMessage ? 'border-primary' : open ? 'border-primary' : 'border-transparent',
        selectedLabel() ? 'text-ink' : 'text-ink/45',
      ]"
      aria-haspopup="listbox"
      :aria-expanded="open"
      :aria-describedby="errorId"
      @click="toggleOpen"
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
      <div v-for="option in options" :key="option.value">
        <button
          type="button"
          role="option"
          :aria-selected="value === option.value"
          class="flex w-full items-center justify-between gap-12 rounded-[10px] px-12 py-8 text-left font-body text-body text-ink"
          :class="value === option.value ? 'bg-surface' : 'hover:bg-surface/60'"
          @click="select(option)"
        >
          {{ option.label }}
          <Check
            v-if="value === option.value && !option.chips"
            class="size-[16px] shrink-0 text-primary"
            aria-hidden="true"
          />
        </button>
        <!-- The option's own second-level choice (e.g. "which social
        network") — expands right here, inside the still-open popup,
        instead of a second field appearing below the select once closed. -->
        <div
          v-if="option.chips && expandedChipsFor === option.value"
          class="flex flex-wrap gap-4 px-12 pb-8 pt-2"
        >
          <button
            v-for="chip in option.chips"
            :key="chip.value"
            type="button"
            class="flex items-center gap-8 rounded-[10px] px-12 py-8 text-left font-body text-body text-ink"
            :class="chipField?.value.value === chip.value ? 'bg-surface' : 'hover:bg-surface/60'"
            @click="selectChip(option, chip)"
          >
            {{ chip.label }}
            <Check
              v-if="chipField?.value.value === chip.value"
              class="size-[16px] shrink-0 text-primary"
              aria-hidden="true"
            />
          </button>
        </div>
      </div>
    </div>

    <!-- Рендерится только при ошибке — поле "растёт" вниз вместо того, чтобы
    всегда резервировать место под текст ошибки. -->
    <p v-if="combinedErrorMessage" :id="errorId" class="font-body text-body text-primary">
      {{ combinedErrorMessage }}
    </p>
  </div>
</template>
