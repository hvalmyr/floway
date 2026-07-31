import type { Course, CourseDetail, CourseModule } from "~/types/api";

/**
 * Mock data standing in for the Go backend until it exposes a public
 * "course with nested blocks/lessons" endpoint (see the TODO on
 * `CourseDetail` in ~/types/api.ts). Shaped identically to what that
 * endpoint should return, so swapping useApi() over to the real API later
 * doesn't require touching any component.
 *
 * Per-lesson "Темы"/"Вы научитесь" copy below is placeholder text for
 * layout purposes only — TODO: заменить на согласованный учебный план от
 * методистов школы перед запуском.
 */

function buketyLessons(): CourseModule["lessons"] {
  const topics = [
    "Введение в профессию флориста, материалы и инструменты",
    "Основы композиции: форма, ритм, баланс",
    "Круглый букет: техника спирали",
    "Каскадный и асимметричный букет",
    "Работа с сезонными и полевыми цветами",
    "Упаковка и оформление букета",
    "Итоговая работа и разбор композиций",
  ];
  return topics.map((title, i) => ({
    id: i + 1,
    courseBlockId: 1,
    number: i + 1,
    title: `Занятие ${i + 1}. ${title}`,
    topics: "TODO: согласовать точный список тем занятия с методистом школы.",
    outcomes: "TODO: согласовать формулировку результата занятия с методистом школы.",
    durationHours: 4,
  }));
}

function kompozitsiiLessons(): CourseModule["lessons"] {
  const topics = [
    "Основы флористической композиции в вазе",
    "Композиции на флористической губке",
    "Интерьерная композиция для дома и офиса",
    "Итоговая композиция и разбор работ",
  ];
  return topics.map((title, i) => ({
    id: 100 + i + 1,
    courseBlockId: 2,
    number: i + 1,
    title: `Занятие ${i + 1}. ${title}`,
    topics: "TODO: согласовать точный список тем занятия с методистом школы.",
    outcomes: "TODO: согласовать формулировку результата занятия с методистом школы.",
    durationHours: 4,
  }));
}

const osnovyFloristikiDetail: CourseDetail = {
  id: 1,
  slug: "osnovy-floristiki",
  title: "Основы флористики",
  shortDescription: "Обучаем современной флористике с нуля — бережно, понятно, с практикой и поддержкой.",
  fullDescription:
    "Курс для тех, кто хочет с нуля освоить флористику: от первого букета до уверенной самостоятельной работы с цветом, формой и композицией. Обучение построено на практике — большую часть времени вы проведёте с живыми цветами, а не с теорией.",
  status: "active",
  coverImage: "",
  gallery: [],
  sortOrder: 1,
  modules: [
    {
      id: 1,
      courseId: 1,
      title: "Букеты",
      lessonsCount: 7,
      hours: 30,
      price: 38500,
      sortOrder: 1,
      lessons: buketyLessons(),
    },
    {
      id: 2,
      courseId: 1,
      title: "Композиции",
      lessonsCount: 4,
      // TODO: бриф расходится сам с собой — на главной указано 16 часов для
      // блока «Композиции», а в блоке цены на странице курса — 17 часов.
      // Уточнить у заказчика правильное значение; сейчас взято 17ч (из более
      // детального описания блока цены).
      hours: 17,
      price: 22000,
      sortOrder: 2,
      lessons: kompozitsiiLessons(),
    },
  ],
};

// Профильные курсы: в брифе для них ещё не заполнены кол-во занятий/часов —
// компонент карточки должен рендерить «—» для отсутствующих числовых полей,
// а не выдумывать значения.
const profileCourses: CourseDetail[] = [
  {
    id: 2,
    slug: "kommercheskaya-floristika",
    title: "Коммерческая флористика",
    shortDescription: "",
    fullDescription: "TODO: описание курса ожидает контента от заказчика.",
    status: "active",
    coverImage: "",
    gallery: [],
    sortOrder: 2,
    modules: [],
  },
  {
    id: 3,
    slug: "aktualnaya-floristika",
    title: "Актуальная флористика",
    shortDescription: "",
    fullDescription: "TODO: описание курса ожидает контента от заказчика.",
    status: "active",
    coverImage: "",
    gallery: [],
    sortOrder: 3,
    modules: [],
  },
  {
    id: 4,
    slug: "svadebnaya-floristika",
    title: "Свадебная флористика",
    shortDescription: "",
    fullDescription: "TODO: описание курса ожидает контента от заказчика.",
    status: "active",
    coverImage: "",
    gallery: [],
    sortOrder: 4,
    modules: [],
  },
];

const allCourses: CourseDetail[] = [osnovyFloristikiDetail, ...profileCourses];

export function mockGetCourses(): Course[] {
  return allCourses.map(({ modules: _modules, ...course }) => course);
}

export function mockGetCourse(slug: string): CourseDetail | null {
  return allCourses.find((c) => c.slug === slug) ?? null;
}
