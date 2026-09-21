<script setup lang="ts">
useSeoMeta({
  title: "Как мы работаем — flo-way",
  description: "Этапы новогоднего оформления: от заявки до демонтажа.",
});

const steps = [
  {
    title: "Заявка",
    description: "Оставляете заявку на сайте или пишете в мессенджер — уточняем детали и формат.",
  },
  {
    title: "Бесплатный выезд и замеры",
    description: "Приезжаем на объект, оцениваем пространство и снимаем размеры.",
  },
  {
    title: "Визуал",
    description: "Готовим визуализацию будущего оформления, чтобы вы увидели результат заранее.",
  },
  {
    title: "Коммерческое предложение",
    description: "Присылаем предложение в виде презентации — с составом работ и стоимостью.",
  },
  {
    title: "Согласование",
    description: "Обсуждаем детали и вносим правки, пока результат не устроит полностью.",
  },
  {
    title: "Монтаж",
    description: "Устанавливаем оформление в согласованные сроки, аккуратно и без лишнего шума.",
  },
  {
    title: "Демонтаж",
    description: "Снимаем оформление после сезона или праздника — бережно и без следов.",
  },
];
</script>

<template>
  <div class="container flex flex-col gap-48 py-48 sm:py-64 lg:py-80">
    <SectionHeading color="ink">
      Как мы работаем
      <template #lead>Семь понятных этапов — от заявки до демонтажа.</template>
    </SectionHeading>

    <div class="flex flex-col gap-24">
      <div
        v-for="(step, i) in steps"
        :key="step.title"
        class="flex flex-col gap-8 rounded-md p-24"
        :class="i % 2 === 1 ? 'border-2 border-primary bg-white' : 'bg-surface'"
      >
        <div class="flex items-center gap-16">
          <span class="font-display text-h4 text-primary">{{ i + 1 }}</span>
          <h2 class="font-display text-h4 text-ink">{{ step.title }}</h2>
        </div>
        <p class="font-body text-body text-ink">{{ step.description }}</p>
      </div>
    </div>

    <div class="mx-auto max-w-[720px] w-full">
      <ApplyForm title="Оставить заявку" />
    </div>
  </div>
</template>
