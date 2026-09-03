<script setup lang="ts">
import type { Component } from "vue";
import { featureIconComponent } from "~/constants/feature-icons";
import type { CourseSectionWithCourses } from "~/types/api";

useSeoMeta({
  title: "Фловей — школа флористики в Москве",
  description:
    "Курсы и мастер-классы по флористике в Москве: с нуля до уверенной самостоятельной работы.",
});

const api = useApi();
const { text } = await usePageContent();

const { data: courseSectionsData } = await useAsyncData("home-course-sections", () =>
  api.getCourseSections(),
);
const courseSections = computed(() => courseSectionsData.value ?? []);

/**
 * One card per VISIBLE block — a course with a single block (or none at
 * all, in which case the backend hands back one synthetic block built from
 * the course's own fields, see model.Course's Go doc comment) renders as
 * one card; a course with several blocks (e.g. "Основы флористики" with a
 * "Букеты" and a "Композиции" block) renders one card per block, each
 * titled with the course's name. `block.blockName` is blank for the
 * synthetic case, so `blockLabel` naturally comes out empty there and
 * CourseCard just shows lessonCount/timeLength — no separate branch needed
 * for "course with no blocks" vs. "course with exactly one named block".
 */
function sectionCards(section: CourseSectionWithCourses) {
  return section.courses.flatMap((course) =>
    course.blocks.map((block, index) => ({
      key: `${course.id}-${index}`,
      name: course.name,
      blockLabel: block.blockName || undefined,
      lessonCount: block.lessonCount || undefined,
      timeLength: block.timeLength || undefined,
      coverImage: block.blockCover || undefined,
      displayStyle: block.displayStyle,
      to: `/courses/${course.slug}`,
    })),
  );
}

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

const { data: galleryPhotosData } = await useAsyncData("home-gallery-photos", () =>
  api.getGalleryPhotos(),
);
const galleryPhotos = computed(
  () => galleryPhotosData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: teachersData } = await useAsyncData("home-teachers", () => api.getTeachers());
const teachers = computed(
  () => teachersData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: faqData } = await useAsyncData("home-faq", () => api.getFAQItems());
const faqItems = computed(() => faqData.value ?? []);
const openFaqIds = ref<Array<string | number>>([0]);

