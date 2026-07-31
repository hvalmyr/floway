<script setup lang="ts">
import { useField } from "vee-validate";
import { useId } from "vue";

/**
 * Text/email input wired to vee-validate via useField(name). Must be used
 * inside a component tree where useForm() ran (e.g. ApplyForm.vue) — the
 * field automatically attaches to that form's context.
 *
 * @example
 * <UiInput name="name" label="Ваше имя" required autocomplete="name" />
 */
const props = withDefaults(
  defineProps<{
    name: string;
    label: string;
    type?: string;
    placeholder?: string;
    required?: boolean;
    autocomplete?: string;
  }>(),
  { type: "text", required: false },
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
    <input
      :id="inputId"
      v-model="value"
      :type="type"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorMessage ? errorId : undefined"
      class="h-[56px] rounded-sm border-[1.5px] border-line bg-white px-16 text-body text-ink-900 outline-none placeholder:text-ink-400 focus:border-primary-600 focus:shadow-[0_0_0_3px_var(--color-primary-50)]"
      @blur="handleBlur"
    />
    <p v-if="errorMessage" :id="errorId" class="text-small text-error">{{ errorMessage }}</p>
  </div>
</template>
