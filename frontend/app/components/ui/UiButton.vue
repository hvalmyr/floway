<script setup lang="ts">
import { NuxtLink } from "#components";

/**
 * Button/link primitive. Renders a <NuxtLink> when `to` is set, otherwise a
 * <button>. Exactly one size/font everywhere — only bg/border/text color
 * (and hover behavior) differ between the two variants:
 *
 * - `primary`: white text on blue fill. Hover inverts to blue text on white
 *   fill; the blue border is there at rest too, just the same color as the
 *   fill so it's invisible until the fill flips to white.
 * - `outline` (brand brown): brown text on white fill with a brown border.
 *   Hover fills brown with white text and drops the border.
 *
 * `pill` still controls shape (stadium vs. large-rounded-rect) since that's
 * not restricted to one value — nav/hero/card CTAs are pill, form submit
 * buttons use `:pill="false"`.
 *
 * `:is` takes the imported NuxtLink component, not the string "NuxtLink" —
 * a bare string doesn't resolve at runtime (Nuxt's auto-import only
 * rewrites literal <NuxtLink> tags at compile time, so `resolveDynamicComponent`
 * can't find it), which would silently render an inert <nuxtlink> custom
 * element instead of a real link on every `to`-having button on the site.
 *
 * @example
 * <UiButton variant="primary" to="/masterclasses">Мастер-классы</UiButton>
 * <UiButton variant="outline" to="/#about">О школе</UiButton>
 * <UiButton type="submit" :pill="false" block :loading="pending">Отправить заявку</UiButton>
 */
const props = withDefaults(
  defineProps<{
    variant?: "primary" | "outline";
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
    :is="to ? NuxtLink : 'button'"
    :to="to"
    :type="to ? undefined : type"
    :disabled="!to && (disabled || loading)"
    :aria-busy="loading || undefined"
    :aria-disabled="to && (disabled || loading) ? 'true' : undefined"
    class="inline-flex h-[44px] items-center justify-center gap-8 border-2 px-24 font-display text-button font-bold transition-colors duration-200 sm:h-[64px] sm:px-40"
    :class="[
      pill ? 'rounded-pill' : 'rounded-lg',
      block ? 'w-full' : '',
      (disabled || loading) && 'cursor-not-allowed opacity-40',
      variant === 'primary' &&
        'border-primary bg-primary text-white hover:bg-white hover:text-primary',
      variant === 'outline' &&
        'border-ink bg-white text-ink hover:border-transparent hover:bg-ink hover:text-white',
    ]"
    @click="onClick"
  >
    <svg
      v-if="loading"
      class="size-16 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8v4a4 4 0 0 0-4 4H4z" />
    </svg>
    <slot />
  </component>
</template>
