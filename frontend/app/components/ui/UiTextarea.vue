<script setup lang="ts">
import { useField } from "vee-validate";
import { useId } from "vue";

/**
 * Multiline text field wired to vee-validate via useField(name). See
 * UiInput.vue for the shared usage pattern.
 *
 * @example
 * <UiTextarea name="comment" label="Комментарий" :rows="4" />
 */
const props = withDefaults(
  defineProps<{
    name: string;
    label: string;
    placeholder?: string;
    required?: boolean;
    rows?: number;
  }>(),
  { required: false, rows: 4 },
);

const { value, errorMessage, handleBlur } = useField<string>(() => props.name);
const inputId = useId();
const errorId = useId();
</script>

<template>
  <div class="flex flex-col gap-8">
    <label :for="inputId" class="text-small text-ink-700">
      {{ label }}<span v-if="required" aria-hidden="true" class="text-error"> *</span>
    </label>
    <textarea
      :id="inputId"
      v-model="value"
      :rows="rows"
      :placeholder="placeholder"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorMessage ? errorId : undefined"
      class="resize-y rounded-sm border-[1.5px] border-line bg-white px-16 py-12 text-body text-ink-900 outline-none placeholder:text-ink-400 focus:border-primary-600 focus:shadow-[0_0_0_3px_var(--color-primary-50)]"
      @blur="handleBlur"
    />
    <p v-if="errorMessage" :id="errorId" class="text-small text-error">{{ errorMessage }}</p>
  </div>
</template>
