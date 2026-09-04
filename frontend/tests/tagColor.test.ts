import { describe, expect, it } from "vitest";
import { readableTextColor } from "../app/lib/tagColor";

describe("readableTextColor", () => {
  it("picks black text on a light background", () => {
    expect(readableTextColor("#ffffff")).toBe("#000000");
    expect(readableTextColor("#f3f4f6")).toBe("#000000");
  });

  it("picks white text on a dark background", () => {
    expect(readableTextColor("#000000")).toBe("#ffffff");
    expect(readableTextColor("#1a1a2e")).toBe("#ffffff");
  });

  it("handles a hex color without a leading #", () => {
    expect(readableTextColor("ffffff")).toBe("#000000");
  });

  it("falls back to black for malformed input", () => {
    expect(readableTextColor("")).toBe("#000000");
    expect(readableTextColor("orange")).toBe("#000000");
    expect(readableTextColor("#fff")).toBe("#000000");
  });
});
