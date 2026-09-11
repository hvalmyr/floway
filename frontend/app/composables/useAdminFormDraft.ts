import type { Ref } from "vue";

/**
 * Persists a set of admin form refs to localStorage as the user types, and
 * restores them on mount — so a reload (or an accidental tab close) mid-edit
 * doesn't lose what was typed. Call `clearDraft()` once the form is
 * successfully submitted or the edit is cancelled, otherwise the stale draft
 * would resurrect itself on the next visit.
 *
 * @example
 * const editingId = ref<number | null>(null);
 * const form = ref(emptyForm());
 * const { clearDraft } = useAdminFormDraft("blog-posts", { editingId, form });
 */
export function useAdminFormDraft<S extends Record<string, Ref<unknown>>>(key: string, state: S) {
  const storageKey = `admin-draft:${key}`;

  if (import.meta.client) {
    try {
      const raw = localStorage.getItem(storageKey);
      if (raw) {
        const parsed = JSON.parse(raw) as Record<string, unknown>;
        for (const field of Object.keys(state)) {
          if (field in parsed) state[field]!.value = parsed[field];
        }
      }
    } catch {
      // Corrupt or inaccessible storage — start from whatever the caller
      // already initialized the refs to.
    }
  }

  let timeout: ReturnType<typeof setTimeout> | null = null;
  function persist() {
    if (!import.meta.client) return;
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(() => {
      const snapshot: Record<string, unknown> = {};
      for (const field of Object.keys(state)) snapshot[field] = state[field]!.value;
      try {
        localStorage.setItem(storageKey, JSON.stringify(snapshot));
      } catch {
        // Storage full/unavailable — draft-saving is best-effort.
      }
    }, 400);
  }

  for (const field of Object.keys(state)) {
    watch(state[field]!, persist, { deep: true });
  }

  function clearDraft() {
    if (timeout) clearTimeout(timeout);
    if (!import.meta.client) return;
    localStorage.removeItem(storageKey);
  }

  return { clearDraft };
}
