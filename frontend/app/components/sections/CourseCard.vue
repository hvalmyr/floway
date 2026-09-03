<script setup lang="ts">
import { NuxtLink } from "#components";
import { displayStyleColorClasses } from "~/constants/display-style-colors";
import type { CourseBlockDisplayStyle } from "~/types/api";

/**
 * Course card used in the homepage's course-sections block — one instance
 * per visible block (the caller, index.vue's `sectionCards()`, produces one
 * card per block; a course with a single block, real or synthetic, just
 * means one call). Order: title (the course's name — same across every
 * card of a multi-block course), then up to 3 stacked caption lines
 * (blockLabel/lessonCount/timeLength, each own line, shown only when set),
 * then the block's cover, then the CTA.
 *
 * `displayStyle` is the admin-chosen background/text color pair (one of the
 * design system's 4 standard colors combined into a readable pair) — picked
 * per block, not cycled automatically. The card hugs its own content (no
 * fixed min-height); `mt-auto` sits on the cover image (not the button) so
 * the image+CTA pair lands flush with the bottom edge as a unit regardless
 * of how many caption lines precede it — otherwise cards with a different
 * caption-line count in the same row would show the photo at different
 * heights even with the button itself still pinned to the bottom. Cards in
 * the same row end up equal height because the parent's `flex flex-wrap`
 * stretches them by default; this component just no longer forces a height
 * of its own on top of that.
 * The CTA's `transparent` outline pulls its text/border from `currentColor`
 * (see UiButton.vue), so it automatically matches whichever of the 4
 * `colorClasses` below is active instead of needing a color per style here.
 *
 * @example
 * <CourseCard
 *   name="Основы флористики"
 *   block-label='Блок "Букеты"'
 *   display-style="blue-beige"
 *   lesson-count="7 занятий"
 *   time-length="30 часов"
 *   cover-image="courses/buketi.jpg"
 *   to="/courses/osnovy-floristiki"
 * />
 */
const props = withDefaults(
  defineProps<{
    name: string;
    description?: string;
    displayStyle?: CourseBlockDisplayStyle;
    coverImage?: string;
    /** The block's own name/label (e.g. `Блок "Букеты"`) — blank for a
     * synthetic single block, so it's simply omitted for a blockless
     * course. */
    blockLabel?: string;
    lessonCount?: string;
    timeLength?: string;
    to: string;
    ctaLabel?: string;
  }>(),
  {
    displayStyle: "blue-beige",
    ctaLabel: "Узнать больше",
  },
);

const colorClasses = displayStyleColorClasses;
</script>

<template>
  <UiCard variant="custom" :class="[colorClasses[displayStyle], 'flex flex-col']">
    <template #title>{{ name }}</template>

    <div v-if="blockLabel || lessonCount || timeLength" class="mb-16 flex flex-col gap-4">
      <p v-if="blockLabel" class="font-body text-body">{{ blockLabel }}</p>
      <p
        v-if="lessonCount || timeLength"
        class="flex flex-row items-center gap-8 font-body text-body"
      >
        <span v-if="lessonCount">{{ lessonCount }}</span>
        <span
          v-if="lessonCount && timeLength"
          aria-hidden="true"
          class="inline-block size-[5px] shrink-0 rounded-full bg-current"
        />
        <span v-if="timeLength">{{ timeLength }}</span>
      </p>
    </div>
    <p v-if="description" class="mb-24 whitespace-pre-line font-body text-body">
      {{ description }}
    </p>

    <component :is="NuxtLink" :to="to" class="mt-auto mb-24 block aspect-square w-full">
      <NuxtImg
        v-if="coverImage"
        :src="resolveOptimizedMediaUrl(coverImage)"
        format="webp"
        :alt="name"
        class="size-full rounded-sm object-cover"
        sizes="400:100vw sm:50vw lg:400px"
        loading="lazy"
      />
      <div v-else class="size-full rounded-sm border-2 border-current" />
    </component>

    <UiButton variant="outline" transparent :to="to" class="w-full justify-center">{{
      ctaLabel
    }}</UiButton>
  </UiCard>
</template>
