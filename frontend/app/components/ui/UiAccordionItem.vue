<script setup lang="ts">
import { Plus } from "lucide-vue-next";
import { computed, inject, useId } from "vue";

/**
 * One row of a <UiAccordion>. Must be rendered inside <UiAccordion>. Uses a
 * CSS grid-template-rows 0fr→1fr transition (no JS height measurement); the
 * +/× icon rotates 45deg on open.
 *
 * @example
 * <UiAccordionItem id="lesson-1" title="Занятие 1. Введение в профессию">
 *   <p>Темы: ...</p>
 * </UiAccordionItem>
 */
const props = defineProps<{
  id: string | number;
  title: string;
}>();

const accordion = inject<{
  isOpen: (id: string | number) => boolean;
  toggle: (id: string | number) => void;
}>("uiAccordion");

if (!accordion) {
  throw new Error("UiAccordionItem must be used inside a UiAccordion");
}

const buttonId = useId();
const panelId = useId();
const open = computed(() => accordion.isOpen(props.id));
</script>

<template>
  <div class="overflow-hidden rounded-md bg-surface">
    <h3 class="m-0">
      <button
        :id="buttonId"
        type="button"
        class="flex w-full items-center justify-between gap-16 px-24 py-24 text-left font-display text-h4 text-ink-900 sm:px-32"
        :aria-expanded="open"
        :aria-controls="panelId"
        @click="accordion.toggle(id)"
      >
        <span>{{ title }}</span>
        <span
          class="grid size-24 shrink-0 place-items-center rounded-full border-[1.5px] border-primary-500 text-primary-700 transition-transform duration-200 motion-reduce:transition-none"
          :class="open ? 'rotate-45' : ''"
          aria-hidden="true"
        >
          <Plus class="size-16" />
        </span>
      </button>
    </h3>
    <div
      :id="panelId"
      role="region"
      :aria-labelledby="buttonId"
      class="grid transition-[grid-template-rows] duration-200 ease-out motion-reduce:transition-none"
      :style="{ gridTemplateRows: open ? '1fr' : '0fr' }"
    >
      <div class="overflow-hidden">
        <div class="px-24 pb-24 text-body text-ink-700 sm:px-32">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>
