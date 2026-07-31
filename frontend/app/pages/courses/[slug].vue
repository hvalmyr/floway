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
  Object.fromEntries(course.value.modules.map((m) => [m.id, m.lessons.length ? [m.lessons[0].id] : []])),
);

function priceInfo(module: (typeof course.value.modules)[number]) {
  return [
    { label: "Занятия", value: String(module.lessonsCount) },
    { label: "Часы", value: String(module.hours) },
    {
      label: "Цена",
      value: `${module.price.toLocaleString("ru-RU")} ₽`,
      oldValue: module.oldPrice ? `${module.oldPrice.toLocaleString("ru-RU")} ₽` : undefined,
    },
  ];
}
</script>

<template>
  <div v-if="course">
    <Hero>
      <template #title>Курс «{{ course.title }}»</template>
      <template #lead>{{ course.fullDescription }}</template>
      <template #actions>
        <UiButton variant="primary" size="lg" to="#apply" class="w-full sm:w-auto">Оставить заявку</UiButton>
      </template>
    </Hero>

    <section v-if="course.modules.length" class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-24">
        <SectionHeading>Стоимость обучения</SectionHeading>
        <div v-for="module in course.modules" :key="module.id" class="flex flex-col gap-8">
          <p class="font-display text-h4 text-ink-900">Блок «{{ module.title }}»</p>
          <UiInfoRow :items="priceInfo(module)" />
        </div>
      </div>
    </section>

    <section v-if="course.modules.length" class="bg-surface py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>
          Учебный план
          <template #lead>Что происходит на занятиях, блок за блоком.</template>
        </SectionHeading>

        <div v-for="module in course.modules" :key="module.id" class="flex flex-col gap-24">
          <h3 class="font-display text-h3 text-ink-900">Блок «{{ module.title }}»</h3>
          <UiAccordion v-model="openLessonIds[module.id]">
            <UiAccordionItem v-for="lesson in module.lessons" :key="lesson.id" :id="lesson.id" :title="lesson.title">
              <p class="mb-8"><strong class="font-display text-ink-900">Темы:</strong> {{ lesson.topics }}</p>
              <p class="mb-8">
                <strong class="font-display text-ink-900">Вы научитесь:</strong> {{ lesson.outcomes }}
              </p>
              <p><strong class="font-display text-ink-900">Продолжительность:</strong> {{ lesson.durationHours }} ч.</p>
            </UiAccordionItem>
          </UiAccordion>
        </div>
      </div>
    </section>
    <section v-else class="py-64">
      <div class="container">
        <p class="text-body text-ink-700">
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
