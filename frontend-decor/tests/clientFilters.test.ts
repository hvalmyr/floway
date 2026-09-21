import { describe, expect, it } from "vitest";
import { filterClients, sortClients } from "../app/lib/clientFilters";
import type { ClientListItem } from "../app/types/api";

let nextId = 1;

function makeItem(overrides: Partial<ClientListItem> = {}): ClientListItem {
  const id = nextId++;
  return {
    id,
    name: "Мария",
    phone: "+79991234567",
    email: "maria@example.com",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    productTags: [],
    clientTypeTags: [],
    requestCount: 1,
    ...overrides,
  };
}

describe("filterClients", () => {
  it("returns everything when no filters are set", () => {
    const items = [makeItem(), makeItem()];
    expect(
      filterClients(items, { query: "", productTagIds: new Set(), clientTypeTagIds: new Set() }),
    ).toHaveLength(2);
  });

  describe("search", () => {
    it("matches by name", () => {
      const items = [makeItem({ name: "Иван Петров" })];
      expect(
        filterClients(items, {
          query: "петров",
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
        }),
      ).toHaveLength(1);
    });

    it("matches by phone", () => {
      const items = [makeItem({ phone: "+79991112233" })];
      expect(
        filterClients(items, {
          query: "9991112233",
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
        }),
      ).toHaveLength(1);
    });

    it("matches by email", () => {
      const items = [makeItem({ email: "ivan@example.com" })];
      expect(
        filterClients(items, {
          query: "ivan@",
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
        }),
      ).toHaveLength(1);
    });

    it("excludes non-matching clients", () => {
      const items = [makeItem({ name: "Иван" })];
      expect(
        filterClients(items, {
          query: "zzz",
          productTagIds: new Set(),
          clientTypeTagIds: new Set(),
        }),
      ).toHaveLength(0);
    });
  });

  describe("tag filters", () => {
    it("keeps clients matching any selected product tag", () => {
      const items = [
        makeItem({ productTags: [{ id: 1, name: "Курс", color: "#f3f4f6" }] }),
        makeItem({ productTags: [{ id: 2, name: "МК", color: "#f3f4f6" }] }),
      ];
      const result = filterClients(items, {
        query: "",
        productTagIds: new Set([1]),
        clientTypeTagIds: new Set(),
      });
      expect(result).toHaveLength(1);
      expect(result[0]?.productTags[0]?.id).toBe(1);
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
      const result = filterClients(items, {
        query: "",
        productTagIds: new Set([1]),
        clientTypeTagIds: new Set([10]),
      });
      expect(result).toHaveLength(1);
    });
  });
});

describe("sortClients", () => {
  it("sorts by activity (updatedAt) descending by default", () => {
    const older = makeItem({ updatedAt: "2026-01-01T00:00:00Z" });
    const newer = makeItem({ updatedAt: "2026-01-05T00:00:00Z" });
    const result = sortClients([older, newer], "activity");
    expect(result[0]).toBe(newer);
    expect(result[1]).toBe(older);
  });

  it("sorts alphabetically by name", () => {
    const b = makeItem({ name: "Борис" });
    const a = makeItem({ name: "Анна" });
    const result = sortClients([b, a], "name");
    expect(result).toEqual([a, b]);
  });

  it("does not mutate the input array", () => {
    const items = [
      makeItem({ updatedAt: "2026-01-01T00:00:00Z" }),
      makeItem({ updatedAt: "2026-01-05T00:00:00Z" }),
    ];
    const original = [...items];
    sortClients(items, "activity");
    expect(items).toEqual(original);
  });
});
