import type { Masterclass } from "~/types/api";

/**
 * Mock data used only when NUXT_PUBLIC_USE_MOCKS=true (see useApi.ts) —
 * both GET /api/v1/masterclasses and GET /api/v1/masterclasses/{slug} are
 * real backend routes now. Content below comes directly from the brief,
 * except where noted.
 */
const masterclasses: Masterclass[] = [
  {
    id: 1,
    slug: "krugliy-buket",
    title: "Круглый букет",
    shortDescription:
      "Классический круглый букет, который подойдёт для первого знакомства с флористикой.",
    fullDescription:
      "Соберёте плотный круглый букет по спиральной технике — базовый навык, на котором строится вся остальная флористика.",
    endingText: "Заберёте букет с собой и разберётесь, как ухаживать за ним дома.",
    duration: "2–3 часа",
    priceGroup: 3000,
    priceIndividual: 4500,
    // TODO: сноска про разницу групповой/индивидуальной цены сформулирована
    // по смыслу брифа — согласовать точную формулировку с заказчиком.
    priceDescription:
      "Индивидуальное занятие — если удобнее заниматься один на один с преподавателем, а не в группе.",
    coverImage: "",
    status: "active",
  },
  {
    id: 2,
    slug: "raskidistiy-buket",
    title: "Раскидистый букет",
    shortDescription: "Воздушный асимметричный букет свободной формы.",
    fullDescription:
      "Освоите более сложную, «природную» форму букета — с разной высотой и раскидистым силуэтом вместо плотного купола.",
    endingText: "Заберёте букет с собой и получите рекомендации по уходу.",
    duration: "4 часа",
    priceGroup: 5500,
    coverImage: "",
    status: "active",
  },
  {
    id: 3,
    slug: "kruglaya-kompozitsiya",
    title: "Круглая композиция",
    shortDescription:
      "Компактная композиция во флористической губке — украшение для дома или подарок.",
    fullDescription: "Соберёте аккуратную круглую композицию в кашпо на флористической губке.",
    endingText: "Заберёте готовую композицию с собой.",
    duration: "2–3 часа",
    priceGroup: 3500,
    coverImage: "",
    status: "active",
  },
  {
    id: 4,
    slug: "raskidistaya-kompozitsiya",
    title: "Раскидистая композиция",
    shortDescription: "Крупная композиция свободной формы для интерьера.",
    fullDescription:
      "Соберёте объёмную композицию с раскидистым силуэтом — эффектный акцент для интерьера.",
    endingText: "Заберёте готовую композицию с собой.",
    duration: "3–4 часа",
    priceGroup: 6000,
    coverImage: "",
    status: "active",
  },
  {
    id: 5,
    slug: "interyerny-buket",
    title: "Интерьерный букет",
    shortDescription: "Крупный букет для дома, офиса или в подарок.",
    fullDescription: "Соберёте крупный интерьерный букет — с акцентом на объём и сочетание фактур.",
    endingText: "Заберёте букет с собой.",
    duration: "3–4 часа",
    priceGroup: 6000,
    coverImage: "",
    status: "active",
  },
  {
    id: 6,
    slug: "interyernaya-kompozitsiya",
    title: "Интерьерная композиция",
    shortDescription: "Композиция для интерьера — детали уточняются.",
    fullDescription: "TODO: в брифе текст описания обрезан — уточнить у заказчика.",
    endingText: "",
    // TODO: в макете этот мастер-класс обрезан — нет ни длительности, ни
    // цены. Не выдумываем значения, оставляем незаполненными до уточнения.
    coverImage: "",
    status: "active",
  },
];

export function mockGetMasterClasses(): Masterclass[] {
  return masterclasses;
}

export function mockGetMasterClass(slug: string): Masterclass | null {
  return masterclasses.find((m) => m.slug === slug) ?? null;
}
