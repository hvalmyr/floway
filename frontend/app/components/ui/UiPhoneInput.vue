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
    <label :for="inputId" class="font-display text-small font-bold text-primary">
      {{ label }}<span v-if="required" aria-hidden="true"> *</span>
    </label>
    <input
      :id="inputId"
      :value="value"
      type="tel"
      inputmode="tel"
      placeholder="+7 (000) 000 00 00"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorMessage ? errorId : undefined"
      class="h-[64px] rounded-lg border-2 border-primary bg-white px-24 font-body text-body text-ink outline-none placeholder:text-primary/50"
      @input="onInput"
      @blur="handleBlur"
    />
    <p v-if="errorMessage" :id="errorId" class="text-small font-bold text-ink">
      {{ errorMessage }}
    </p>
  </div>
</template>
