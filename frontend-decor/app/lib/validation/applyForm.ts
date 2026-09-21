import { isValidPhoneNumber } from "libphonenumber-js/min";
import { z } from "zod";

const baseFields = {
  name: z.string().trim().min(1, "Пожалуйста, укажите имя").max(100, "Слишком длинное имя"),
  phone: z
    .string()
    .trim()
    .min(1, "Пожалуйста, укажите номер телефона")
    .refine((v) => isValidPhoneNumber(v, "RU"), "Проверьте номер телефона — введите полный номер"),
  email: z.union([z.string().trim().email("Проверьте адрес почты"), z.literal("")]).optional(),
  consent: z.boolean().refine((v) => v === true, "Нужно согласие на обработку персональных данных"),
  contactMethod: z.enum(["call", "telegram", "whatsapp", "max"], {
    required_error: "Пожалуйста, выберите способ связи",
  }),
  source: z.enum(["referral", "ads", "internet", "social", "maps"], {
    required_error: "Пожалуйста, укажите, как вы о нас узнали",
  }),
  format: z.enum(["season", "event"], {
    required_error: "Пожалуйста, выберите формат",
  }),
};

/** Used when the form has no preset object type (the home page) — the
 * visitor picks one from ApplyForm's own select. */
export const applyFormSchemaWithObjectType = z.object({
  ...baseFields,
  objectTypeId: z.string().min(1, "Пожалуйста, выберите тип объекта"),
});

/** Used when embedded on a landing page — objectTypeId comes from the page
 * itself (a prop), not a form field, so there's nothing to validate here. */
export const applyFormSchema = z.object(baseFields);

export type ApplyFormValues = z.infer<typeof applyFormSchemaWithObjectType>;