function capitalizeName(name: string): string {
  return name
    .split(" ")
    .map((part) => (part ? part[0]!.toUpperCase() + part.slice(1) : part))
    .join(" ");
}
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
        <UiButton variant="outline" to="/#trial">Пробное занятие</UiButton>
      </template>
      <template v-if="text('home_hero_image')" #media>
        <UiHeroPicture :src="resolveOptimizedMediaUrl(text('home_hero_image'))" alt="" />
      </template>
    </Hero>

    <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>
          {{ text("home_features_heading", "Почему стоит учиться в школе «Фловей»?") }}
          <template #lead>
            {{
              text(
                "home_features_lead",
                "Мы создали школу, в которой удобно учиться, легко развиваться и получать реальные практические навыки флористики.",
              )
            }}
          </template>
        </SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section v-if="galleryPhotos.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container">
        <PhotoCarousel :photos="galleryPhotos" />
      </div>
    </section>

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

    <section
      v-for="(section, sIndex) in courseSections"
      :id="sIndex === 0 ? 'courses' : undefined"
      :key="section.id"
      class="py-48 sm:py-64 lg:py-80"
      :class="sIndex === 0 ? 'scroll-mt-64 lg:scroll-mt-96' : ''"
    >
      <div class="container flex flex-col gap-48">
        <SectionHeading :color="sIndex % 2 === 0 ? 'primary' : 'ink'">
          {{ section.heading }}
          <template #lead>{{ section.description }}</template>
        </SectionHeading>
        <div class="flex flex-wrap justify-center gap-24 lg:gap-32">
          <CourseCard
            v-for="card in sectionCards(section)"
            :key="card.key"
            :name="card.name"
            :display-style="card.displayStyle"
            :block-label="card.blockLabel"
            :lesson-count="card.lessonCount"
            :time-length="card.timeLength"
            :cover-image="card.coverImage"
            :to="card.to"
            class="w-full sm:w-[calc(50%-12px)] lg:w-[calc(33.333%-22px)]"
          />
        </div>
      </div>
    </section>

    <GiftCertificateMarquee />

    <section id="trial" class="scroll-mt-64 py-48 sm:py-64 lg:scroll-mt-96 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">
          {{ text("trial_section_heading", "Попробуйте флористику на практике") }}
          <template #lead>
            {{
              text(
                "trial_section_description",
                "Вы научитесь собирать круглый букет в спиральной технике: поймёте принцип работы со спиралью, научитесь уверенно удерживать букет в руках во время сборки и правильно подвязывать букет. Кроме практики, вы сможете познакомиться с нашим педагогом, узнать, как проходят занятия в школе, задать все интересующие вопросы и понять, подходит ли вам обучение.",
              )
            }}
          </template>
        </SectionHeading>

        <div class="grid grid-cols-1 gap-32 md:grid-cols-2 md:gap-64">
          <div class="order-2 flex w-full flex-col gap-24 md:order-1">
            <div
              class="flex flex-col gap-16 rounded-md bg-white/55 px-16 py-24 backdrop-blur backdrop-saturate-150"
            >
              <h3 class="font-body text-h4 text-ink">
                {{ text("trial_heading", "Пробное занятие") }}
              </h3>
              <!-- Single \n (not \n\n) so duration+price render as one tight
              paragraph with a <br> between them, not two separately-spaced
              ones — otherwise the gap between them (markdown paragraph
              margin) was bigger than the gap to the heading above, so
              duration read as grouped with "Пробное занятие" instead of
              with the price right under it. -->
              <MarkdownContent
                :source="
                  text('trial_description', 'Продолжительность: 2,5 часа.\nСтоимость: 3 000 ₽.')
                "
              />
            </div>
            <ApplyForm context="trial_lesson" title="" bare class="w-full" />
          </div>
          <!-- TODO: заменить на видео с пробным уроком, когда оно будет готово (пока фото). -->
          <!-- Портретное 9:16, растянуто до ширины колонки — но не выше 80%
          экрана: max-h ограничивает высоту, а aspect-ratio при этом сжимает
          и ширину пропорционально, так что соотношение сторон не ломается
          (стандартное поведение aspect-ratio + max-height у браузера, без
          JS). mx-auto центрирует, когда из-за max-h ширина не дотягивает до
          полной колонки. -->
          <NuxtImg
            v-if="text('home_trial_image')"
            :src="resolveOptimizedMediaUrl(text('home_trial_image'))"
            format="webp"
            alt=""
            class="order-1 mx-auto aspect-[9/16] max-h-[80vh] max-w-full rounded-lg object-cover md:sticky md:top-96 md:order-2"
            sizes="400:100vw md:50vw"
            loading="lazy"
          />
          <div
            v-else
            class="order-1 mx-auto aspect-[9/16] max-h-[80vh] max-w-full rounded-lg bg-primary md:sticky md:top-96 md:order-2"
          />
        </div>
      </div>
    </section>

    <section class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading>Педагоги</SectionHeading>
        <div class="flex flex-wrap justify-center gap-24">
          <div
            v-for="(teacher, index) in teachers"
            :key="teacher.id"
            class="flex w-full flex-col items-center gap-16 md:w-[calc((100%-48px)/3)]"
          >
            <NuxtImg
              v-if="teacher.photo"
              :src="resolveOptimizedMediaUrl(teacher.photo)"
              format="webp"
              :alt="teacher.name"
              class="aspect-square w-full rounded-lg object-cover"
              sizes="400:100vw md:33vw"
              loading="lazy"
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
              class="w-full rounded-md bg-white/55 py-12 text-center font-body text-h4 font-medium backdrop-blur backdrop-saturate-150"
              :class="index % 2 === 0 ? 'text-primary' : 'text-ink'"
            >
              {{ capitalizeName(teacher.name) }}
            </p>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>
          Отзывы
          <template #lead
            >Смотрите отзывы о школе на Яндекс Картах — заходите оставить свой.</template
          >
        </SectionHeading>
        <!-- Официальный виджет отзывов Яндекс Карт — сам виджет задаёт
        style width:100%/height:100% (заказчик прислал), поэтому обёртка
        ниже даёт ему конкретную высоту, внутри которой он растягивается. -->
        <div class="mx-auto h-[950px] w-full max-w-[400px]">
          <iframe
            src="https://yandex.ru/maps-reviews-widget/83657275642?comments"
            title="Отзывы о школе «Фловей» на Яндекс Картах"
            loading="lazy"
            style="
              width: 100%;
              height: 100%;
              border: 1px solid #e6e6e6;
              border-radius: 8px;
              box-sizing: border-box;
            "
          />
        </div>
      </div>
    </section>

    <section class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">
          Вопросы и ответы
          <template #lead>Отвечаем на часто задаваемые вопросы.</template>
        </SectionHeading>
        <UiAccordion v-model="openFaqIds">
          <UiAccordionItem v-for="(item, i) in faqItems" :key="i" :id="i" :title="item.question">
            <MarkdownContent :source="item.answer" />
          </UiAccordionItem>
        </UiAccordion>
      </div>
    </section>
  </div>
</template>
