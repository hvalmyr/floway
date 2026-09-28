<script setup lang="ts">
import type { AboutItem } from "~/types/api";

/** Homepage "О школе" block — see HomeFeaturesSection.vue's doc comment for the pattern. */
defineProps<{
  aboutItems: AboutItem[];
}>();
</script>

<template>
  <section id="about" class="scroll-mt-64 py-48 sm:py-64 lg:scroll-mt-96 lg:py-80">
    <div class="container flex flex-col gap-48">
      <SectionHeading>О школе</SectionHeading>
      <!-- Бежевая карточка-обёртка вокруг списка (flex-column, не грид);
      каждый пункт — отдельная белая строка, растянутая на всю ширину
      обёртки. -->
      <div class="rounded-lg bg-surface/55 p-24 backdrop-blur backdrop-saturate-150 sm:p-32">
        <div class="flex flex-col items-start gap-24">
          <div
            v-for="item in aboutItems"
            :key="item.id"
            class="flex w-full flex-col items-start gap-16 rounded-md bg-white p-32"
          >
            <UiBadge>{{ item.badge }}</UiBadge>
            <p class="whitespace-pre-line font-body text-body text-ink">
              {{ item.description }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
