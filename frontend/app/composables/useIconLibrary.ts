import type { Icon } from "~/types/api";

/**
 * The uploaded icon library (see AppIcon.vue for how a Feature/PageContent
 * icon value resolves to one of these vs. a built-in FEATURE_ICONS key).
 * `useAsyncData`'s key-based dedup means every call on one page — one per
 * <AppIcon>/<AdminIconPicker> instance — shares a single fetch, not one
 * request per icon.
 */
export function useIconLibrary() {
  const api = useApi();
  const { data, refresh } = useAsyncData("icon-library", () => api.getIcons());

  function findIcon(id: number): Icon | undefined {
    return data.value?.find((icon) => icon.id === id);
  }

  return { icons: computed(() => data.value ?? []), findIcon, refresh };
}
