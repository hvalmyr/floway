import { describe, expect, it } from "vitest";
import { formatPhoneMask } from "../app/composables/usePhoneMask";

describe("formatPhoneMask", () => {
  it("formats a full number with the leading 8 into the +7 mask", () => {
    expect(formatPhoneMask("89261234567")).toBe("+7 926 123 45 67");
  });

  it("formats a full number with the leading 7", () => {
    expect(formatPhoneMask("79261234567")).toBe("+7 926 123 45 67");
  });

  it("formats partial input progressively", () => {
    expect(formatPhoneMask("926")).toBe("+7 926");
    expect(formatPhoneMask("9261")).toBe("+7 926 1");
  });

  it("ignores non-digit characters", () => {
    expect(formatPhoneMask("+7 (926) 123-45-67")).toBe("+7 926 123 45 67");
  });

  it("keeps extra digits beyond 11 appended raw, past the point AsYouType can format them", () => {
    expect(formatPhoneMask("892612345678888")).toBe("+7 92612345678888");
  });
});
