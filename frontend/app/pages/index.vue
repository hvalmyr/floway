<script setup lang="ts">
import type { Component } from "vue";
import { featureIconComponent } from "~/constants/feature-icons";

useSeoMeta({
  title: "ФлоВей — школа флористики в Москве",
  description:
    "Курсы и мастер-классы по флористике в Москве: с нуля до уверенной самостоятельной работы.",
});

const api = useApi();
const { text } = await usePageContent();

const { data: mainCourse } = await useAsyncData("home-course-osnovy", () =>
  api.getCourse("osnovy-floristiki"),
);
const { data: courses } = await useAsyncData("home-courses", () => api.getCourses());

const profileSlugs = [
  "kommercheskaya-floristika",
  "aktualnaya-floristika",
  "svadebnaya-floristika",
];
const profileVariants = ["surface-ink", "surface-primary", "surface-white"] as const;
const profileCourses = computed(
  () =>
    courses.value
      ?.filter((c) => profileSlugs.includes(c.slug))
      .sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: featuresData } = await useAsyncData("home-features", () => api.getFeatures("home"));
const features = computed(
  () =>
    featuresData.value
      ?.slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((f) => ({
        icon: featureIconComponent(f.icon),
        title: f.title,
        description: f.description,
      }))
      .filter((f): f is typeof f & { icon: Component } => f.icon !== undefined) ?? [],
);

const { data: aboutItemsData } = await useAsyncData("home-about-items", () => api.getAboutItems());
const aboutItems = computed(
  () => aboutItemsData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: teachersData } = await useAsyncData("home-teachers", () => api.getTeachers());
const teachers = computed(
  () => teachersData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: faqData } = await useAsyncData("home-faq", () => api.getFAQItems());
const faqItems = computed(() => faqData.value ?? []);
const openFaqIds = ref<Array<string | number>>([0]);
</script>

