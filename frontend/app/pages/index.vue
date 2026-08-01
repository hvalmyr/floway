<script setup lang="ts">
import IconEightyTwenty from "~/components/ui/IconEightyTwenty.vue";
import IconFlexStart from "~/components/ui/IconFlexStart.vue";
import IconLevels from "~/components/ui/IconLevels.vue";
import IconMapPin from "~/components/ui/IconMapPin.vue";
import IconPeopleGroup from "~/components/ui/IconPeopleGroup.vue";
import IconTwoFormats from "~/components/ui/IconTwoFormats.vue";

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

const features = [
  {
    icon: IconTwoFormats,
    title: "два формата обучения",
    description:
      "Обучение в группе по расписанию или свободное посещение — выбирайте формат, который подходит именно вам.",
  },
  {
    icon: IconPeopleGroup,
    title: "комфортное обучение",
    description: "Небольшие группы из 5–7 человек позволяют преподавателю уделить внимание каждому ученику.",
  },
  {
    icon: IconMapPin,
    title: "центр москвы",
    description:
      "Школа находится в центре Москвы, в пешей доступности от станций метро Арбатская, Смоленская, Баррикадная и Краснопресненская.",
  },
  {
    icon: IconFlexStart,
    title: "гибкий старт",
    description: "Не нужно ждать набора новой группы — начните обучение, когда будете готовы.",
  },
  {
    icon: IconEightyTwenty,
    title: "максимум практики",
    description: "Большую часть обучения занимает работа с живыми цветами. Только 20% курса — теория.",
  },
  {
    icon: IconLevels,
    title: "для любого уровня",
    description: "Обучаем как новичков, так и практикующих флористов.",
  },
];

const aboutItems = [
  {
    badge: "есть свободный формат",
    description:
      "Формат свободного посещения позволяет начать обучение без ожидания набора группы. Вы сами выбираете удобные даты и время занятий. Каждое занятие начинается с индивидуального объяснения новой темы, после чего основное время посвящено практике под постоянным сопровождением преподавателя.",
  },
  {
    badge: "есть групповой формат",
    description:
      "Мы проводим обучение в небольших группах. Ближайшие даты набора можно посмотреть в календаре событий. При раннем бронировании курса действует скидка 10%.",
  },
  {
    badge: "маленькие группы до 7 человек",
    description:
      "Мы принципиально не работаем с большими группами. На групповых занятиях обучается до 5–7 человек, а в формате свободного посещения преподаватель одновременно сопровождает не более 5–7 учеников. Такой подход позволяет уделить внимание каждому, ответить на вопросы и помочь закрепить материал на практике.",
  },
  {
    badge: "только живые цветы",
    description:
      "Мы убеждены, что работать со спиральной техникой важно сразу на живых цветах. Каждый стебель отличается толщиной, гибкостью, формой и особенностями строения, поэтому именно живой материал учит чувствовать правильное положение цветов в руке и понимать, как формируется спираль. Спиральная техника — это не только положение стеблей, но и умение работать с живым материалом. Поэтому её невозможно полноценно освоить на материале, который не передаёт естественные особенности растений. Чем раньше появляется этот навык, тем увереннее флорист чувствует себя в дальнейшей работе с цветами.",
  },
  {
    badge: "работаем с 2013 года",
    description:
      "Школа ФлоВей работает с 2013 года. За это время мы обучили флористике сотни учеников — от новичков до тех, кто решил связать с цветами свою профессию.",
  },
];

const teachers = [
  { name: "Алёна Рыбкина", accent: true },
  { name: "Елена Зайцева", accent: false },
  { name: "Анна Лобус", accent: true },
];

