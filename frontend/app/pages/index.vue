<script setup lang="ts">
import { Armchair, CalendarClock, Flower2, GraduationCap, Layers, MapPin } from "lucide-vue-next";

useSeoMeta({
  title: "ФлоВей — школа флористики в Москве",
  description: "Курсы и мастер-классы по флористике в Москве: с нуля до уверенной самостоятельной работы.",
});

const api = useApi();

const { data: mainCourse } = await useAsyncData("home-course-osnovy", () => api.getCourse("osnovy-floristiki"));
const { data: courses } = await useAsyncData("home-courses", () => api.getCourses());

const profileSlugs = ["kommercheskaya-floristika", "aktualnaya-floristika", "svadebnaya-floristika"];
const profileVariants = ["surface-ink", "surface-primary", "surface-white"] as const;
const profileCourses = computed(
  () => courses.value?.filter((c) => profileSlugs.includes(c.slug)).sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

// TODO: тексты ниже написаны по смыслу брифа (иконки/заголовки/лиды заданы,
// короткие описания — нет: "текст взять из макета/аналогичный") — сверить с
// заказчиком перед публикацией.
const features = [
  {
    icon: Layers,
    title: "Два формата обучения",
    description: "Групповые занятия или свободный график — выбираете сами, что удобнее.",
  },
  {
    icon: Armchair,
    title: "Комфортное обучение",
    description: "Светлая студия, все материалы и инструменты — приходите налегке.",
  },
  {
    icon: MapPin,
    title: "Центр Москвы",
    description: "Учимся в шаговой доступности от метро, в самом центре города.",
  },
  {
    icon: CalendarClock,
    title: "Гибкий старт",
    description: "Новые потоки стартуют регулярно — не нужно ждать месяцами.",
  },
  {
    icon: Flower2,
    title: "Максимум практики",
    description: "Большую часть времени вы работаете с живыми цветами, а не с теорией.",
  },
  {
    icon: GraduationCap,
    title: "Для любого уровня",
    description: "Берём с нуля и помогаем дойти до уверенной самостоятельной работы.",
  },
];

const aboutItems = [
  { badge: "Есть свободный формат", description: "Занимайтесь в своём темпе, подстраивая расписание под себя." },
  { badge: "Есть групповой формат", description: "Учитесь вместе с другими и обменивайтесь опытом на занятиях." },
  { badge: "Маленькие группы до 7 человек", description: "Больше внимания преподавателя каждому ученику." },
  { badge: "Только живые цветы", description: "Никаких имитаций — работаем с настоящим материалом с первого занятия." },
];
</script>

<template>
  <div>
    <Hero>
      <template #title>Мы рядом с первого букета</template>
      <template #lead>
        Обучаем современной флористике с нуля — бережно, понятно, с практикой и поддержкой.
      </template>
      <template #actions>
        <UiButton variant="primary" size="lg" to="/#courses" class="w-full sm:w-auto">Курсы</UiButton>
        <UiButton variant="outline" size="lg" to="/#about" class="w-full sm:w-auto">О школе</UiButton>
      </template>
    </Hero>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>
          Почему стоит учиться в школе «ФлоВей»
          <template #lead>Собрали то, что действительно важно перед стартом обучения.</template>
        </SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section id="about" class="scroll-mt-64 bg-surface py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <div class="flex flex-col items-center gap-16 text-center">
          <UiBadge>работаем с 2013 года</UiBadge>
          <SectionHeading>
            О школе
            <template #lead>
              Больше десяти лет учим флористике так, чтобы результат был виден уже на первом занятии.
            </template>
          </SectionHeading>
        </div>
        <div class="grid grid-cols-1 gap-24 md:grid-cols-2">
          <div v-for="item in aboutItems" :key="item.badge" class="flex flex-col items-start gap-16 rounded-md bg-white p-32">
            <UiBadge>{{ item.badge }}</UiBadge>
            <p class="text-body text-ink-700">{{ item.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <section id="courses" class="scroll-mt-64 py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>
          «Основы флористики»
          <template #lead>{{ mainCourse?.shortDescription }}</template>
        </SectionHeading>
        <div v-if="mainCourse" class="grid grid-cols-1 gap-24 md:grid-cols-2 lg:gap-32">
          <CourseCard
            v-for="module in mainCourse.modules"
            :key="module.id"
            :title="`Блок «${module.title}»`"
            :variant="module.title === 'Букеты' ? 'surface-primary' : 'surface-ink'"
            :lessons-count="module.lessonsCount"
            :hours="module.hours"
            :price="module.price"
            :to="`/courses/${mainCourse.slug}`"
          />
        </div>
      </div>
    </section>

    <section class="bg-surface py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>
          Профильные курсы
          <template #lead>Когда основы освоены — можно двигаться вглубь одного направления.</template>
        </SectionHeading>
        <div class="grid grid-cols-1 gap-24 md:grid-cols-2 lg:grid-cols-3 lg:gap-32">
          <CourseCard
            v-for="(course, i) in profileCourses"
            :key="course.id"
            :title="course.title"
            :variant="profileVariants[i % profileVariants.length]"
            :to="`/courses/${course.slug}`"
          />
        </div>
      </div>
    </section>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm context="trial_lesson" title="Записаться на пробное занятие" />
        </div>
      </div>
    </section>
  </div>
</template>
