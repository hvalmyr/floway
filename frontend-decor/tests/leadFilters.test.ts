import { describe, expect, it } from "vitest";
import { filterLeads, formatLeadExcerpt, lastContactAt, sortLeads } from "../app/lib/leadFilters";
import type { LeadListItem } from "../app/types/api";

let nextId = 1;

function makeItem(overrides: Partial<LeadListItem> = {}): LeadListItem {
  const id = nextId++;
  return {
    id,
    name: "Мария",
    phone: "+79991234567",
    email: "maria@example.com",
    contactMethod: "call",
    source: "ads",
    requestType: "course",
    relatedSlug: "",
    status: "new",
    createdAt: "2026-01-01T00:00:00Z",
    clientId: id,
    needsStatusReview: false,
    client: {
      id,
      name: "Мария",
      phone: "+79991234567",
      email: "maria@example.com",
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-01T00:00:00Z",
    },
    productTags: [],
    clientTypeTags: [],
    ...overrides,
  };
}

describe("filterLeads", () => {
  it("returns everything when no filters are set", () => {
    const items = [makeItem(), makeItem()];
    expect(
      filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set(),
        clientTypeTagIds: new Set(),
        staleOnly: false,
      }),
    ).toHaveLength(2);
  });

  describe("status filter", () => {
    it("passes everything through when the status set is empty", () => {
      const items = [makeItem({ status: "new" }), makeItem({ status: "closed_won" })];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set(),
        clientTypeTagIds: new Set(),
        staleOnly: false,
      });
      expect(result).toHaveLength(2);
    });

    it("keeps only the single selected status", () => {
      const items = [makeItem({ status: "new" }), makeItem({ status: "closed_won" })];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(["closed_won"]),
        productTagIds: new Set(),
        clientTypeTagIds: new Set(),
        staleOnly: false,
      });
      expect(result).toHaveLength(1);
      expect(result[0]?.status).toBe("closed_won");
    });

    it("keeps any of multiple selected statuses", () => {
      const items = [
        makeItem({ status: "new" }),
        makeItem({ status: "postponed" }),
        makeItem({ status: "closed_lost" }),
      ];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(["new", "postponed"]),
        productTagIds: new Set(),
        clientTypeTagIds: new Set(),
        staleOnly: false,
      });
      expect(result).toHaveLength(2);
    });
  });

  describe("tag filters", () => {
    it("keeps items matching any selected product tag", () => {
      const items = [
        makeItem({ productTags: [{ id: 1, name: "Курс", color: "#f3f4f6" }] }),
        makeItem({ productTags: [{ id: 2, name: "МК", color: "#f3f4f6" }] }),
      ];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set([1]),
        clientTypeTagIds: new Set(),
        staleOnly: false,
      });
      expect(result).toHaveLength(1);
      expect(result[0]?.productTags[0]?.id).toBe(1);
    });

    it("keeps items matching any selected client-type tag, independent of product tags", () => {
      const items = [
        makeItem({ clientTypeTags: [{ id: 10, name: "Постоянный", color: "#f3f4f6" }] }),
        makeItem({ clientTypeTags: [{ id: 20, name: "Корпоратив", color: "#f3f4f6" }] }),
      ];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set(),
        clientTypeTagIds: new Set([20]),
        staleOnly: false,
      });
      expect(result).toHaveLength(1);
      expect(result[0]?.clientTypeTags[0]?.id).toBe(20);
    });

    it("combines both tag filters as AND across types", () => {
      const items = [
        makeItem({
          productTags: [{ id: 1, name: "Курс", color: "#f3f4f6" }],
          clientTypeTags: [{ id: 10, name: "Постоянный", color: "#f3f4f6" }],
        }),
        makeItem({
          productTags: [{ id: 1, name: "Курс", color: "#f3f4f6" }],
          clientTypeTags: [{ id: 20, name: "Корпоратив", color: "#f3f4f6" }],
        }),
      ];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set([1]),
        clientTypeTagIds: new Set([10]),
        staleOnly: false,
      });
      expect(result).toHaveLength(1);
    });
  });

  describe("search", () => {
    it("matches by client name", () => {
      const items = [
        makeItem({
          client: {
            id: 1,
            name: "Иван Петров",
            phone: "1",
            email: "",
            createdAt: "",
            updatedAt: "",
          },
        }),
      ];
      expect(
        filterLeads(items, {
          query: "петров",
          statuses: new Set(),
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
          staleOnly: false,
        }),
      ).toHaveLength(1);
    });

    it("matches by phone", () => {
      const items = [
        makeItem({
          client: {
            id: 1,
            name: "Иван",
            phone: "+79991112233",
            email: "",
            createdAt: "",
            updatedAt: "",
          },
        }),
      ];
      expect(
        filterLeads(items, {
          query: "9991112233",
          statuses: new Set(),
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
          staleOnly: false,
        }),
      ).toHaveLength(1);
    });

    it("matches by email", () => {
      const items = [
        makeItem({
          client: {
            id: 1,
            name: "Иван",
            phone: "1",
            email: "ivan@example.com",
            createdAt: "",
            updatedAt: "",
          },
        }),
      ];
      expect(
        filterLeads(items, {
          query: "ivan@",
          statuses: new Set(),
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
          staleOnly: false,
        }),
      ).toHaveLength(1);
    });

    it("excludes non-matching clients", () => {
      const items = [
        makeItem({
          client: { id: 1, name: "Иван", phone: "1", email: "", createdAt: "", updatedAt: "" },
        }),
      ];
      expect(
        filterLeads(items, {
          query: "zzz",
          statuses: new Set(),
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
          staleOnly: false,
        }),
      ).toHaveLength(0);
    });
  });

  describe("staleOnly (no reply in 3+ days) preset", () => {
    const now = new Date("2026-01-10T00:00:00Z");

    it("excludes an item contacted just under 3 days ago", () => {
      const items = [makeItem({ createdAt: "2026-01-07T00:00:01Z" })]; // 2d 23h 59m 59s ago
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set(),
        clientTypeTagIds: new Set(),
        staleOnly: true,
        now,
      });
      expect(result).toHaveLength(0);
    });

    it("includes an item contacted exactly 3 days ago", () => {
      const items = [makeItem({ createdAt: "2026-01-07T00:00:00Z" })];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set(),
        clientTypeTagIds: new Set(),
        staleOnly: true,
        now,
      });
      expect(result).toHaveLength(1);
    });

    it("uses the latest comment date over the creation date when present", () => {
      const items = [
        makeItem({ createdAt: "2020-01-01T00:00:00Z", latestCommentAt: "2026-01-09T00:00:00Z" }), // commented on yesterday -> not stale
      ];
      const result = filterLeads(items, {
        query: "",
        statuses: new Set(),
        productTagIds: new Set(),
        clientTypeTagIds: new Set(),
        staleOnly: true,
        now,
      });
      expect(result).toHaveLength(0);
    });
  });
});

