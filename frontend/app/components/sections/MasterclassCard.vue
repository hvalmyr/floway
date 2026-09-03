<script setup lang="ts">
import { computed } from "vue";
import type { CourseBlockDisplayStyle, Masterclass } from "~/types/api";

/**
 * Masterclass row for the /masterclasses list — one masterclass per row
 * (not a grid), image (always on the left, sm+) and copy side by side so
 * each entry reads like an editorial feature rather than a cramped card.
 * Stacks image-then-copy on mobile.
 *
 * The button's top margin is a fixed gap, not a `mt-auto`-to-the-bottom
 * trick — a masterclass with a long description/description2/endingText
 * stack can already be taller than the square image next to it, leaving no
 * leftover flex space for `mt-auto` to push against (the button would then
 * end up jammed right under the last paragraph with no gap at all).
 *
 * The image's `min-h-0` matters below `sm` (column layout, no explicit
 * shrink-0 there): flex items default to `min-height: auto`, which for a
 * cover photo taller than the `aspect-[4/5]` box lets its own intrinsic
 * ratio win over the aspect-ratio CSS and stretch the box — see the same
 * fix's comment in CourseCard.vue for a confirmed-live repro.
 *
 * @example
 * <MasterclassCard :masterclass="mc" display-style="beige-blue" @apply="..." />
 */
const props = withDefaults(
  defineProps<{
    masterclass: Masterclass;
    displayStyle?: CourseBlockDisplayStyle;
  }>(),
  { displayStyle: "beige-blue" },
);

// Pulled out of the template as computeds rather than inline compound
// conditions -- `v-if="masterclass.duration || masterclass.price"` on a
// <p> with 3 sibling children (span/div/span) made Vue's SSR/hydration
// produce an empty, class-less paragraph on the server while the client
// rendered it correctly, throwing off hydration alignment for every
// paragraph after it (description2/endingText showed up as "hydration
// text mismatch" purely as a downstream symptom of this). Named computeds
// for both boolean expressions sidestep whatever compiler edge case those
// inline `||`/`&&` expressions hit.
const hasDurationOrPrice = computed(
  () => !!(props.masterclass.duration || props.masterclass.price),
);
const hasDurationAndPrice = computed(
  () => !!(props.masterclass.duration && props.masterclass.price),
);

// The page has ONE shared ApplyForm below every row (not one per card) —
// this tells the page which masterclass's "Записаться" was actually
// clicked, so it can pass that slug into the form (see masterclasses/index.vue).
defineEmits<{ apply: [] }>();

// Masterclass-specific palette — diverges from the shared
// displayStyleColorClasses (course cards keep the plain beige): the
// "beige-blue" style reads as blue text on white here instead of beige,
// and both variants get the site's glass treatment.
const colorClasses: Record<CourseBlockDisplayStyle, string> = {
  "blue-beige": "bg-primary/70 text-surface backdrop-blur backdrop-saturate-150",
  "brown-beige": "bg-ink/70 text-surface backdrop-blur backdrop-saturate-150",
  "beige-blue":
    "border-2 border-primary bg-white/55 text-primary backdrop-blur backdrop-saturate-150",
  "beige-brown": "border-2 border-ink bg-surface/55 text-ink backdrop-blur backdrop-saturate-150",
};
</script>

<template>
  <div
    class="flex flex-col gap-32 rounded-md p-40 sm:flex-row sm:gap-48"
    :class="colorClasses[displayStyle]"
  >
    <NuxtImg
      v-if="masterclass.coverImage"
      :src="resolveOptimizedMediaUrl(masterclass.coverImage)"
      format="webp"
      :alt="masterclass.title"
      class="aspect-[4/5] w-full min-h-0 rounded-sm object-cover sm:w-[38%] sm:shrink-0"
      sizes="400:100vw sm:38vw"
      loading="lazy"
    />
    <div
      v-else
      class="aspect-[4/5] w-full rounded-sm border-2 border-current sm:w-[38%] sm:shrink-0"
    />

    <div class="flex flex-1 flex-col">
      <h3 class="font-display text-h2">{{ masterclass.title }}</h3>
      <p
        v-if="hasDurationOrPrice"
        class="mt-16 flex flex-row items-center gap-8 font-body text-body font-bold"
      >
        <span v-if="masterclass.duration">{{ masterclass.duration }}</span>
        <span
          v-if="hasDurationAndPrice"
          aria-hidden="true"
          class="inline-block size-[5px] shrink-0 rounded-full bg-current"
        />
        <span v-if="masterclass.price">{{ masterclass.price }}</span>
      </p>
      <p class="mt-24 whitespace-pre-line font-body text-body">{{ masterclass.description }}</p>
      <p v-if="masterclass.description2" class="mt-24 whitespace-pre-line font-body text-body">
        {{ masterclass.description2 }}
      </p>
      <!-- Known issue: this paragraph (the last conditional element before
      <UiButton>) still triggers an SSR/hydration mismatch warning for
      masterclasses with endingText set, on every load — a stable wrapper
      div here didn't fix it either. Looks like a Vue 3.5/Nuxt 4 SSR edge
      case around a conditional element adjacent to a component boundary
      inside a v-for. Vue itself marks this class of mismatch "check-only"
      (production doesn't re-render for it), so it's a console-only
      annoyance today, not a visible bug — see git history/PR discussion
      for the investigation before attempting another fix here. -->
      <p v-if="masterclass.endingText" class="mt-24 whitespace-pre-line font-body text-body italic">
        {{ masterclass.endingText }}
      </p>

      <UiButton
        variant="outline"
        transparent
        to="#apply"
        class="mt-32 w-full justify-center sm:w-auto"
        @click="$emit('apply')"
        >Записаться</UiButton
      >
    </div>
  </div>
</template>
