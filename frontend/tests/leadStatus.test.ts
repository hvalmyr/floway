import { describe, expect, it } from "vitest";
import { LEAD_STATUSES, LEAD_STATUS_LABELS } from "../app/lib/leadStatus";

describe("LEAD_STATUS_LABELS", () => {
  it("has exactly one label per status, in both directions", () => {
    // Guards against the exact bug this module exists to prevent: a status
    // added to one of the two without updating the other.
    expect(Object.keys(LEAD_STATUS_LABELS).sort()).toEqual([...LEAD_STATUSES].sort());
  });

  it("has no blank labels", () => {
    for (const status of LEAD_STATUSES) {
      expect(LEAD_STATUS_LABELS[status].trim().length).toBeGreaterThan(0);
    }
  });
});
