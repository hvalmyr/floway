<script setup lang="ts">
import type { CourseBlockWithLessons } from "~/types/api";

const route = useRoute();
const slug = route.params.slug as string;

const api = useApi();
const { data: course } = await useAsyncData(`course-${slug}`, () => api.getCourse(slug));

if (!course.value) {
  throw createError({ statusCode: 404, statusMessage: "Курс не найден", fatal: true });
}

useSeoMeta({
  title: () => `${course.value?.name} — Фловей`,
  description: () => course.value?.description,
});

// Undivided single-block courses (the common case) don't show the block's
// own name anywhere — it's only a meaningful label once there's more than
// one block to tell apart. Same rule for both the info-row badges and the
// curriculum heading below.
const hasNamedBlocks = computed(() => course.value!.blockCount > 1);

function badgeItems(block: CourseBlockWithLessons) {
  const parts = hasNamedBlocks.value ? [block.blockName] : [];
  return [...parts, block.lessonCount, block.timeLength, block.price].filter(Boolean);
}

// The hero photo is the one place this page is allowed to mirror the
// homepage's collapsed single-card representation: for a `singleCard`
// course, that homepage card was built from the course's OWN cover (see
// syntheticBlock() in course_catalog_service.go), so the hero should show
// that same photo rather than the first real block's — otherwise a visitor
// lands on a different image than the one they just clicked. Every other
// section below (info-rows, study plans) always uses the real per-block
// breakdown regardless of singleCard.
const heroCover = computed(() =>
  course.value!.singleCard ? course.value!.coverImage : course.value!.blocks[0]?.blockCover,
);

// First lesson of each block open by default, same as the old per-module
// behavior.
const openLessonIds = ref<Record<number, Array<string | number>>>(
  Object.fromEntries(
    course.value!.blocks.map((b) => [b.id, b.lessons.length ? [b.lessons[0]!.id] : []]),
  ),
);
</script>

<template>
  <div v-if="course">
    <Hero>
      <template #title>Курс “{{ course.name }}”</template>
      <template #lead>{{ course.description }}</template>
      <template #actions>
        <UiButton variant="primary" to="#apply">Оставить заявку</UiButton>
      </template>
      <template v-if="heroCover" #media>
        <UiHeroPicture :src="resolveOptimizedMediaUrl(heroCover)" :alt="course.name" />
      </template>
    </Hero>

    <section v-if="course.blocks.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-24">
        <UiInfoRow
          v-for="(block, i) in course.blocks"
          :key="block.id"
          :items="badgeItems(block)"
          :highlighted="i % 2 === 1"
        />
      </div>
    </section>

    <section v-if="course.blocks.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-64">
        <div v-for="block in course.blocks" :key="block.id" class="flex flex-col gap-24">
          <div class="flex flex-col items-center gap-16 text-center">
            <h2 class="font-display text-h2">
              <template v-if="hasNamedBlocks">
                <span class="text-primary">Учебный план. </span>
                <span class="text-ink">{{ block.blockName }}</span>
              </template>
              <span v-else class="text-primary">{{ block.blockName || "Учебный план" }}</span>
            </h2>
            <p
              v-if="block.description"
              class="mx-auto w-full whitespace-pre-line rounded-md bg-white/55 px-24 py-16 font-body text-body text-ink backdrop-blur backdrop-saturate-150 lg:w-4/5"
            >
              {{ block.description }}
            </p>
          </div>
          <UiAccordion v-model="openLessonIds[block.id]!">
            <UiAccordionItem
              v-for="lesson in block.lessons"
              :key="lesson.id"
              :id="lesson.id"
              :title="lesson.name"
            >
              <MarkdownContent :source="lesson.description" />
            </UiAccordionItem>
          </UiAccordion>
        </div>
      </div>
    </section>
    <section v-else class="py-48">
      <div class="container">
        <p class="font-body text-body text-ink">
          Программа курса пока готовится — оставьте заявку, и мы расскажем подробности при звонке.
        </p>
      </div>
    </section>

    <section
      id="apply"
      class="scroll-mt-64 bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:scroll-mt-96 lg:py-80"
    >
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm
            context="course"
            :related-id="course.id"
            :related-slug="course.slug"
            title="Записаться на курс"
          />
        </div>
      </div>
    </section>
  </div>
</template>
