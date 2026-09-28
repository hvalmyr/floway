<script setup lang="ts">
import type { FAQItem } from "~/types/api";

/** Homepage FAQ block — see HomeFeaturesSection.vue's doc comment for the pattern. */
defineProps<{
  faqItems: FAQItem[];
}>();

const openFaqIds = ref<Array<string | number>>([0]);
</script>

<template>
  <section class="py-48 sm:py-64 lg:py-80">
    <div class="container flex flex-col gap-48">
      <SectionHeading color="primary">
        Вопросы и ответы
        <template #lead>Отвечаем на часто задаваемые вопросы.</template>
      </SectionHeading>
      <UiAccordion v-model="openFaqIds">
        <UiAccordionItem v-for="(item, i) in faqItems" :key="i" :id="i" :title="item.question">
          <MarkdownContent :source="item.answer" />
        </UiAccordionItem>
      </UiAccordion>
    </div>
  </section>
</template>