<template>
  <div>
    <Hero>
      <template #title>{{ text("home_hero_title", "Мы рядом с первого букета") }}</template>
      <template #lead>
        {{
          text(
            "home_hero_lead",
            "Обучаем современной флористике с нуля — бережно, понятно, с практикой и поддержкой.",
          )
        }}
      </template>
      <template #actions>
        <UiButton variant="primary" to="/#courses">Курсы</UiButton>
        <UiButton variant="outline" to="/#about">О школе</UiButton>
      </template>
      <template v-if="text('home_hero_image')" #media>
        <img
          :src="resolveMediaUrl(text('home_hero_image'))"
          alt=""
          class="aspect-[4/3] w-full rounded-lg object-cover"
        />
      </template>
    </Hero>

    <section class="bg-surface py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">
          Почему стоит учиться в школе “ФлоВей”?
          <template #lead
            >Мы создали школу, в которой удобно учиться, легко развиваться и получать реальные
            практические навыки флористики.</template
          >
        </SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section id="about" class="scroll-mt-64 bg-white py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>О школе</SectionHeading>
        <!-- Бежевая карточка-обёртка вокруг списка (flex-column, не грид); каждый
        пункт внутри — отдельная белая строка на всю ширину. -->
        <div class="rounded-lg bg-surface p-24 sm:p-32">
          <div class="flex flex-col gap-24">
            <div
              v-for="item in aboutItems"
              :key="item.id"
              class="flex flex-col items-start gap-16 rounded-md bg-white p-32"
            >
              <UiBadge>{{ item.badge }}</UiBadge>
              <p class="font-body text-body text-ink">{{ item.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="courses" class="scroll-mt-64 py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">
          Курс “Основы флористики”
          <template #lead>{{ mainCourse?.shortDescription }}</template>
        </SectionHeading>
        <div v-if="mainCourse" class="grid grid-cols-1 gap-24 md:grid-cols-2 lg:gap-32">
          <CourseCard
            v-for="module in mainCourse.modules"
            :key="module.id"
            :title="`Блок “${module.title}”`"
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
          <template #lead
            >После прохождения курса “Основы флористики” вы можете продолжить обучение по одному или
            сразу нескольким специализированным направлениям.</template
          >
        </SectionHeading>
        <div class="grid grid-cols-1 gap-24 md:grid-cols-2 lg:grid-cols-3 lg:gap-32">
          <CourseCard
            v-for="(course, i) in profileCourses"
            :key="course.id"
            :title="course.title"
            :variant="profileVariants[i % profileVariants.length]"
            :cover-image="course.coverImage"
            :to="`/courses/${course.slug}`"
          />
        </div>
      </div>
    </section>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">
          {{ text("home_trial_heading", "Попробуйте флористику на практике") }}
          <template #lead>
            {{
              text(
                "home_trial_description",
                "Вы научитесь собирать круглый букет в спиральной технике: поймёте принцип работы со спиралью, научитесь уверенно удерживать букет в руках во время сборки и правильно подвязывать букет. Кроме практики, вы сможете познакомиться с нашим педагогом, узнать, как проходят занятия в школе, задать все интересующие вопросы и понять, подходит ли вам обучение.",
              )
            }}
          </template>
        </SectionHeading>

        <div class="grid grid-cols-1 gap-32 lg:grid-cols-2 lg:gap-64">
          <div class="flex flex-col gap-24">
            <div class="flex flex-col gap-8">
              <h3 class="text-center font-display text-h2 text-ink">
                {{ text("home_trial_lesson_title", "Пробное занятие") }}
              </h3>
              <p class="text-left font-body text-body text-ink">
                {{ text("home_trial_duration", "Продолжительность: 2,5 часа") }}
              </p>
              <p class="text-left font-body text-body text-ink">
                {{ text("home_trial_price", "Стоимость: 3 000 ₽") }}
              </p>
            </div>
            <ApplyForm context="trial_lesson" variant="simple" title="" />
          </div>
          <!-- TODO: заменить на видео с пробным уроком, когда оно будет готово (пока фото). -->
          <!-- max-w ограничивает контейнер видео, чтобы оно не растягивалось на всю колонку. -->
          <img
            v-if="text('home_trial_image')"
            :src="resolveMediaUrl(text('home_trial_image'))"
            alt=""
            class="mx-auto aspect-[9/16] w-full max-w-[280px] rounded-lg object-cover"
          />
          <div v-else class="mx-auto aspect-[9/16] w-full max-w-[280px] rounded-lg bg-primary" />
        </div>
      </div>
    </section>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>Педагоги</SectionHeading>
        <div class="grid grid-cols-1 gap-24 md:grid-cols-3">
          <div
            v-for="(teacher, index) in teachers"
            :key="teacher.id"
            class="flex flex-col items-center gap-16"
          >
            <img
              v-if="teacher.photo"
              :src="resolveMediaUrl(teacher.photo)"
              :alt="teacher.name"
              class="aspect-square w-full rounded-lg object-cover"
            />
            <div
              v-else
              class="aspect-square w-full rounded-lg"
              :class="index % 2 === 0 ? 'bg-primary' : 'bg-surface'"
            />
            <!-- Тот же стиль (шрифт/размер/жирность), что и у заголовков преимуществ,
            текста кнопок и вопросов FAQ: font-display text-h4 font-bold. -->
            <p
              class="font-display text-h4 font-bold"
              :class="index % 2 === 0 ? 'text-primary' : 'text-ink'"
            >
              {{ teacher.name }}
            </p>
          </div>
        </div>
      </div>
    </section>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container">
        <!-- TODO: дизайна карточек отзывов пока нет (в макете только заголовок) — уточнить у заказчика перед вёрсткой блока. -->
        <SectionHeading>Отзывы</SectionHeading>
      </div>
    </section>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">
          Вопросы и ответы
          <template #lead>Отвечаем на часто задаваемые вопросы.</template>
        </SectionHeading>
        <UiAccordion v-model="openFaqIds">
          <UiAccordionItem v-for="(item, i) in faqItems" :key="i" :id="i" :title="item.question">
            <p>{{ item.answer }}</p>
          </UiAccordionItem>
        </UiAccordion>
      </div>
    </section>
  </div>
</template>
