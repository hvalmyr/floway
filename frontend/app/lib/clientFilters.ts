import type { ClientListItem } from "~/types/api";

export interface ClientFilterOptions {
  query: string;
  /** Empty set = no product-tag filter applied. */
  productTagIds: Set<number>;
  /** Empty set = no client-type-tag filter applied. */
  clientTypeTagIds: Set<number>;
}

function matchesAnyTag(itemTagIds: number[], wanted: Set<number>): boolean {
  if (wanted.size === 0) return true;
  return itemTagIds.some((id) => wanted.has(id));
}

export function filterClients(
  items: ClientListItem[],
  opts: ClientFilterOptions,
): ClientListItem[] {
  const query = opts.query.trim().toLowerCase();

  return items.filter((item) => {
    if (
      !matchesAnyTag(
        item.productTags.map((t) => t.id),
        opts.productTagIds,
      )
    )
      return false;
    if (
      !matchesAnyTag(
        item.clientTypeTags.map((t) => t.id),
        opts.clientTypeTagIds,
      )
    )
      return false;

    if (query) {
      const haystack = `${item.name} ${item.phone} ${item.email}`.toLowerCase();
      if (!haystack.includes(query)) return false;
    }

    return true;
  });
}

export type ClientSortMode = "activity" | "name";

/** Default sort mirrors the backend's own ordering (most recently active
 * first, by updatedAt); the alternate sort is alphabetical by name. */
export function sortClients(items: ClientListItem[], by: ClientSortMode): ClientListItem[] {
  const copy = [...items];
  if (by === "name") {
    return copy.sort((a, b) => a.name.localeCompare(b.name, "ru"));
  }
  return copy.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
}
