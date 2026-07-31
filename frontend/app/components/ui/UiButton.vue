<script setup lang="ts">
/**
 * Pill-shaped button/link primitive. Renders a <NuxtLink> when `to` is set,
 * otherwise a <button>. Variants follow docs/floway-design.md §5.1.
 *
 * @example
 * <UiButton variant="primary" size="lg" to="/masterclasses">Мастер-классы</UiButton>
 * <UiButton variant="outline" size="md">О школе</UiButton>
 * <UiButton type="submit" block :loading="pending" :disabled="pending">Отправить заявку</UiButton>
 * <UiButton variant="outline-inverted" size="md" to="/courses/osnovy-floristiki">Узнать больше</UiButton>
 */
const props = withDefaults(
  defineProps<{
    variant?: "primary" | "outline" | "outline-inverted";
    size?: "md" | "lg";
    to?: string;
    type?: "button" | "submit" | "reset";
    disabled?: boolean;
    loading?: boolean;
    /** Always full-width, regardless of breakpoint (used for the apply-form submit button). */
    block?: boolean;
  }>(),
  {
    variant: "primary",
    size: "md",
    to: undefined,
    type: "button",
    disabled: false,
    loading: false,
    block: false,
  },
);

function onClick(event: MouseEvent) {
  if (props.to && (props.disabled || props.loading)) {
    event.preventDefault();
  }
}
</script>

<template>
  <component
    :is="to ? 'NuxtLink' : 'button'"
    :to="to"
    :type="to ? undefined : type"
    :disabled="!to && (disabled || loading)"
    :aria-busy="loading || undefined"
    :aria-disabled="to && (disabled || loading) ? 'true' : undefined"
    class="inline-flex items-center justify-center gap-8 rounded-pill font-display text-button font-bold transition-[background-color,transform,border-color] duration-200"
    :class="[
      block ? 'w-full' : '',
      size === 'lg' ? 'h-[44px] px-24 sm:h-[64px] sm:px-40' : 'h-[44px] px-24 sm:h-[56px] sm:px-32',
      variant === 'primary' &&
        'bg-primary-600 text-white hover:-translate-y-px hover:bg-primary-700 active:translate-y-0 disabled:translate-y-0 disabled:cursor-not-allowed disabled:bg-surface disabled:text-ink-400',
      variant === 'outline' &&
        'border-2 border-ink-900 bg-transparent text-ink-900 hover:bg-primary-50 disabled:cursor-not-allowed disabled:border-ink-400 disabled:text-ink-400 disabled:hover:bg-transparent',
      variant === 'outline-inverted' &&
        'border-2 border-white bg-transparent text-white hover:bg-white/10 disabled:cursor-not-allowed disabled:border-ink-400 disabled:text-ink-400 disabled:hover:bg-transparent',
    ]"
    @click="onClick"
  >
    <svg v-if="loading" class="size-16 animate-spin" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8v4a4 4 0 0 0-4 4H4z" />
    </svg>
    <slot />
  </component>
</template>
