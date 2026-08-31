import type { ComputedRef, Ref } from "vue";

/**
 * Checkbox-select state for admin list bulk actions (delete/status-change
 * multiple rows at once). `items` should be the currently VISIBLE rows
 * (i.e. already search/filtered) — selecting "all" only selects what's on
 * screen, not rows hidden by a filter.
 */
export function useAdminBulkSelect<T extends { id: number }>(items: Ref<T[]> | ComputedRef<T[]>) {
  const selectedIds = ref<Set<number>>(new Set());

  function isSelected(id: number) {
    return selectedIds.value.has(id);
  }

  function toggle(id: number) {
    const next = new Set(selectedIds.value);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    selectedIds.value = next;
  }

  const allSelected = computed(
    () => items.value.length > 0 && items.value.every((item) => selectedIds.value.has(item.id)),
  );

  function toggleAll() {
    selectedIds.value = allSelected.value ? new Set() : new Set(items.value.map((item) => item.id));
  }

  function clear() {
    selectedIds.value = new Set();
  }

  return { selectedIds, isSelected, toggle, allSelected, toggleAll, clear };
}
