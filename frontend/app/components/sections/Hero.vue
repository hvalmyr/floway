<script setup lang="ts">
 /**
 * Shared hero layout for all 4 page types: H1 + lead + CTAs on the left,
 * photo/video media on the right; stacks to text → CTA → media below `lg`.
 *
 * @example
 * <Hero>
 *   <template #title>Мы рядом с первого букета</template>
 *   <template #lead>Обучаем современной флористике с нуля — бережно, понятно, с практикой и поддержкой.</template>
 *   <template #actions>
 *     <UiButton variant="primary" size="lg" to="/#courses">Курсы</UiButton>
 *     <UiButton variant="outline" size="lg" to="/#about">О школе</UiButton>
 *   </template>
 * </Hero>
 */
</script>

<template>
  <section class="py-64 sm:py-96 lg:py-120">
    <div class="container flex flex-col gap-32 lg:flex-row lg:items-center lg:gap-64">
      <div class="flex flex-col items-start gap-24 lg:w-1/2">
        <h1 class="font-display text-h1 text-ink">
          <slot name="title" />
        </h1>
        <p v-if="$slots.lead" class="font-body text-body-l text-ink">
          <slot name="lead" />
        </p>
        <div v-if="$slots.actions" class="flex w-full flex-col gap-16 sm:w-auto sm:flex-row">
          <slot name="actions" />
        </div>
      </div>
      <div class="lg:w-1/2">
        <slot name="media">
          <UiMediaPlaceholder aspect="4/3" />
        </slot>
      </div>
    </div>
  </section>
</template>
