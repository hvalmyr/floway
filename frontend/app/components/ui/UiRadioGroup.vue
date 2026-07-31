<script setup lang="ts">
import { useField } from "vee-validate";
import { useId } from "vue";

/**
 * Custom radio group wired to vee-validate via useField(name). The whole row
 * (not just the dot) is clickable and at least 44px tall for touch targets.
 *
 * @example
 * <UiRadioGroup
 *   name="contactMethod"
 *   label="Как с вами связаться?"
 *   :options="[
 *     { value: 'call', label: 'Позвоните мне' },
 *     { value: 'telegram', label: 'Напишите мне в Telegram' },
 *   ]"
 * />
 */
const props = defineProps<{
  name: string;
  label: string;
  options: { value: string; label: string }[];
}>();

const { value, errorMessage } = useField<string>(() => props.name);
const groupName = useId();
const errorId = useId();
</script>

<template>
  <fieldset
    class="flex flex-col gap-4"
    :aria-describedby="errorMessage ? errorId : undefined"
    :aria-invalid="!!errorMessage"
  >
    <legend class="mb-4 text-small text-ink-700">{{ label }}</legend>
    <label
      v-for="option in options"
      :key="option.value"
      class="flex min-h-[44px] cursor-pointer items-center gap-12 rounded-sm px-12 py-8 text-body text-ink-900 hover:bg-primary-50"
    >
      <span class="relative grid size-[20px] shrink-0 place-items-center rounded-full border-[1.5px] border-primary-300">
        <input
          v-model="value"
          type="radio"
          :name="groupName"
          :value="option.value"
          class="peer absolute inset-0 size-full cursor-pointer appearance-none opacity-0"
        />
        <span
          class="size-[10px] rounded-full bg-primary-600 transition-transform duration-150 motion-reduce:transition-none"
          :class="value === option.value ? 'scale-100' : 'scale-0'"
          aria-hidden="true"
        />
        <span
          class="pointer-events-none absolute -inset-2 rounded-full peer-focus-visible:outline peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-ink-900"
        />
      </span>
      {{ option.label }}
    </label>
    <p v-if="errorMessage" :id="errorId" class="text-small text-error">{{ errorMessage }}</p>
  </fieldset>
</template>
