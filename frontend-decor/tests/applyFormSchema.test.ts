import { describe, expect, it } from "vitest";
import { applyFormSchema } from "../app/lib/validation/applyForm";

const validPayload = {
  name: "Мария",
  phone: "+7 (926) 123 45 67",
  email: "",
  consent: true as const,
  contactMethod: "call" as const,
  source: "referral" as const,
};

describe("applyFormSchema", () => {
  it("accepts a fully valid payload with an empty optional email", () => {
    expect(applyFormSchema.safeParse(validPayload).success).toBe(true);
  });

  it("accepts a valid email when provided", () => {
    const result = applyFormSchema.safeParse({ ...validPayload, email: "test@example.com" });
    expect(result.success).toBe(true);
  });

  it("rejects an invalid, non-empty email", () => {
    const result = applyFormSchema.safeParse({ ...validPayload, email: "not-an-email" });
    expect(result.success).toBe(false);
  });

  it("rejects a phone with fewer than 11 digits", () => {
    const result = applyFormSchema.safeParse({ ...validPayload, phone: "+7 (926) 123" });
    expect(result.success).toBe(false);
  });

  it("accepts an 11-digit phone even without the mask's punctuation", () => {
    // UiPhoneInput's live mask means a submitted value is never actually
    // unpunctuated in practice, but the schema itself only cares that all
    // 11 digits are there, not the exact formatting around them.
    const result = applyFormSchema.safeParse({ ...validPayload, phone: "89261234567" });
    expect(result.success).toBe(true);
  });

  it("rejects an empty name", () => {
    const result = applyFormSchema.safeParse({ ...validPayload, name: "" });
    expect(result.success).toBe(false);
  });

  it("rejects a missing consent", () => {
    const { consent: _consent, ...rest } = validPayload;
    const result = applyFormSchema.safeParse(rest);
    expect(result.success).toBe(false);
  });

  it("rejects consent: false", () => {
    const result = applyFormSchema.safeParse({ ...validPayload, consent: false });
    expect(result.success).toBe(false);
  });

  it("rejects a missing contactMethod", () => {
    const { contactMethod: _contactMethod, ...rest } = validPayload;
    const result = applyFormSchema.safeParse(rest);
    expect(result.success).toBe(false);
  });

  it("rejects a missing source", () => {
    const { source: _source, ...rest } = validPayload;
    const result = applyFormSchema.safeParse(rest);
    expect(result.success).toBe(false);
  });
});
