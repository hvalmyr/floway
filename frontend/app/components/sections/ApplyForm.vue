<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import { applyFormSchema, simpleApplyFormSchema } from "~/lib/validation/applyForm";
import type { ContactMethod, LeadRequestType, LeadSource } from "~/types/api";

/**
 * Reusable lead-capture form embedded on the home page (trial lesson, in its
 * simplified 3-field `variant="simple"` form), a course page, and a
 * masterclass page (both `variant="full"`, the default, with both radio
 * groups). `context` is sent as the lead's requestType so every submission
 * records where it came from.
 *
 * @example
 * <ApplyForm context="course" :related-id="course.id" title="Записаться на курс" />
 * <ApplyForm context="trial_lesson" variant="simple" title="Пробное занятие" />
 */
const props = withDefaults(
  defineProps<{
    context: LeadRequestType;
    relatedId?: number;
    title?: string;
    variant?: "full" | "simple";
  }>(),
  { title: "Оставить заявку", relatedId: undefined, variant: "full" },
);

const api = useApi();

// Schema (and therefore the fields vee-validate tracks) is picked once from
// a prop that doesn't change after mount, so calling useForm() a single time
// here doesn't violate the rules of hooks. The two schemas produce
// differently-shaped values, so the submit handler below reads fields
// defensively rather than relying on one static inferred type.
const schema = props.variant === "simple" ? simpleApplyFormSchema : applyFormSchema;
const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(schema),
  initialValues: { name: "", phone: "", email: "" },
});

const status = ref<"idle" | "success" | "error">("idle");
const submitError = ref("");

const contactMethodOptions = [
  { value: "call", label: "Позвоните мне" },
  { value: "telegram", label: "Напишите мне в Telegram" },
  { value: "whatsapp", label: "Напишите мне в Whatsapp" },
  { value: "max", label: "Напишите мне в Max" },
];

const sourceOptions = [
  { value: "referral", label: "По рекомендации" },
  { value: "ads", label: "Реклама" },
  { value: "internet", label: "В интернете" },
  { value: "social", label: "В социальных сетях" },
  { value: "maps", label: "В картах" },
];

const onSubmit = handleSubmit(async (values) => {
  status.value = "idle";
  submitError.value = "";
  const raw = values as Record<string, string | undefined>;
  try {
    await api.submitApplication({
      name: raw.name!,
      phone: raw.phone!,
      email: raw.email || undefined,
      // Простая форма (пробное занятие) не спрашивает способ связи/источник —
      // бэкенду они всё равно нужны, подставляем разумные значения по умолчанию.
      contactMethod: (raw.contactMethod as ContactMethod | undefined) ?? "call",
      source: (raw.source as LeadSource | undefined) ?? "internet",
      requestType: props.context,
      relatedId: props.relatedId,
    });
    status.value = "success";
  } catch (err) {
    status.value = "error";
    submitError.value =
      (err as { message?: string } | undefined)?.message ?? "Не удалось отправить заявку";
  }
});
</script>

<template>
  <div class="rounded-lg bg-surface p-32 sm:p-48">
    <h2 v-if="title" class="mb-24 text-center font-display text-h2 text-primary">{{ title }}</h2>

    <div
      v-if="status === 'success'"
      class="flex flex-col items-center gap-16 py-32 text-center"
      role="status"
    >
      <p class="font-display text-h3 text-primary">Заявка отправлена!</p>
      <p class="font-body text-body text-ink">Мы свяжемся с вами в ближайшее время.</p>
    </div>

    <form v-else class="flex flex-col gap-24" novalidate @submit="onSubmit">
      <UiInput name="name" label="Ваше имя" required autocomplete="name" />
      <UiPhoneInput name="phone" label="Номер телефона" required />
      <UiInput name="email" label="Почта" type="email" autocomplete="email" />
      <template v-if="variant === 'full'">
        <UiRadioGroup
          name="contactMethod"
          label="Как с вами связаться?"
          :options="contactMethodOptions"
        />
        <UiRadioGroup name="source" label="Откуда вы о нас узнали?" :options="sourceOptions" />
      </template>

      <p class="font-body text-body text-ink">
        Нажимая на кнопку, вы даёте согласие на обработку персональных данных и соглашаетесь с
        <NuxtLink to="/privacy" class="text-primary underline"
          >политикой конфиденциальности</NuxtLink
        >.
      </p>

      <p v-if="status === 'error'" class="font-body text-body font-bold text-ink" role="alert">
        {{ submitError }}
      </p>

      <UiButton type="submit" :pill="false" block :loading="isSubmitting" :disabled="isSubmitting">
        {{ variant === "simple" ? "Записаться на занятие" : "Отправить заявку" }}
      </UiButton>
    </form>
  </div>
</template>
