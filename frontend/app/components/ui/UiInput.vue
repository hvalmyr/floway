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
    <label :for="inputId" class="font-display text-small font-bold text-primary">
      {{ label }}<span v-if="required" aria-hidden="true"> *</span>
    </label>
    <input
      :id="inputId"
      v-model="value"
      :type="type"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorMessage ? errorId : undefined"
      class="h-[64px] rounded-lg border-2 border-primary bg-white px-24 font-body text-body text-ink outline-none placeholder:text-primary/50"
      @blur="handleBlur"
    />
    <p v-if="errorMessage" :id="errorId" class="text-small font-bold text-ink">{{ errorMessage }}</p>
  </div>
</template>
