<script setup lang="ts">
/**
 * Button/link primitive. Renders a <NuxtLink> when `to` is set, otherwise a
 * <button>. Nav/hero/card CTAs are full pill (default); form submit buttons
 * use `:pill="false"` for the large-rounded-rect shape seen in the apply
 * forms mockups.
 *
 * @example
 * <UiButton variant="primary" size="lg" to="/masterclasses">Мастер-классы</UiButton>
 * <UiButton variant="outline" size="md">О школе</UiButton>
 * <UiButton type="submit" :pill="false" block :loading="pending">Отправить заявку</UiButton>
 */
const props = withDefaults(
  defineProps<{
    variant?: "primary" | "outline" | "outline-inverted";
    size?: "md" | "lg";
    to?: string;
    type?: "button" | "submit" | "reset";
    disabled?: boolean;
    loading?: boolean;
    /** Always full-width, regardless of breakpoint (used for apply-form submit buttons). */
    block?: boolean;
    /** Stadium/pill shape (nav, hero, card CTAs). Set false for the flatter form-submit shape. */
    pill?: boolean;
  }>(),
  {
    variant: "primary",
    size: "md",
    to: undefined,
    type: "button",
    disabled: false,
    loading: false,
    block: false,
    pill: true,
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
    class="inline-flex items-center justify-center gap-8 font-display text-button font-bold transition-[background-color,transform,border-color,opacity] duration-200"
    :class="[
      pill ? 'rounded-pill' : 'rounded-lg',
      block ? 'w-full' : '',
      size === 'lg' ? 'h-[44px] px-24 sm:h-[64px] sm:px-40' : 'h-[44px] px-24 sm:h-[56px] sm:px-32',
      (disabled || loading) && 'cursor-not-allowed opacity-40 hover:translate-y-0',
      variant === 'primary' && 'bg-primary text-white hover:-translate-y-px active:translate-y-0',
      variant === 'outline' && 'border-2 border-ink bg-transparent text-ink hover:bg-surface',
      variant === 'outline-inverted' && 'border-2 border-white bg-transparent text-white hover:bg-white/10',
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
