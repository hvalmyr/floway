import type { Ref } from "vue";
import type { FaqPage, PageFaq, PageFaqItem } from "~/types/api";

/**
 * Admin data layer for one page-scoped FAQ block (masterclasses, gift
 * certificate) — settings (title/description/visible) plus its Q&A items,
 * fetched together via the same public GET the page itself uses (there's no
 * separate admin-only read). Mirrors useAdminResource's create/update/remove
 * shape for items; settings save separately via saveSettings().
 */
export function useAdminPageFaq(page: FaqPage) {
  const api = useApiClient();
  const settings = ref({ title: "", description: "", visible: false });
  const items = ref<PageFaqItem[]>([]) as Ref<PageFaqItem[]>;
  const loading = ref(false);
  const error = ref("");
  const saving = ref(false);

  async function fetchAll() {
    loading.value = true;
    error.value = "";
    try {
      const data = await api<PageFaq>(`/api/v1/page-faq/${page}`);
      settings.value = { title: data.title, description: data.description, visible: data.visible };
      items.value = data.items;
    } catch {
      error.value = "Не удалось загрузить данные";
    } finally {
      loading.value = false;
    }
  }

  async function saveSettings() {
    saving.value = true;
    try {
      const updated = await api<PageFaq>(`/api/v1/page-faq/${page}`, {
        method: "PUT",
        body: settings.value,
      });
      settings.value = {
        title: updated.title,
        description: updated.description,
        visible: updated.visible,
      };
    } finally {
      saving.value = false;
    }
  }

  async function create(payload: Partial<PageFaqItem>) {
    const created = await api<PageFaqItem>(`/api/v1/page-faq/${page}/items`, {
      method: "POST",
      body: payload,
    });
    items.value = [...items.value, created];
    return created;
  }

  async function update(id: number, payload: Partial<PageFaqItem>) {
    const updated = await api<PageFaqItem>(`/api/v1/page-faq/${page}/items/${id}`, {
      method: "PUT",
      body: payload,
    });
    items.value = items.value.map((item) => (item.id === id ? updated : item));
    return updated;
  }

  async function remove(id: number) {
    await api(`/api/v1/page-faq/${page}/items/${id}`, { method: "DELETE" });
    items.value = items.value.filter((item) => item.id !== id);
  }

  return {
    settings,
    items,
    loading,
    error,
    saving,
    fetchAll,
    saveSettings,
    create,
    update,
    remove,
  };
}
