import type { SiteButton, SiteButtonVariant } from "~/types/api";

/**
 * Looks up a static, standalone CTA button's admin-edited text/style/link
 * by key (see SiteButton in types/api.ts) — same pattern as
 * usePageContent()'s `text()`. `fallback` is what renders if the key is
 * somehow missing (migration not run yet, typo) — pass the current
 * hardcoded copy/variant/url as a safety net rather than a dead button.
 *
 * @example
 * const { button } = await useSiteButtons();
 * // in template:
 * // <UiButton :variant="button('home_hero_courses', {text: 'Курсы', variant: 'primary', url: '/#courses'}).variant"
 * //           :to="button('home_hero_courses', {text: 'Курсы', variant: 'primary', url: '/#courses'}).url">
 * //   {{ button('home_hero_courses', {text: 'Курсы', variant: 'primary', url: '/#courses'}).text }}
 * // </UiButton>
 */
export async function useSiteButtons() {
  const api = useApi();
  const { data } = await useAsyncData("site-buttons", () => api.getSiteButtons());

  const byKey = computed(() =>
    Object.fromEntries((data.value ?? []).map((item) => [item.key, item])),
  );

  function button(
    key: string,
    fallback: { text: string; variant: SiteButtonVariant; url: string },
  ): Pick<SiteButton, "text" | "variant" | "url"> {
    const found = byKey.value[key];
    return found ? { text: found.text, variant: found.variant, url: found.url } : fallback;
  }

  return { buttons: byKey, button };
}
