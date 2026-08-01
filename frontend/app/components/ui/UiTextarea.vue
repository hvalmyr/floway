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
    <label :for="inputId" class="font-display text-small font-bold text-primary">
      {{ label }}<span v-if="required" aria-hidden="true"> *</span>
    </label>
    <textarea
      :id="inputId"
      v-model="value"
      :rows="rows"
      :placeholder="placeholder"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorMessage ? errorId : undefined"
      class="resize-y rounded-lg border-2 border-primary bg-white px-24 py-16 font-body text-body text-ink outline-none placeholder:text-primary/50"
      @blur="handleBlur"
    />
    <p v-if="errorMessage" :id="errorId" class="text-small font-bold text-ink">
      {{ errorMessage }}
    </p>
  </div>
</template>
