import { z } from "zod";

const phoneMaskPattern = /^\+7 \(\d{3}\) \d{3} \d{2} \d{2}$/;

const baseFields = {
  name: z.string().trim().min(2, "Введите имя").max(100, "Слишком длинное имя"),
  phone: z.string().regex(phoneMaskPattern, "Введите номер телефона полностью"),
  email: z.union([z.string().trim().email("Введите корректный email"), z.literal("")]).optional(),
};

/** Full form (course/masterclass pages): name/phone/email + both radio groups. */
export const applyFormSchema = z.object({
  ...baseFields,
  contactMethod: z.enum(["call", "telegram", "whatsapp", "max"], {
    required_error: "Выберите способ связи",
  }),
  source: z.enum(["referral", "ads", "internet", "social", "maps"], {
    required_error: "Выберите, откуда вы о нас узнали",
  }),
});

/**
 * Simplified form (home page trial lesson, per mockup): just name/phone/
 * email, no radio groups. contactMethod/source are still required by the
 * backend's Lead model — ApplyForm.vue fills in sensible defaults for them.
 */
export const simpleApplyFormSchema = z.object(baseFields);

export type ApplyFormValues = z.infer<typeof applyFormSchema>;
export type SimpleApplyFormValues = z.infer<typeof simpleApplyFormSchema>;
