/**
 * Looks up freeform site copy by key (see PageContent in types/api.ts).
 * `fallback` is what renders if the key is somehow missing (migration not
 * run yet, typo) — pass the current hardcoded copy as a safety net rather
 * than leaving a blank spot on the page.
 *
 * @example
 * const { text } = await usePageContent();
 * // in template: {{ text('home_hero_title', 'Мы рядом с первого букета') }}
 */
export async function usePageContent() {
  const api = useApi();
  const { data } = await useAsyncData("page-content", () => api.getPageContent());

  const byKey = computed(() =>
    Object.fromEntries((data.value ?? []).map((item) => [item.key, item.value])),
  );

  function text(key: string, fallback = ""): string {
    return byKey.value[key] ?? fallback;
  }

  return { content: byKey, text };
}
