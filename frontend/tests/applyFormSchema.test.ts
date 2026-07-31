import { describe, expect, it } from "vitest";
import { applyFormSchema } from "../app/lib/validation/applyForm";

const validPayload = {
  name: "Мария",
  phone: "+7 (926) 123 45 67",
  email: "",
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

  it("rejects a phone that doesn't match the +7 mask", () => {
    const result = applyFormSchema.safeParse({ ...validPayload, phone: "89261234567" });
    expect(result.success).toBe(false);
  });

  it("rejects a name shorter than 2 characters", () => {
    const result = applyFormSchema.safeParse({ ...validPayload, name: "М" });
    expect(result.success).toBe(false);
  });

  it("rejects a missing contactMethod", () => {
    const { contactMethod: _contactMethod, ...rest } = validPayload;
    const result = applyFormSchema.safeParse(rest);
    expect(result.success).toBe(false);
  });
});
