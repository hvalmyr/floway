<script setup lang="ts">
import type { Teacher } from "~/types/api";

/** Homepage "Педагоги" block — see HomeFeaturesSection.vue's doc comment for the pattern. */
defineProps<{
  teachers: Teacher[];
  glassClass: string;
}>();

function capitalizeName(name: string): string {
  return name
    .split(" ")
    .map((part) => (part ? part[0]!.toUpperCase() + part.slice(1) : part))
    .join(" ");
}
</script>

<template>
  <section class="py-48 sm:py-64 lg:py-80">
    <div class="container flex flex-col gap-48">
      <SectionHeading>Педагоги</SectionHeading>
      <div class="flex flex-wrap justify-center gap-24">
        <div
          v-for="(teacher, index) in teachers"
          :key="teacher.id"
          class="flex w-full flex-col items-center gap-16 md:w-[calc((100%-48px)/3)]"
        >
          <UiContentImage
            v-if="teacher.photo"
            :src="resolveOptimizedMediaUrl(teacher.photo)"
            :alt="teacher.name"
            class="aspect-square w-full min-h-0 rounded-lg object-cover"
            sizes="400:100vw md:33vw"
          />
          <div
            v-else
            class="aspect-square w-full rounded-lg"
            :class="index % 2 === 0 ? 'bg-primary' : 'bg-surface'"
          />
          <!-- Тот же размер, что и у заголовков преимуществ, текста кнопок
          и вопросов FAQ (text-h4), но шрифт Non Bureau (не Soyuz Grotesk) и
          Medium, а не Bold. -->
          <p
            class="w-full rounded-md py-12 text-center font-body text-h4 font-medium"
            :class="[glassClass, index % 2 === 0 ? 'text-primary' : 'text-ink']"
          >
            {{ capitalizeName(teacher.name) }}
          </p>
        </div>
      </div>
    </div>
  </section>
</template>
