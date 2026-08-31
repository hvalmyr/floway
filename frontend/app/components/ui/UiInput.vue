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
    <label :for="inputId" class="font-body text-body font-light text-ink">
      {{ label }}<span v-if="required" aria-hidden="true"> *</span>
    </label>
    <input
      :id="inputId"
      v-model="value"
      :type="type"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorId"
      class="h-[56px] rounded-[15px] border-2 bg-surface px-[20px] font-body text-body text-ink outline-none placeholder:text-ink/45"
      :class="errorMessage ? 'border-primary' : 'border-transparent'"
      @blur="handleBlur"
    />
    <!-- Рендерится только при ошибке — поле "растёт" вниз вместо того, чтобы
    всегда резервировать место под текст ошибки. -->
    <p v-if="errorMessage" :id="errorId" class="font-body text-body text-primary">
      {{ errorMessage }}
    </p>
  </div>
</template>