describe("lastContactAt", () => {
  it("falls back to createdAt when there is no comment yet", () => {
    const item = makeItem({ createdAt: "2026-01-01T00:00:00Z", latestCommentAt: undefined });
    expect(lastContactAt(item)).toBe("2026-01-01T00:00:00Z");
  });

  it("prefers the latest comment date when one exists", () => {
    const item = makeItem({
      createdAt: "2026-01-01T00:00:00Z",
      latestCommentAt: "2026-01-05T00:00:00Z",
    });
    expect(lastContactAt(item)).toBe("2026-01-05T00:00:00Z");
  });
});

describe("formatLeadExcerpt", () => {
  it("prefers the resolved program name over the raw slug", () => {
    const item = makeItem({
      requestType: "course",
      relatedSlug: "aktualnaya-floristika",
      relatedName: "Актуальная флористика",
      source: "ads",
    });
    expect(formatLeadExcerpt(item)).toBe("Курс · Актуальная флористика · Реклама");
  });

  it("falls back to the raw slug when the program couldn't be resolved", () => {
    const item = makeItem({
      requestType: "course",
      relatedSlug: "aktualnaya-floristika",
      relatedName: undefined,
      source: "ads",
    });
    expect(formatLeadExcerpt(item)).toBe("Курс · aktualnaya-floristika · Реклама");
  });

  it("omits the program segment when there is no related slug", () => {
    const item = makeItem({ requestType: "trial_lesson", relatedSlug: "", source: "internet" });
    expect(formatLeadExcerpt(item)).toBe("Пробный урок · Интернет");
  });
});

describe("sortLeads", () => {
  it("sorts by createdAt descending by default", () => {
    const older = makeItem({ createdAt: "2026-01-01T00:00:00Z" });
    const newer = makeItem({ createdAt: "2026-01-05T00:00:00Z" });
    const result = sortLeads([older, newer], "createdAt");
    expect(result[0]).toBe(newer);
    expect(result[1]).toBe(older);
  });

  it("sorts by next reminder ascending, nulls last", () => {
    const noReminder = makeItem({ createdAt: "2026-01-05T00:00:00Z" });
    const dueSoon = makeItem({
      createdAt: "2026-01-01T00:00:00Z",
      nextReminderAt: "2026-01-03T00:00:00Z",
    });
    const dueLater = makeItem({
      createdAt: "2026-01-01T00:00:00Z",
      nextReminderAt: "2026-01-10T00:00:00Z",
    });
    const result = sortLeads([noReminder, dueLater, dueSoon], "nextAction");
    expect(result).toEqual([dueSoon, dueLater, noReminder]);
  });

  it("tie-breaks items with no reminder by newest-first", () => {
    const older = makeItem({ createdAt: "2026-01-01T00:00:00Z" });
    const newer = makeItem({ createdAt: "2026-01-05T00:00:00Z" });
    const result = sortLeads([older, newer], "nextAction");
    expect(result).toEqual([newer, older]);
  });

  it("does not mutate the input array", () => {
    const items = [
      makeItem({ createdAt: "2026-01-01T00:00:00Z" }),
      makeItem({ createdAt: "2026-01-05T00:00:00Z" }),
    ];
    const original = [...items];
    sortLeads(items, "createdAt");
    expect(items).toEqual(original);
  });
});
