<script setup lang="ts">
 /**
 * Shared hero layout for all 4 page types: H1 + lead + CTAs on the left,
 * 1:1 photo/video on the right; below `lg` (tablets and phones) it stacks
 * to media → text → CTA instead, via `order` utilities. CTAs always stack
 * full-width in a column (never side by side), per the mockup, and sit
 * flush with the media's bottom edge on `lg`+ (the text column stretches to
 * the media's height via `items-stretch` + `mt-auto` on the actions
 * wrapper). `.container` caps overall width so this doesn't stretch
 * awkwardly on very wide screens.
 *
 * @example
 * <Hero>
 *   <template #title>Мы рядом с первого букета</template>
 *   <template #lead>Обучаем современной флористике с нуля — бережно, понятно, с практикой и поддержкой.</template>
 *   <template #actions>
 *     <UiButton variant="primary" to="/#courses">Курсы</UiButton>
 *     <UiButton variant="outline" to="/#trial">Пробное занятие</UiButton>
 *   </template>
 * </Hero>
 */
</script>

<template>
  <section class="py-48 sm:py-64 lg:py-80">
    <div class="container flex flex-col gap-32 lg:flex-row lg:items-stretch lg:gap-64">
      <div class="order-2 flex flex-col items-start gap-24 lg:order-1 lg:w-1/2">
        <h1 class="font-display text-h1 text-ink">
          <slot name="title" />
        </h1>
        <p
          v-if="$slots.lead"
          class="w-full whitespace-pre-line rounded-md bg-white/55 px-24 py-16 font-body text-body text-ink backdrop-blur backdrop-saturate-150"
        >
          <slot name="lead" />
        </p>
        <div v-if="$slots.actions" class="mt-auto flex w-full flex-col gap-16">
          <slot name="actions" />
        </div>
      </div>
      <div class="order-1 lg:order-2 lg:w-1/2">
        <slot name="media">
          <UiMediaPlaceholder aspect="1/1" />
        </slot>
      </div>
    </div>
  </section>
</template>
