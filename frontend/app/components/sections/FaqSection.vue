<script setup lang="ts">
/**
 * A single page's own FAQ block — title, intro text and a Q&A accordion.
 * Distinct from the site-wide FAQ on the homepage (see pages/index.vue),
 * which has no title/description/visibility of its own. Used by the course
 * page (CourseFAQItem, faqTitle/faqDescription/faqVisible) and by the
 * masterclasses/gift-certificate pages (PageFaq, fetched via
 * api.getPageFaq()) — both shapes carry the same id/question/answer fields,
 * so a structural type here covers either without importing one specific to
 * the other. The caller is responsible for only mounting this when
 * visible && items.length.
 */
const props = defineProps<{
  title: string;
  description: string;
  items: { id: number; question: string; answer: string }[];
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
