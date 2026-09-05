<script setup lang="ts">
import type { CourseFAQItem } from "~/types/api";

/**
 * A single course's own FAQ block — distinct from the site-wide FAQ on the
 * homepage (see pages/index.vue). Rendered on the course page right after
 * the apply form; the caller (courses/[slug].vue) is responsible for only
 * mounting this when course.faqVisible && course.faqItems.length.
 */
const props = defineProps<{
  title: string;
  description: string;
  items: CourseFAQItem[];
}>();

const openIds = ref<Array<string | number>>(props.items.length ? [props.items[0]!.id] : []);
</script>

<template>
  <div class="flex flex-col gap-48">
    <SectionHeading color="primary">
      {{ title }}
      <template v-if="description" #lead>{{ description }}</template>
    </SectionHeading>
    <UiAccordion v-model="openIds">
      <UiAccordionItem v-for="item in items" :key="item.id" :id="item.id" :title="item.question">
        <MarkdownContent :source="item.answer" />
      </UiAccordionItem>
    </UiAccordion>
  </div>
</template>
