<script setup lang="ts">
import { provide } from "vue";

/**
 * Accordion container. Manages which item(s) are open via a v-model array of
 * ids and shares toggle/isOpen with descendant <UiAccordionItem> via provide.
 *
 * @example
 * <UiAccordion v-model="openLessonIds">
 *   <UiAccordionItem v-for="l in lessons" :key="l.id" :id="l.id" :title="l.title">
 *     {{ l.topics }}
 *   </UiAccordionItem>
 * </UiAccordion>
 */
const props = withDefaults(
  defineProps<{
    modelValue: Array<string | number>;
    allowMultiple?: boolean;
  }>(),
  { allowMultiple: true },
);

const emit = defineEmits<{
  "update:modelValue": [Array<string | number>];
}>();

function isOpen(id: string | number) {
  return props.modelValue.includes(id);
}

function toggle(id: string | number) {
  if (isOpen(id)) {
    emit(
      "update:modelValue",
      props.modelValue.filter((v) => v !== id),
    );
    return;
  }
  emit("update:modelValue", props.allowMultiple ? [...props.modelValue, id] : [id]);
}

provide("uiAccordion", { isOpen, toggle });
</script>

<template>
  <div class="flex flex-col gap-16">
    <slot />
  </div>
</template>
