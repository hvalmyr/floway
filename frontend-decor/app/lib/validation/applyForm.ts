import { isValidPhoneNumber } from "libphonenumber-js/min";
import { z } from "zod";

const baseFields = {
  name: z.string().trim().min(1, "Пожалуйста, укажите имя").max(100, "Слишком длинное имя"),
  // Two distinct messages (empty vs. incomplete) rather than one regex
  // failure message — matches the redesigned form's per-state copy, and is
  // just a clearer error either way. isValidPhoneNumber (not a digit-count
  // check) so numbers from any country validate correctly — students come
  // from Belarus, Europe, South America too, and length alone doesn't
  // distinguish "too short" from "wrong country's format".
  phone: z
    .string()
    .trim()
    .min(1, "Пожалуйста, укажите номер телефона")
    .refine((v) => isValidPhoneNumber(v, "RU"), "Проверьте номер телефона — введите полный номер"),
  email: z.union([z.string().trim().email("Проверьте адрес почты"), z.literal("")]).optional(),
  // Required in the redesigned form (it wasn't a field at all before) —
  // the backend has no way to record consent otherwise, and the checkbox
  // is the actual legal basis for storing the rest of the submitted data.
  // `z.boolean().refine(...)` rather than `z.literal(true)` on purpose: the
  // literal type would make TS reject the checkbox's own unchecked (false)
  // initial state as a type error, even though "unchecked" is exactly what
  // an unfilled-out form should start as.
  consent: z.boolean().refine((v) => v === true, "Нужно согласие на обработку персональных данных"),
};

/**
 * Every ApplyForm instance (home page trial lesson, course pages,
 * masterclass pages): name/phone/email/consent + both selects, all required.
 */
export const applyFormSchema = z.object({
  ...baseFields,
  contactMethod: z.enum(["call", "telegram", "whatsapp", "max"], {
    required_error: "Пожалуйста, выберите способ связи",
  }),
  source: z.enum(["referral", "ads", "internet", "social", "maps"], {
    required_error: "Пожалуйста, укажите, как вы о нас узнали",
  }),
});

export type ApplyFormValues = z.infer<typeof applyFormSchema>;
