<script setup lang="ts">
/**
 * Small mockup of a public course card (see CourseCard.vue), rendered with
 * whatever bg/text colors the caller passes — used by the custom-styles
 * admin page (live preview while picking colors) and by the course/course
 * block edit forms (preview of whichever style, custom or standard, is
 * currently selected) so an admin can see the actual card look before
 * saving, not just a color swatch.
 */
withDefaults(
  defineProps<{
    bgColor: string;
    textColor: string;
    title?: string;
  }>(),
  { title: "Название курса" },
);
</script>

<template>
  <div
    class="flex w-full max-w-[260px] flex-col gap-16 rounded-md p-24"
    :style="{ backgroundColor: bgColor, color: textColor }"
  >
    <p class="text-base font-semibold">{{ title }}</p>
    <p class="flex items-center gap-8 text-sm">
      <span>7 занятий</span>
      <span class="inline-block size-[5px] shrink-0 rounded-full bg-current" />
      <span>30 часов</span>
    </p>
    <div class="aspect-square w-full rounded-sm border-2 border-current" />
    <div
      class="w-full rounded-sm border-2 border-current px-16 py-8 text-center text-sm font-medium"
    >
      Узнать больше
    </div>
  </div>
</template>
