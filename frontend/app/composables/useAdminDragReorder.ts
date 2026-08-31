import type { Ref } from "vue";

/**
 * Native HTML5 drag-and-drop row reordering — replaces manually typing a
 * numeric sortOrder. `items` must already be sorted by sortOrder (as
 * returned by the list API). `persistItem` is called once per moved row
 * with its full, updated object (PUT here is a full replace, so it must
 * carry every field, not just sortOrder) after a drop settles.
 *
 * @example
 * const { draggingIndex, onDragStart, onDragOver, onDrop } =
 *   useAdminDragReorder(items, (item) => update(item.id, item));
 * // <tr draggable="true" @dragstart="onDragStart(index)"
 * //     @dragover.prevent="onDragOver(index)" @drop.prevent="onDrop">
 */
export function useAdminDragReorder<T extends { id: number; sortOrder: number }>(
  items: Ref<T[]>,
  persistItem: (item: T) => Promise<unknown>,
) {
  const draggingIndex = ref<number | null>(null);

  function onDragStart(index: number) {
    draggingIndex.value = index;
  }

  function onDragOver(index: number) {
    if (draggingIndex.value === null || draggingIndex.value === index) return;
    const reordered = [...items.value];
    const [moved] = reordered.splice(draggingIndex.value, 1);
    reordered.splice(index, 0, moved!);
    items.value = reordered;
    draggingIndex.value = index;
  }

  async function onDrop() {
    if (draggingIndex.value === null) return;
    draggingIndex.value = null;
    const changedIndexes = items.value
      .map((item, index) => (item.sortOrder === index ? null : index))
      .filter((index): index is number => index !== null);
    items.value = items.value.map((item, index) => ({ ...item, sortOrder: index }));
    await Promise.all(changedIndexes.map((index) => persistItem(items.value[index]!)));
  }

  return { draggingIndex, onDragStart, onDragOver, onDrop };
}
