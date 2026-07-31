import type { Config } from "tailwindcss";

/**
 * Design-system theme for «ФлоВей». Values come from docs/floway-design.md.
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
      colors: {
        primary: {
          50: "#EAF2F7",
          300: "#A9C9DD",
          500: "#82B1CC",
          600: "#5F93B3",
          700: "#42708F",
        },
        ink: {
          900: "#41342A",
          700: "#6B5D51",
          400: "#A79A8E",
        },
        surface: {
          DEFAULT: "#F7F5F3",
          100: "#FBFAF9",
        },
        line: "#E4DFDA",
        success: "#5B9A6F",
        error: "#C15B4A",
      },
      fontFamily: {
        display: ["Nunito", "sans-serif"],
        body: ["Onest", "sans-serif"],
      },
      fontSize: {
        h1: ["var(--text-h1)", { lineHeight: "1.08", fontWeight: "800" }],
        h2: ["var(--text-h2)", { lineHeight: "1.15", fontWeight: "800" }],
        h3: ["var(--text-h3)", { lineHeight: "1.25", fontWeight: "700" }],
        h4: ["var(--text-h4)", { lineHeight: "1.3", fontWeight: "700" }],
        "body-l": ["var(--text-body-l)", { lineHeight: "1.5" }],
        body: ["var(--text-body)", { lineHeight: "1.6" }],
        small: ["var(--text-small)", { lineHeight: "1.4", fontWeight: "600" }],
        button: ["var(--text-button)", { lineHeight: "1" }],
      },
    },
  },
  plugins: [],
} satisfies Config;
