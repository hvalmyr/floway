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
  // `.optional()` here (not z.enum's own `required_error`) — "required" is
  // enforced below via superRefine instead. A bare `z.enum` rejects
  // `undefined` with an `invalid_type` issue, which makes Zod abort the
  // *whole object's* parse the instant one fires; an aborted object parse
  // skips every superRefine below it entirely, including the one that
  // checks socialNetwork/objectTypeId. So as long as these stayed required
  // this way, leaving them all unset at once — the normal state of a
  // freshly opened form — meant the submit button's first click reported
  // errors for every *other* field but silently dropped e.g. the
  // socialNetwork one; it only appeared once the rest were already filled
  // in on a second attempt. `.optional()` can never abort on `undefined`,
  // so the object parse always reaches "dirty" at worst, and superRefine
  // always runs.
  contactMethod: z.enum(["call", "telegram", "whatsapp", "max"]).optional(),
  source: z.enum(["referral", "ads", "internet", "social", "maps"]).optional(),
  format: z.enum(["season", "event"]).optional(),
  // Only shown (and required) once `source` is "social" — see
  // ApplyForm.vue's conditional UiSelect.
  socialNetwork: z.enum(["telegram", "vk", "pinterest", "instagram", "youtube"]).optional(),
};

function validateRequiredFields(
  data: { contactMethod?: string; source?: string; format?: string; socialNetwork?: string },
  ctx: z.RefinementCtx,
) {
  if (!data.contactMethod) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ["contactMethod"],
      message: "Пожалуйста, выберите способ связи",
    });
  }
  if (!data.source) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ["source"],
      message: "Пожалуйста, укажите, как вы о нас узнали",
    });
  }
  if (!data.format) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ["format"],
      message: "Пожалуйста, выберите формат",
    });
  }
  if (data.source === "social" && !data.socialNetwork) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ["socialNetwork"],
      message: "Пожалуйста, укажите соц. сеть",
    });
  }
}

/** Used when the form has no preset object type (the home page) — the
 * visitor picks one from ApplyForm's own select. */
export const applyFormSchemaWithObjectType = z
  .object({
    ...baseFields,
    // Same `.optional()` + superRefine reasoning as above.
    objectTypeId: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    if (!data.objectTypeId) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["objectTypeId"],
        message: "Пожалуйста, выберите тип объекта",
      });
    }
    validateRequiredFields(data, ctx);
  });

/** Used when embedded on a landing page — objectTypeId comes from the page
 * itself (a prop), not a form field, so there's nothing to validate here. */
export const applyFormSchema = z.object(baseFields).superRefine(validateRequiredFields);

export type ApplyFormValues = z.infer<typeof applyFormSchemaWithObjectType>;
