<script setup lang="ts">
import type { Component } from "vue";
import {
  GiftCertificateMarquee,
  HomeAboutSection,
  HomeCoursesSection,
  HomeFaqSection,
  HomeFeaturesSection,
  HomeGallerySection,
  HomeReviewsSection,
  HomeTeachersSection,
  HomeTrialSection,
} from "#components";
import type { CourseSectionWithCourses, HomeSectionKey, SiteButtonVariant } from "~/types/api";

useSeoMeta({
  title: "Фловей — школа флористики в Москве",
  description:
    "Курсы и мастер-классы по флористике в Москве: с нуля до уверенной самостоятельной работы.",
});

const api = useApi();
const { text } = await usePageContent();
const { glassClass } = await useTreeMode();

// These 9 calls hit independent backend endpoints, so they're fired
// together and awaited via Promise.all instead of one `await` per call —
// awaiting each individually made SSR wait through several sequential
// round-trips before sending any HTML, directly delaying TTFB (and with it
// the hero image's LCP) on the page that gets the most traffic.
const courseSectionsAsync = useAsyncData("home-course-sections", () => api.getCourseSections());
const customDisplayStylesAsync = useAsyncData("home-custom-display-styles", () =>
  api.getCustomDisplayStyles(),
);
const featuresAsync = useAsyncData("home-features", () => api.getFeatures("home"));
const aboutItemsAsync = useAsyncData("home-about-items", () => api.getAboutItems());
const galleryPhotosAsync = useAsyncData("home-gallery-photos", () => api.getGalleryPhotos());
const teachersAsync = useAsyncData("home-teachers", () => api.getTeachers());
const faqAsync = useAsyncData("home-faq", () => api.getFAQItems());
const homeSectionsAsync = useAsyncData("home-sections", () => api.getHomeSections());
const siteButtonsAsync = useAsyncData("site-buttons", () => api.getSiteButtons());

await Promise.all([
  courseSectionsAsync,
  customDisplayStylesAsync,
  featuresAsync,
  aboutItemsAsync,
  galleryPhotosAsync,
  teachersAsync,
  faqAsync,
  homeSectionsAsync,
  siteButtonsAsync,
]);

const { data: courseSectionsData } = courseSectionsAsync;
const courseSections = computed(() => courseSectionsData.value ?? []);

const { data: customDisplayStylesData } = customDisplayStylesAsync;
const customDisplayStylesById = computed(
  () => new Map((customDisplayStylesData.value ?? []).map((s) => [s.id, s])),
);

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
      price: block.price || undefined,
      coverImage: block.blockCover || undefined,
      displayStyle: block.displayStyle,
      customColors: block.customDisplayStyleId
        ? customDisplayStylesById.value.get(block.customDisplayStyleId)
        : undefined,
      to: `/courses/${course.slug}`,
    })),
  );
}

const { data: featuresData } = featuresAsync;
const features = computed(
  () =>
    featuresData.value
      ?.slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((f) => ({ icon: f.icon, title: f.title, description: f.description })) ?? [],
);

const { data: aboutItemsData } = aboutItemsAsync;
const aboutItems = computed(
  () => aboutItemsData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: galleryPhotosData } = galleryPhotosAsync;
const galleryPhotos = computed(
  () => galleryPhotosData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: teachersData } = teachersAsync;
const teachers = computed(
  () => teachersData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: faqData } = faqAsync;
const faqItems = computed(() => faqData.value ?? []);

/**
 * Order/visibility of every homepage block below Hero, admin-managed at
 * /admin/page-content/home-sections. Hero itself is always shown, first,
 * and isn't part of home_sections at all.
 */
const { data: homeSectionsData } = homeSectionsAsync;
const visibleHomeSections = computed(() =>
  (homeSectionsData.value ?? [])
    .filter((section) => section.visible)
    .slice()
    .sort((a, b) => a.sortOrder - b.sortOrder),
);

const HOME_SECTION_COMPONENTS: Record<HomeSectionKey, Component> = {
  features: HomeFeaturesSection,
  gallery: HomeGallerySection,
  courses: HomeCoursesSection,
  gift_certificates: GiftCertificateMarquee,
  trial: HomeTrialSection,
  about: HomeAboutSection,
  teachers: HomeTeachersSection,
  reviews: HomeReviewsSection,
  faq: HomeFaqSection,
};

const sectionProps = computed<Record<HomeSectionKey, Record<string, unknown>>>(() => ({
  features: { text, features: features.value },
  gallery: { photos: galleryPhotos.value },
  courses: { courseSections: courseSections.value, sectionCards },
  gift_certificates: {},
  trial: { text, glassClass: glassClass.value },
  about: { aboutItems: aboutItems.value },
  teachers: { teachers: teachers.value, glassClass: glassClass.value },
  reviews: {},
  faq: { faqItems: faqItems.value },
}));

/**
 * Text/style/link for the site's static standalone CTA buttons — admin-
 * managed at /admin/site-buttons. `fallback` renders if a key is somehow
 * missing (migration not run yet). Fetched in the Promise.all batch above
 * (not via useSiteButtons()) so it doesn't add a sequential round-trip on
 * top of the homepage's other 8 parallel calls.
 */
const { data: siteButtonsData } = siteButtonsAsync;
const siteButtonsByKey = computed(
  () => new Map((siteButtonsData.value ?? []).map((b) => [b.key, b])),
);
function siteButton(
  key: string,
  fallback: { text: string; variant: SiteButtonVariant; url: string },
) {
  const found = siteButtonsByKey.value.get(key);
  return found ? { text: found.text, variant: found.variant, url: found.url } : fallback;
}
const heroCoursesButton = computed(() =>
  siteButton("home_hero_courses", { text: "Курсы", variant: "primary", url: "/#courses" }),
);
const heroTrialButton = computed(() =>
  siteButton("home_hero_trial", { text: "Пробное занятие", variant: "outline", url: "/#trial" }),
);
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
        <UiButton :variant="heroCoursesButton.variant" :to="heroCoursesButton.url">{{
          heroCoursesButton.text
        }}</UiButton>
        <UiButton :variant="heroTrialButton.variant" :to="heroTrialButton.url">{{
          heroTrialButton.text
        }}</UiButton>
      </template>
      <template v-if="text('home_hero_image')" #media>
        <UiHeroPicture :src="resolveOptimizedMediaUrl(text('home_hero_image'))" alt="" />
      </template>
    </Hero>

    <component
      :is="HOME_SECTION_COMPONENTS[section.key]"
      v-for="section in visibleHomeSections"
      :key="section.id"
      v-bind="sectionProps[section.key]"
    />
  </div>
</template>
