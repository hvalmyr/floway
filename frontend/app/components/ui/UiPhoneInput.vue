<script setup lang="ts">
import { useField } from "vee-validate";
import { computed, useId } from "vue";
import { detectPhoneCountry, formatPhoneMask, phoneTemplateFor } from "~/composables/usePhoneMask";

/**
 * Phone input with a live international mask, wired to vee-validate via
 * useField(name). Reuses formatPhoneMask() from usePhoneMask.ts so the
 * masking logic has a single implementation. Defaults to Russian formatting
 * for bare digit entry, but typing a leading "+" and any country's calling
 * code reformats to that country's own pattern (students come from Belarus,
 * Europe, South America too, not just Russia).
 *
 * The rest of the detected country's number shape shows as a muted
 * continuation of what's typed — e.g. typing "7" shows "+7" in normal ink
 * and " 000 000 00 00" ghosted in behind it, live-updating per keystroke
 * (switches to Belgium's shape mid-typing for "+32...", and so on). Built
 * from two overlapping layers rather than the native `placeholder`
 * attribute, which disappears the moment there's any value — an invisible
 * span mirrors the real value char-for-char (so it takes up the exact same
 * width, keeping the ghosted remainder pixel-aligned regardless of font
 * metrics), followed by the muted template tail.
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

const ghostRemainder = computed(() => {
  const current = value.value ?? "";
  const template = phoneTemplateFor(detectPhoneCountry(current));
  return template.slice(current.length);
});

function onInput(event: Event) {
  const target = event.target as HTMLInputElement;
  setValue(formatPhoneMask(target.value));
}

// formatPhoneMask already strips non-digit characters out of the stored
// value on every input, but that's after-the-fact cleanup — a letter would
// still flash into the field for a frame before being stripped back out.
// Blocking it here means it's never actually insertable by typing at all.
// Scoped to plain typing (insertType "insertText") specifically — paste
// ("insertFromPaste") is left alone and goes through the same cleanup, so
// pasting "+7 (912) 345-67-89 " still works instead of being rejected
// outright just because it contains non-digit formatting characters.
function onBeforeInput(event: InputEvent) {
  if (event.inputType !== "insertText" || !event.data) return;
  if (/[^\d+]/.test(event.data)) event.preventDefault();
}
</script>

<template>
  <div class="flex flex-col gap-8">
    <label :for="inputId" class="font-body text-body font-light text-ink">
      {{ label }}<span v-if="required" aria-hidden="true"> *</span>
    </label>
    <div class="relative h-[56px] overflow-hidden rounded-[15px] bg-surface">
      <div
        class="pointer-events-none absolute inset-0 flex items-center whitespace-pre px-[20px] font-body text-body"
        aria-hidden="true"
      >
        <span class="invisible">{{ value }}</span
        ><span class="text-ink/40">{{ ghostRemainder }}</span>
      </div>
      <input
        :id="inputId"
        :value="value"
        type="tel"
        inputmode="tel"
        :aria-invalid="!!errorMessage"
        :aria-describedby="errorId"
        class="absolute inset-0 h-full w-full rounded-[15px] border-2 bg-transparent px-[20px] font-body text-body text-ink outline-none"
        :class="errorMessage ? 'border-primary' : 'border-transparent'"
        @input="onInput"
        @beforeinput="onBeforeInput"
        @blur="handleBlur"
      />
    </div>
    <!-- Рендерится только при ошибке — поле "растёт" вниз вместо того, чтобы
    всегда резервировать место под текст ошибки. -->
    <p v-if="errorMessage" :id="errorId" class="font-body text-body text-primary">
      {{ errorMessage }}
    </p>
  </div>
</template>
