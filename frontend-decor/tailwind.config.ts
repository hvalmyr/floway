import type { Config } from "tailwindcss";

/**
 * Design-system theme for «Фловей». Values come from docs/floway-design.md.
 *
 * `spacing`/`borderRadius` go through `extend` (not a full `theme` override)
 * on purpose: the pre-existing admin panel already uses Tailwind's default
 * scale (`py-2`, `gap-3`, `mt-6`, bare `rounded`, ...). Extending only adds/
 * overrides the specific keys the design system needs (4/8/12/16/24/32/40/
 * 48/64/80/96/120/160 for spacing; sm/md/lg/pill for radius) and leaves every
 * other default key intact, so admin markup keeps rendering as before.
 */
export default {
  content: [
    "./app/components/**/*.{vue,js,ts}",
    "./app/layouts/**/*.vue",
    "./app/pages/**/*.vue",
    "./app/plugins/**/*.{js,ts}",
    "./app/composables/**/*.{js,ts}",
    "./app/constants/**/*.{js,ts}",
    "./app/lib/**/*.{js,ts}",
    "./app/middleware/**/*.{js,ts}",
    "./app/utils/**/*.{js,ts}",
    "./app/app.vue",
    "./app/error.vue",
  ],
  theme: {
    screens: {
      sm: "480px",
      md: "768px",
      lg: "1024px",
      xl: "1280px",
    },
    container: {
      center: true,
      screens: { xl: "1280px" },
      padding: {
        DEFAULT: "16px",
        sm: "24px",
        md: "40px",
        lg: "64px",
        xl: "64px",
      },
    },
    extend: {
      spacing: {
        4: "4px",
        8: "8px",
        12: "12px",
        16: "16px",
        24: "24px",
        32: "32px",
        40: "40px",
        48: "48px",
        64: "64px",
        80: "80px",
        96: "96px",
        120: "120px",
        160: "160px",
      },
      borderRadius: {
        sm: "12px",
        md: "20px",
        lg: "28px",
        pill: "999px",
      },
      // Ровно 4 цвета по требованию заказчика — никаких оттенков/шейдов.
      // Состояния валидации форм используют primary, не отдельный цвет.
      colors: {
        primary: "#82B1CC",
        surface: "#F7F5F3",
        ink: "#41342A",
      },
      fontFamily: {
        display: ["Soyuz Grotesk", "sans-serif"],
        body: ["Non Bureau", "sans-serif"],
      },
      // Exactly 4 sizes by client direction: hero-block headings (h1),
      // section headings (h2), everything else that reads as "a small
      // heading" — item headings, buttons, cards, badges, course cards,
      // header nav links (h4/button, already the same size) — and body copy
      // for descriptions and the form (body). h3/body-l/small were dropped;
      // anything that used them moved to the nearest of these four.
      fontSize: {
        h1: ["var(--text-h1)", { lineHeight: "1.08", fontWeight: "800" }],
        h2: ["var(--text-h2)", { lineHeight: "1.15", fontWeight: "800" }],
        h4: ["var(--text-h4)", { lineHeight: "1.3", fontWeight: "700" }],
        body: ["var(--text-body)", { lineHeight: "1.6" }],
        // Одинаковый шрифт/размер/жирность с заголовками преимуществ, именами
        // педагогов и вопросами FAQ — везде text-h4 font-display font-bold, у
        // кнопок только line-height туже (однострочный текст).
        button: ["var(--text-h4)", { lineHeight: "1", fontWeight: "700" }],
      },
    },
  },
  plugins: [],
} satisfies Config;
