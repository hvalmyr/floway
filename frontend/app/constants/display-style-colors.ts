import type { CourseBlockDisplayStyle } from "~/types/api";

/**
 * Shared by CourseCard.vue and MasterclassCard.vue so every colored card on
 * the site (course blocks, masterclasses) draws from the same 4 background/
 * text pairs — one visual grammar, not two components that happen to look
 * similar. The 2 beige-background styles get a border in the same color as
 * their text — otherwise a beige card has no visible edge against the
 * page's own off-white background.
 */
export const displayStyleColorClasses: Record<CourseBlockDisplayStyle, string> = {
  "blue-beige": "bg-primary text-surface",
  "brown-beige": "bg-ink text-surface",
  "beige-blue": "border-2 border-primary bg-surface text-primary",
  "beige-brown": "border-2 border-ink bg-surface text-ink",
};

/** Cycles between the 2 beige-background styles by index — used where
 * there's no per-item admin-chosen style (masterclasses, unlike course
 * blocks). Beige for both keeps every masterclass row visually consistent;
 * only the accent (blue vs. brown) alternates. */
export const masterclassDisplayStyleCycle: CourseBlockDisplayStyle[] = [
  "beige-blue",
  "beige-brown",
];
