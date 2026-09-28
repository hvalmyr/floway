<script setup lang="ts">
import type {
  CourseBlockDisplayStyle,
  CourseSectionWithCourses,
  CustomDisplayStyle,
} from "~/types/api";

interface SectionCard {
  key: string;
  name: string;
  blockLabel?: string;
  lessonCount?: string;
  timeLength?: string;
  price?: string;
  coverImage?: string;
  displayStyle: CourseBlockDisplayStyle;
  customColors?: CustomDisplayStyle;
  to: string;
}

/**
 * Homepage course-sections block — see HomeFeaturesSection.vue's doc
 * comment for the pattern. `id="courses"` stays on the first sub-section
 * regardless of where this whole block sits in the admin-managed order,
 * since Hero's "Курсы" button always scrolls to `/#courses`.
 */
defineProps<{
  courseSections: CourseSectionWithCourses[];
  sectionCards: (section: CourseSectionWithCourses) => SectionCard[];
}>();
</script>

<template>
  <section
    v-for="(section, sIndex) in courseSections"
    :id="sIndex === 0 ? 'courses' : undefined"
    :key="section.id"
    class="py-48 sm:py-64 lg:py-80"
    :class="sIndex === 0 ? 'scroll-mt-64 lg:scroll-mt-96' : ''"
  >
    <div class="container flex flex-col gap-48">
      <SectionHeading :color="sIndex % 2 === 0 ? 'primary' : 'ink'">
        {{ section.heading }}
        <template #lead>{{ section.description }}</template>
      </SectionHeading>
      <div class="flex flex-wrap justify-center gap-24 lg:gap-32">
        <CourseCard
          v-for="card in sectionCards(section)"
          :key="card.key"
          :name="card.name"
          :display-style="card.displayStyle"
          :custom-colors="card.customColors"
          :block-label="card.blockLabel"
          :lesson-count="card.lessonCount"
          :time-length="card.timeLength"
          :price="card.price"
          :cover-image="card.coverImage"
          :to="card.to"
          class="w-full sm:w-[calc(50%-12px)] lg:w-[calc(33.333%-22px)]"
        />
      </div>
    </div>
  </section>
</template>
