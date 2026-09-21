import type { Ref } from "vue";

/**
 * Pointer-based row reordering — replaces manually typing a numeric
 * sortOrder. Built on the Pointer Events API (not native HTML5 DnD, whose
 * dragstart/dragover/drop events touch and pen input never fire) so it
 * works with mouse, touch, and stylus alike. `items` must already be
 * sorted by sortOrder (as returned by the list API). `persistItem` is
 * called once per moved row with its full, updated object (PUT here is a
 * full replace, so it must carry every field, not just sortOrder) once a
 * drag settles.
 *
 * @example
 * const { draggingIndex, onPointerDown } =
 *   useAdminDragReorder(items, (item) => update(item.id, item));
 * // <tr :data-row-index="index" :class="draggingIndex === index ? 'opacity-50' : ''">
 * //   <AdminDragHandle @pointerdown="onPointerDown(index, $event)" />
 * // </tr>
 * // AdminDragHandle must have `touch-action: none` so mobile browsers hand
 * // the gesture to JS instead of scrolling the page.
 */
export function useAdminDragReorder<T extends { id: number; sortOrder: number }>(
  items: Ref<T[]>,
  persistItem: (item: T) => Promise<unknown>,
) {
  const draggingIndex = ref<number | null>(null);

  function rowIndexAt(x: number, y: number): number | null {
    const row = document.elementFromPoint(x, y)?.closest<HTMLElement>("[data-row-index]");
    if (!row) return null;
    const index = Number(row.dataset.rowIndex);
    return Number.isNaN(index) ? null : index;
  }

  function onPointerMove(event: PointerEvent) {
    if (draggingIndex.value === null) return;
    event.preventDefault();
    const index = rowIndexAt(event.clientX, event.clientY);
    if (index === null || index === draggingIndex.value) return;
    const reordered = [...items.value];
    const [moved] = reordered.splice(draggingIndex.value, 1);
    reordered.splice(index, 0, moved!);
    items.value = reordered;
    draggingIndex.value = index;
  }

  async function onPointerUp() {
    window.removeEventListener("pointermove", onPointerMove);
    window.removeEventListener("pointerup", onPointerUp);
    window.removeEventListener("pointercancel", onPointerUp);
    if (draggingIndex.value === null) return;
    draggingIndex.value = null;
    const changedIndexes = items.value
      .map((item, index) => (item.sortOrder === index ? null : index))
      .filter((index): index is number => index !== null);
    items.value = items.value.map((item, index) => ({ ...item, sortOrder: index }));
    await Promise.all(changedIndexes.map((index) => persistItem(items.value[index]!)));
  }

  function onPointerDown(index: number, event: PointerEvent) {
    if (event.pointerType === "mouse" && event.button !== 0) return;
    event.preventDefault();
    draggingIndex.value = index;
    window.addEventListener("pointermove", onPointerMove);
    window.addEventListener("pointerup", onPointerUp);
    window.addEventListener("pointercancel", onPointerUp);
  }

  return { draggingIndex, onPointerDown };
}
