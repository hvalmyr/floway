<script setup lang="ts">
const route = useRoute();
const slug = route.params.slug as string;

const api = useApi();
const { data: course } = await useAsyncData(`course-${slug}`, () => api.getCourse(slug));

if (!course.value) {
  throw createError({ statusCode: 404, statusMessage: "Курс не найден", fatal: true });
}

useSeoMeta({
  title: () => `${course.value?.title} — ФлоВей`,
  description: () => course.value?.shortDescription,
});

// Первое занятие каждого блока открыто по умолчанию.
const openLessonIds = ref<Record<number, Array<string | number>>>(
  Object.fromEntries(
    course.value.modules.map((m) => [m.id, m.lessons.length ? [m.lessons[0].id] : []]),
  ),
);

function priceRow(module: (typeof course.value.modules)[number]) {
  const lessonsWord = module.lessonsCount === 1 ? "занятие" : "занятий";
  const price = `${module.price.toLocaleString("ru-RU")} ₽`;
  return [
    `Блок «${module.title}»`,
    `${module.lessonsCount} ${lessonsWord}`,
    `${module.hours} часов`,
    module.oldPrice ? `${module.oldPrice.toLocaleString("ru-RU")} ₽ → ${price}` : price,
  ];
}
</script>

<template>
  <div v-if="course">
    <Hero>
      <template #title>Курс «{{ course.title }}»</template>
      <template #lead>{{ course.fullDescription }}</template>
      <template #actions>
        <UiButton variant="primary" size="lg" to="#apply" class="w-full sm:w-auto"
          >Оставить заявку</UiButton
        >
      </template>
    </Hero>

    <section v-if="course.modules.length" class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-24">
        <UiInfoRow
          v-for="(module, i) in course.modules"
          :key="module.id"
          :items="priceRow(module)"
          :highlighted="i % 2 === 1"
        />
      </div>
    </section>

    <section v-if="course.modules.length" class="bg-surface py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-64">
        <div v-for="module in course.modules" :key="module.id" class="flex flex-col gap-24">
          <div class="flex flex-col gap-16">
            <h2 class="font-display text-h2">
              <span class="text-primary">Учебный план.</span>
              <span class="text-ink">Блок «{{ module.title }}»</span>
            </h2>
            <p v-if="module.description" class="font-body text-body-l text-ink">
              {{ module.description }}
            </p>
          </div>
          <UiAccordion v-model="openLessonIds[module.id]">
            <UiAccordionItem
              v-for="lesson in module.lessons"
              :key="lesson.id"
              :id="lesson.id"
              :title="lesson.title"
            >
              <p class="mb-8">
                <strong class="font-display text-ink">Темы:</strong> {{ lesson.topics }}
              </p>
              <p class="mb-8">
                <strong class="font-display text-ink">Вы научитесь:</strong> {{ lesson.outcomes }}
              </p>
              <p>
                <strong class="font-display text-ink">Продолжительность:</strong>
                {{ lesson.durationHours }} часа.
              </p>
            </UiAccordionItem>
          </UiAccordion>
        </div>
      </div>
    </section>
    <section v-else class="py-64">
      <div class="container">
        <p class="font-body text-body text-ink">
          Программа курса пока готовится — оставьте заявку, и мы расскажем подробности при звонке.
        </p>
      </div>
    </section>

    <section id="apply" class="scroll-mt-64 py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm context="course" :related-id="course.id" title="Записаться на курс" />
        </div>
      </div>
    </section>
  </div>
</template>
