<script setup lang="ts">
import { useField } from "vee-validate";
import { useId } from "vue";

/**
 * Single checkbox wired to vee-validate via useField(name), boolean-valued.
 * Label content comes through the default slot (not a `label` prop) since
 * the one real use case — data-processing consent — needs an inline link
 * inside the label text.
 *
 * @example
 * <UiCheckbox name="consent">
 *   Отправляя форму, вы соглашаетесь с
 *   <NuxtLink to="/privacy" class="text-primary underline">политикой конфиденциальности</NuxtLink>.
 * </UiCheckbox>
 */
const props = defineProps<{ name: string }>();

const { value, errorMessage, handleBlur } = useField<boolean>(() => props.name, undefined, {
  initialValue: false,
});
const inputId = useId();
const errorId = useId();
</script>

<template>
  <div class="flex flex-col gap-8">
    <label :for="inputId" class="flex cursor-pointer items-start gap-12">
      <span
        class="relative mt-2 grid size-[22px] shrink-0 place-items-center rounded-[6px] border-2 transition-colors"
        :class="[
          value ? 'border-primary bg-primary' : 'bg-white',
          !value && errorMessage ? 'border-primary' : '',
          !value && !errorMessage ? 'border-ink' : '',
        ]"
      >
        <input
          :id="inputId"
          v-model="value"
          type="checkbox"
          :aria-invalid="!!errorMessage"
          :aria-describedby="errorId"
          class="peer absolute inset-0 size-full cursor-pointer appearance-none"
          @blur="handleBlur"
        />
        <svg
          v-if="value"
          class="pointer-events-none size-12"
          viewBox="0 0 28 28"
          fill="none"
          aria-hidden="true"
        >
          <path
            d="M 12 16 L 0 16 L 0 12 L 12 12 L 12 0 L 16 0 L 16 12 L 28 12 L 28 16 L 16 16 L 16 28 L 12 28 L 12 16 Z"
            fill="white"
            transform="rotate(45 14 14)"
          />
        </svg>
        <span
          class="pointer-events-none absolute -inset-2 rounded-sm peer-focus-visible:outline peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-ink"
        />
      </span>
      <span class="font-body text-body text-ink"><slot /></span>
    </label>
    <!-- Рендерится только при ошибке — поле "растёт" вниз вместо того, чтобы
    всегда резервировать место под текст ошибки. -->
    <p v-if="errorMessage" :id="errorId" class="font-body text-body text-primary">
      {{ errorMessage }}
    </p>
  </div>
</template>
