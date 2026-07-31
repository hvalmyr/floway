<script setup lang="ts">
import { useField } from "vee-validate";
import { useId } from "vue";
import { formatPhoneMask } from "~/composables/usePhoneMask";

/**
 * Phone input with the +7 (000) 000 00 00 mask, wired to vee-validate via
 * useField(name). Reuses formatPhoneMask() from usePhoneMask.ts so the
 * masking logic has a single implementation.
 *
 * @example
 * <UiPhoneInput name="phone" label="Номер телефона" required />
 */
const props = withDefaults(
  defineProps<{
    name: string;
    label: string;
    required?: boolean;
  }>(),
  { required: false },
);

const { value, errorMessage, handleBlur, setValue } = useField<string>(() => props.name);
const inputId = useId();
const errorId = useId();

function onInput(event: Event) {
  const target = event.target as HTMLInputElement;
  setValue(formatPhoneMask(target.value));
}
</script>

<template>
  <div class="flex flex-col gap-8">
    <label :for="inputId" class="text-small text-ink-700">
      {{ label }}<span v-if="required" aria-hidden="true" class="text-error"> *</span>
    </label>
    <input
      :id="inputId"
      :value="value"
      type="tel"
      inputmode="tel"
      placeholder="+7 (___) ___ __ __"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorMessage ? errorId : undefined"
      class="h-[56px] rounded-sm border-[1.5px] border-line bg-white px-16 text-body text-ink-900 outline-none placeholder:text-ink-400 focus:border-primary-600 focus:shadow-[0_0_0_3px_var(--color-primary-50)]"
      @input="onInput"
      @blur="handleBlur"
    />
    <p v-if="errorMessage" :id="errorId" class="text-small text-error">{{ errorMessage }}</p>
  </div>
</template>
