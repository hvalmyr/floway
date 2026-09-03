import type { Ref } from "vue";
import type { PageContent } from "~/types/api";

/**
 * Same data layer as the flat /admin/page-content list, but scoped to a
 * fixed set of keys — used by the grouped admin pages (hero block, homepage
 * copy, contacts/legal) so each shows only the page_content rows relevant
 * to it instead of every key in the table.
 */
export function useAdminPageContent(keys: string[]) {
  const api = useApiClient();
  const allItems = ref<PageContent[]>([]) as Ref<PageContent[]>;
  const loading = ref(false);
  const error = ref("");
  const savingKey = ref<string | null>(null);
  const savedKey = ref<string | null>(null);

  const items = computed(() =>
    keys
      .map((key) => allItems.value.find((item) => item.key === key))
      .filter((item): item is PageContent => item !== undefined),
  );

  async function fetchAll() {
    loading.value = true;
    error.value = "";
    try {
      allItems.value = (await api<PageContent[]>("/api/v1/page-content")) ?? [];
    } catch {
      error.value = "Не удалось загрузить данные";
    } finally {
      loading.value = false;
    }
  }

  async function save(item: PageContent) {
    savingKey.value = item.key;
    savedKey.value = null;
    try {
      const updated = await api<PageContent>(`/api/v1/page-content/${item.key}`, {
        method: "PUT",
        body: { value: item.value },
      });
      item.value = updated.value;
      savedKey.value = item.key;
    } catch {
      error.value = `Не удалось сохранить «${item.label}»`;
    } finally {
      savingKey.value = null;
    }
  }

  function onImageUploaded(item: PageContent, url: string) {
    item.value = url;
    save(item);
  }

  return { items, loading, error, savingKey, savedKey, fetchAll, save, onImageUploaded };
}
