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
 *
 * contactMethod/source are `.optional()` here — NOT `z.enum(..., {
 * required_error })` — and "required" is instead enforced below via
 * `superRefine`. A bare `z.enum` rejects `undefined` with an `invalid_type`
 * issue, which makes Zod abort the *whole object's* parse the instant one
 * fires; an aborted object parse skips every `superRefine` below it
 * entirely, including the one that checks `socialNetwork`. So as long as
 * contactMethod/source stayed required this way, leaving them (and
 * socialNetwork) all unset at once — the normal state of a freshly opened
 * form — meant the submit button's first click reported errors for every
 * *other* field but silently dropped the socialNetwork one; it only
 * appeared once contactMethod/source were already filled in on a second
 * attempt. `.optional()` can never abort on `undefined`, so the object
 * parse always reaches "dirty" at worst, and superRefine always runs.
 */
export const applyFormSchema = z
  .object({
    ...baseFields,
    contactMethod: z.enum(["call", "telegram", "whatsapp", "max"]).optional(),
    source: z.enum(["referral", "ads", "internet", "social", "maps"]).optional(),
    // Only shown (and required) once `source` is "social" — see
    // ApplyForm.vue's conditional UiSelect.
    socialNetwork: z.enum(["telegram", "vk", "pinterest", "instagram", "youtube"]).optional(),
  })
  .superRefine((data, ctx) => {
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
    if (data.source === "social" && !data.socialNetwork) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["socialNetwork"],
        message: "Пожалуйста, укажите соц. сеть",
      });
    }
  });

export type ApplyFormValues = z.infer<typeof applyFormSchema>;
