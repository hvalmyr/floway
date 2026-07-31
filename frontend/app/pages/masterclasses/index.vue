<script setup lang="ts">
import { Clock, Flower2, Gift, Palette, Sparkles, Users } from "lucide-vue-next";

useSeoMeta({
  title: "Мастер-классы по флористике в Москве — ФлоВей",
  description: "Мастер-классы по флористике в свободном графике: букеты и композиции, все материалы включены.",
});

const api = useApi();
const { data: masterclasses } = await useAsyncData("masterclasses-list", () => api.getMasterClasses());

// TODO: набор преимуществ ниже придуман по аналогии с блоком на главной
// ("аналогично главной, другой набор" в брифе, без конкретных формулировок) —
// сверить с заказчиком перед публикацией.
const features = [
  { icon: Sparkles, title: "Без подготовки", description: "Не нужно ничего уметь заранее — всему научим на месте." },
  { icon: Flower2, title: "Только живые цветы", description: "Работаете с настоящим материалом, а не с имитацией." },
  { icon: Gift, title: "Готовый подарок", description: "Забираете работу с собой — отличный повод для мастер-класса перед праздником." },
  { icon: Clock, title: "Один день — один результат", description: "Приходите на 2–4 часа и уходите с готовой работой." },
  { icon: Users, title: "Одному или компанией", description: "Подходит и для индивидуального занятия, и для группы друзей." },
  { icon: Palette, title: "Простор для творчества", description: "Помогаем с техникой, но форму и настроение выбираете вы." },
];
</script>

<template>
  <div>
    <Hero>
      <template #title>Мастер-классы по флористике в свободном графике</template>
      <template #lead>
        Разовое занятие без предварительной подготовки: соберёте букет или композицию и заберёте с собой.
      </template>
      <template #actions>
        <UiButton variant="primary" size="lg" to="#apply" class="w-full sm:w-auto">Оставить заявку</UiButton>
      </template>
    </Hero>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>Почему стоит выбрать мастер-класс «ФлоВей»</SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section class="bg-surface py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-64">
        <MasterclassCard v-for="mc in masterclasses" :key="mc.id" :masterclass="mc" link-to-detail />
      </div>
    </section>

    <section id="apply" class="scroll-mt-64 py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm context="masterclass" title="Оставить заявку на мастер-класс" />
        </div>
      </div>
    </section>
  </div>
</template>
