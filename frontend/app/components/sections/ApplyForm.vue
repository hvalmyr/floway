<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import { applyFormSchema } from "~/lib/validation/applyForm";
import type { LeadRequestType } from "~/types/api";

/**
 * Reusable lead-capture form embedded on the home page (trial lesson), a
 * course page, and a masterclass page. `context` is sent as the lead's
 * requestType so every submission records where it came from.
 *
 * @example
 * <ApplyForm context="course" :related-id="course.id" title="Записаться на курс" />
 */
const props = withDefaults(
  defineProps<{
    context: LeadRequestType;
    relatedId?: number;
    title?: string;
  }>(),
  { title: "Оставить заявку", relatedId: undefined },
);

const api = useApi();

const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(applyFormSchema),
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
  try {
    await api.submitApplication({
      name: values.name,
      phone: values.phone,
      email: values.email || undefined,
      contactMethod: values.contactMethod,
      source: values.source,
      requestType: props.context,
      relatedId: props.relatedId,
    });
    status.value = "success";
  } catch (err) {
    status.value = "error";
    submitError.value = (err as { message?: string } | undefined)?.message ?? "Не удалось отправить заявку";
  }
});
</script>

<template>
  <div class="rounded-lg bg-surface p-32 sm:p-48">
    <h2 class="mb-24 font-display text-h2 text-ink-900">{{ title }}</h2>

    <div v-if="status === 'success'" class="flex flex-col items-center gap-16 py-32 text-center" role="status">
      <p class="font-display text-h3 text-success">Заявка отправлена!</p>
      <p class="text-body text-ink-700">Мы свяжемся с вами в ближайшее время.</p>
    </div>

    <form v-else class="flex flex-col gap-24" novalidate @submit="onSubmit">
      <UiInput name="name" label="Ваше имя" required autocomplete="name" />
      <UiPhoneInput name="phone" label="Номер телефона" required />
      <UiInput name="email" label="Почта" type="email" autocomplete="email" />
      <UiRadioGroup name="contactMethod" label="Как с вами связаться?" :options="contactMethodOptions" />
      <UiRadioGroup name="source" label="Откуда вы о нас узнали?" :options="sourceOptions" />

      <p class="text-small text-ink-400">
        Нажимая «Отправить заявку», вы соглашаетесь с
        <NuxtLink to="/privacy" class="text-primary-700 underline">
          политикой обработки персональных данных</NuxtLink
        >.
      </p>

      <p v-if="status === 'error'" class="text-small text-error" role="alert">{{ submitError }}</p>

      <UiButton type="submit" block :loading="isSubmitting" :disabled="isSubmitting">Отправить заявку</UiButton>
    </form>
  </div>
</template>
