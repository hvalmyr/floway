import type { Course, CourseDetail, CourseModule } from "~/types/api";

/**
 * Mock data standing in for the Go backend until it exposes a public
 * "course with nested blocks/lessons" endpoint (see the TODO on
 * `CourseDetail` in ~/types/api.ts). Shaped identically to what that
 * endpoint should return, so swapping useApi() over to the real API later
 * doesn't require touching any component.
 *
 * Titles/intro copy below come straight from the course mockup
 * (desktop-course.png) where visible; per-lesson "Темы"/"Вы научитесь" text
 * is only confirmed for "Занятие 1" of each block — the rest is placeholder,
 * clearly TODO-marked, pending the real program from the school's methodists.
 */

function buketyLessons(): CourseModule["lessons"] {
  // Заголовки 5–7 в макете буквально повторяются («Раскидистый букет» трижды)
  // — воспроизведено как есть, не исправлено произвольно.
  const titles = [
    "Спиральная техника",
    "Пропорции в букете",
    "Флористические стили",
    "Колористика во флористике",
    "Раскидистый букет",
    "Раскидистый букет",
    "Раскидистый букет",
  ];
  return titles.map((title, i) => ({
    id: i + 1,
    courseBlockId: 1,
    number: i + 1,
    title: `Занятие ${i + 1}. ${title}`,
    topics:
      i === 0
        ? "Спиральная техника сборки букета, простые способы расстановки цветов, флористические инструменты и аксессуары."
        : "TODO: согласовать точный список тем занятия с методистом школы.",
    outcomes:
      i === 0
        ? "Собирать свой первый круглый букет в спиральной технике, правильно подвязывать букет, познакомитесь с основными флористическими инструментами и аксессуарами, а также начнёте отрабатывать постановку руки для уверенной сборки букетов."
        : "TODO: согласовать формулировку результата занятия с методистом школы.",
    durationHours: 4,
  }));
}

function kompozitsiiLessons(): CourseModule["lessons"] {
  const titles = ["Основы композиции в вазе", "Композиции на флористической губке", "Интерьерная композиция", "Раскидистая композиция"];
  return titles.map((title, i) => ({
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
  shortDescription:
    "Курс подойдет тем, кто хочет освоить флористику с нуля и получить прочную базу знаний и практических навыков. Вы научитесь создавать современные букеты и композиции, освоите основные техники и принципы работы с цветами.",
  fullDescription:
    "Курс подойдет тем, кто хочет освоить флористику с нуля и получить прочную базу знаний и практических навыков. Вы научитесь создавать современные букеты и композиции, освоите основные техники и принципы работы с цветами.",
  status: "active",
  coverImage: "",
  gallery: [],
  sortOrder: 1,
  modules: [
    {
      id: 1,
      courseId: 1,
      title: "Букеты",
      description:
        "Блок посвящён основным техникам сборки современных букетов. Вы научитесь работать со спиральной техникой, создавать разные формы букетов, правильно подбирать цветы и зелень, а также добиваться гармоничной формы и пропорций.",
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
      // TODO: согласовать с методистом реальный вводный текст блока «Композиции» (в макете не показан).
      description: "Блок посвящён технике сборки цветочных композиций в вазе и на флористической губке.",
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
// компонент карточки должен рендерить «?» для отсутствующих числовых полей,
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
