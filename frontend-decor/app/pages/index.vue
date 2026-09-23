<script setup lang="ts">
useSeoMeta({
  title: "flo-way — новогоднее оформление домов, офисов и мероприятий",
  description:
    "Оформляем дома, офисы, магазины, рестораны и отели к Новому году в Москве и Московской области — на сезон или на один праздник.",
});

const api = useApi();
const { text } = await usePageContent();

const { data: landingPagesData } = await useAsyncData("home-landing-pages", () =>
  api.getLandingPages(),
);
const landingPages = computed(
  () => landingPagesData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: aboutItemsData } = await useAsyncData("home-about-items", () => api.getAboutItems());
const aboutItems = computed(
  () => aboutItemsData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: featuresData } = await useAsyncData("home-features", () => api.getFeatures("home"));
const features = computed(
  () =>
    featuresData.value
      ?.slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((f) => ({ icon: f.icon, title: f.title, description: f.description })) ?? [],
);

const { data: galleryPhotosData } = await useAsyncData("home-gallery-photos", () =>
  api.getGalleryPhotos(),
);
const galleryPhotos = computed(
  () => galleryPhotosData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: faqData } = await useAsyncData("home-faq", () => api.getFAQItems());
const faqItems = computed(() => faqData.value ?? []);
const openFaqIds = ref<Array<string | number>>([0]);
</script>

<template>
  <div>
    <Hero>
      <template #title>{{
        text("home_hero_title", "Новогоднее оформление домов, офисов и мероприятий")
      }}</template>
      <template #lead>
        {{
          text(
            "home_hero_lead",
            "Оформляем на сезон (декабрь-январь) или на один праздник — рассчитываем стоимость индивидуально после бесплатного выезда и замеров.",
          )
        }}
      </template>
      <template #actions>
        <UiButton variant="primary" to="#apply">Оставить заявку</UiButton>
      </template>
      <template v-if="text('home_hero_image')" #media>
        <UiHeroPicture :src="resolveOptimizedMediaUrl(text('home_hero_image'))" alt="" />
      </template>
    </Hero>

    <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>
          Два формата оформления
          <template #lead>Выбирайте то, что подходит именно вам.</template>
        </SectionHeading>
        <div class="flex flex-wrap justify-center gap-24 lg:gap-32">
          <div
            class="flex w-full flex-col items-center gap-16 rounded-md bg-white p-32 text-center md:w-[calc((100%_-_24px)/2)] lg:w-[calc((100%_-_64px)/3)]"
          >
            <IconFirBranch class="h-40 w-auto text-primary" />
            <h3 class="font-display text-h4 text-primary">На сезон</h3>
            <p class="font-body text-body text-ink">
              Оформляем объект на весь новогодний сезон — с декабря по январь. Подходит для домов,
              офисов, магазинов и ресторанов.
            </p>
          </div>
          <div
            class="flex w-full flex-col items-center gap-16 rounded-md bg-white p-32 text-center md:w-[calc((100%_-_24px)/2)] lg:w-[calc((100%_-_64px)/3)]"
          >
            <IconSnowflake class="h-40 w-auto text-primary" />
            <h3 class="font-display text-h4 text-primary">На праздник</h3>
            <p class="font-body text-body text-ink">
              Разовое оформление под конкретное мероприятие — корпоратив, зимний юбилей, открытие —
              с демонтажом после праздника.
            </p>
          </div>
        </div>
      </div>
    </section>

    <section class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="ink">
          Для кого мы оформляем
          <template #lead
            >Частные дома, офисы, магазины, рестораны и отели в Москве и Московской области.
            Торговые центры не оформляем — слишком крупные объекты.</template
          >
        </SectionHeading>
        <div class="rounded-lg bg-surface/55 p-24 backdrop-blur backdrop-saturate-150 sm:p-32">
          <div class="flex flex-col items-start gap-24">
            <div
              v-for="item in aboutItems"
              :key="item.id"
              class="flex w-full flex-col items-start gap-16 rounded-md bg-white p-32"
            >
              <UiBadge>{{ item.badge }}</UiBadge>
              <p class="whitespace-pre-line font-body text-body text-ink">{{ item.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section
      v-if="landingPages.length"
      class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80"
    >
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>Что мы оформляем</SectionHeading>
        <div class="flex flex-wrap justify-center gap-16">
          <UiButton
            v-for="page in landingPages"
            :key="page.id"
            variant="outline"
            :to="`/${page.slug}`"
          >
            {{ page.objectType.name }}
          </UiButton>
        </div>
      </div>
    </section>

    <section class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="ink">
          Как мы работаем
          <template #lead>От заявки до демонтажа — прозрачно и без сюрпризов.</template>
        </SectionHeading>
        <div class="flex flex-col gap-16">
          <UiInfoRow
            v-for="(step, i) in [
              'Заявка на сайте или в мессенджере',
              'Бесплатный выезд на объект и замеры',
              'Визуал будущего оформления',
              'Коммерческое предложение — презентация',
              'Согласование деталей',
              'Монтаж',
              'Демонтаж после сезона или праздника',
            ]"
            :key="step"
            :items="[`${i + 1}. ${step}`]"
            :highlighted="i % 2 === 1"
          />
        </div>
        <UiButton variant="outline" to="/how-it-works" class="mx-auto"
          >Подробнее об этапах</UiButton
        >
      </div>
    </section>

    <section v-if="galleryPhotos.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-32">
        <SectionHeading color="ink">Примеры работ</SectionHeading>
        <PhotoCarousel :photos="galleryPhotos" />
        <UiButton variant="outline" to="/portfolio" class="mx-auto">Всё портфолио</UiButton>
      </div>
    </section>

    <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>Почему выбирают нас</SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-24 text-center">
        <SectionHeading color="ink">Отзывы</SectionHeading>
        <p class="font-body text-body text-ink">
          {{
            text(
              "home_reviews_placeholder",
              "Здесь появятся отзывы наших клиентов. Если вы уже работали с нами — напишите нам в мессенджер, будем рады отзыву.",
            )
          }}
        </p>
      </div>
    </section>

    <section
      id="apply"
      class="scroll-mt-64 bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:scroll-mt-96 lg:py-80"
    >
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm title="Оставить заявку" />
        </div>
      </div>
    </section>

    <section v-if="faqItems.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="accent">
          Частые вопросы
          <template #lead>Отвечаем на самые популярные вопросы.</template>
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
