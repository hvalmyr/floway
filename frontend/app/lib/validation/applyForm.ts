import { z } from "zod";

const phoneMaskPattern = /^\+7 \(\d{3}\) \d{3} \d{2} \d{2}$/;

/** Validation schema for <ApplyForm>. Shared by every apply-form instance on the site. */
export const applyFormSchema = z.object({
  name: z.string().trim().min(2, "Введите имя").max(100, "Слишком длинное имя"),
  phone: z.string().regex(phoneMaskPattern, "Введите номер телефона полностью"),
  email: z.union([z.string().trim().email("Введите корректный email"), z.literal("")]).optional(),
  contactMethod: z.enum(["call", "telegram", "whatsapp", "max"], {
    required_error: "Выберите способ связи",
  }),
  source: z.enum(["referral", "ads", "internet", "social", "maps"], {
    required_error: "Выберите, откуда вы о нас узнали",
  }),
});

export type ApplyFormValues = z.infer<typeof applyFormSchema>;