// TODO: реальные ответы согласованы только для первого вопроса (текст с
// макета); остальные — по смыслу вопроса, требуют подтверждения школы перед
// публикацией.
const faqItems = [
  {
    question: "Как проходит обучение?",
    answer:
      "Обучение состоит из теоретической и практической частей, при этом 80% времени занимают занятия с живыми цветами. Все курсы доступны как в группе, так и в формате свободного посещения.",
  },
  {
    question: "Что такое свободное посещение?",
    answer: "Формат, при котором вы сами выбираете даты и время занятий, не дожидаясь набора группы.",
  },
  {
    question: "Проводите ли вы мастер-классы?",
    answer: "Да — разовые мастер-классы по букетам и композициям, подробности на странице «Мастер-классы».",
  },
  {
    question: "Как записаться на обучение?",
    answer: "Оставьте заявку в форме на сайте — мы свяжемся с вами и поможем выбрать курс или мастер-класс.",
  },
  {
    question: "Нужно ли покупать цветы и инструменты?",
    answer: "Нет, все материалы и инструменты предоставляет школа — приходите налегке.",
  },
  {
    question: "Подойдет ли обучение, если у меня совсем нет опыта?",
    answer: "Да, мы обучаем с нуля — программа рассчитана в том числе на новичков.",
  },
];
const openFaqIds = ref<Array<string | number>>([0]);
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
        <SectionHeading color="primary">
          Почему стоит учиться в школе «ФлоВей»?
          <template #lead>Мы создали школу, в которой удобно учиться, легко развиваться и получать реальные практические навыки флористики.</template>
        </SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section id="about" class="scroll-mt-64 bg-surface py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>О школе</SectionHeading>
        <!-- Список (flex-column), не грид: каждый пункт — отдельная строка на всю ширину. -->
        <div class="flex flex-col gap-24">
          <div v-for="item in aboutItems" :key="item.badge" class="flex flex-col items-start gap-16 rounded-md bg-white p-32">
            <UiBadge>{{ item.badge }}</UiBadge>
            <p class="font-body text-body text-ink">{{ item.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <section id="courses" class="scroll-mt-64 py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">
          Курс «Основы флористики»
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
          <template #lead>После прохождения курса «Основы флористики» вы можете продолжить обучение по одному или сразу нескольким специализированным направлениям.</template>
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
      <div class="container flex flex-col gap-48">
        <div class="flex flex-col gap-16">
          <h2 class="font-display text-h2 text-primary">Попробуйте флористику на практике</h2>
          <p class="max-w-[760px] font-body text-body-l text-ink">
            Вы научитесь собирать круглый букет в спиральной технике: поймёте принцип работы со спиралью, научитесь
            уверенно удерживать букет в руках во время сборки и правильно подвязывать букет. Кроме практики, вы
            сможете познакомиться с нашим педагогом, узнать, как проходят занятия в школе, задать все интересующие
            вопросы и понять, подходит ли вам обучение.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-32 lg:grid-cols-2 lg:gap-64">
          <div class="flex flex-col gap-24">
            <div class="flex flex-col gap-8">
              <h3 class="font-display text-h3 text-ink">Пробное занятие</h3>
              <p class="font-body text-body text-ink">Продолжительность: 2,5 часа</p>
              <p class="font-body text-body text-ink">Стоимость: 3 000 ₽</p>
            </div>
            <ApplyForm context="trial_lesson" variant="simple" title="" />
          </div>
          <!-- TODO: заменить на видео с пробным уроком, когда оно будет готово. -->
          <div class="aspect-[4/3] rounded-lg bg-primary lg:aspect-auto" />
        </div>
      </div>
    </section>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading>Педагоги</SectionHeading>
        <div class="grid grid-cols-1 gap-24 md:grid-cols-3">
          <div v-for="teacher in teachers" :key="teacher.name" class="flex flex-col items-center gap-16">
            <div class="aspect-square w-full rounded-lg" :class="teacher.accent ? 'bg-primary' : 'bg-surface'" />
            <p class="font-display text-h4" :class="teacher.accent ? 'text-primary' : 'text-ink'">
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
