import type { LeadListItem, LeadStatus } from "~/types/api";
import { requestTypeLabels, sourceLabels } from "./leadLabels";

const MS_PER_DAY = 24 * 60 * 60 * 1000;

/** The date this lead's client was last actually heard from — the most
 * recent comment if there is one, otherwise the request's own creation
 * date. Backs both the "no reply in 3+ days" filter preset and the card's
 * displayed date. */
export function lastContactAt(item: LeadListItem): string {
  return item.latestCommentAt ?? item.createdAt;
}

/** A one-line summary of the request itself for the list card. Public
 * submissions carry no free-text message (see ApplyForm), so this is
 * synthesized from structured fields rather than lifted from a comment —
 * comments get their own dedicated feed on the client detail page. */
export function formatLeadExcerpt(item: LeadListItem): string {
  const parts = [requestTypeLabels[item.requestType] ?? item.requestType];
  // Prefer the resolved course/masterclass title; fall back to the raw
  // slug only if resolution came back empty (e.g. the program was since
  // deleted/renamed — see model.Lead.RelatedName on the backend).
  const program = item.relatedName || item.relatedSlug;
  if (program) parts.push(program);
  parts.push(sourceLabels[item.source] ?? item.source);
  return parts.join(" · ");
}

export interface LeadFilterOptions {
  query: string;
  /** Empty set = no status filter applied (show every status). */
  statuses: Set<LeadStatus>;
  /** Empty set = no product-tag filter applied. */
  productTagIds: Set<number>;
  /** Empty set = no client-type-tag filter applied. */
  clientTypeTagIds: Set<number>;
  /** "No reply in 3+ days" preset. */
  staleOnly: boolean;
  now?: Date;
}

function matchesAnyTag(itemTagIds: number[], wanted: Set<number>): boolean {
  if (wanted.size === 0) return true;
  return itemTagIds.some((id) => wanted.has(id));
}

export function filterLeads(items: LeadListItem[], opts: LeadFilterOptions): LeadListItem[] {
  const query = opts.query.trim().toLowerCase();
  const now = opts.now ?? new Date();

  return items.filter((item) => {
    if (opts.statuses.size > 0 && !opts.statuses.has(item.status)) return false;

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

    if (opts.staleOnly) {
      const daysSinceContact =
        (now.getTime() - new Date(lastContactAt(item)).getTime()) / MS_PER_DAY;
      if (daysSinceContact < 3) return false;
    }

    if (query) {
      const c = item.client;
      const haystack = `${c.name} ${c.phone} ${c.email}`.toLowerCase();
      if (!haystack.includes(query)) return false;
    }

    return true;
  });
}

export type LeadSortMode = "createdAt" | "nextAction";

/** Default sort is newest-first by creation date. The alternate mode sorts
 * by the soonest upcoming reminder — items with no open reminder sort last,
 * tie-broken by newest-first. */
export function sortLeads(items: LeadListItem[], by: LeadSortMode): LeadListItem[] {
  const copy = [...items];
  if (by === "createdAt") {
    return copy.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
  }
  return copy.sort((a, b) => {
    if (a.nextReminderAt && b.nextReminderAt) {
      return new Date(a.nextReminderAt).getTime() - new Date(b.nextReminderAt).getTime();
    }
    if (a.nextReminderAt) return -1;
    if (b.nextReminderAt) return 1;
    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
  });
}
