<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import { applyFormSchema } from "~/lib/validation/applyForm";
import type { ContactMethod, LeadRequestType, LeadSource } from "~/types/api";

/**
 * Reusable lead-capture form embedded on the home page (trial lesson), a
 * course page, and a masterclass page — same fields everywhere. `context` is
 * sent as the lead's requestType so every submission records where it came
 * from.
 *
 * `title` stays caller-driven on purpose (not a hardcoded generic string) —
 * the course and masterclass pages each pass their own wording ("...курс"
 * vs "...мастер-класс") so the form reads as being about the specific thing
 * the visitor was just looking at, not a generic contact form. `lead` is an
 * optional one-line subhead under it — empty by default (the trial-lesson
 * embed on the home page already sits under its own section heading and
 * doesn't need a second one), set it per caller where a subhead helps.
 *
 * @example
 * <ApplyForm context="course" :related-id="course.id" :related-slug="course.slug" title="Записаться на курс" />
 * <ApplyForm context="trial_lesson" title="Пробное занятие" bare />
 */
const props = withDefaults(
  defineProps<{
    context: LeadRequestType;
    relatedId?: number;
    /** The course/masterclass slug the visitor was looking at — sent
     * straight through to the lead so the admin panel can show which one
     * without cross-referencing relatedId. Leave unset for trial_lesson. */
    relatedSlug?: string;
    title?: string;
    lead?: string;
    /** Drops the white card + its padding, for a caller that already puts
     * this in its own full-bleed column (the home page's trial-lesson
     * section) instead of a padded card floating in a layout. */
    bare?: boolean;
  }>(),
  { title: "Оставить заявку", lead: "", relatedId: undefined, relatedSlug: undefined, bare: false },
);

const api = useApi();
const { text } = await usePageContent();

const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(applyFormSchema),
  initialValues: { name: "", phone: "", email: "", consent: false },
});

const status = ref<"idle" | "success" | "error">("idle");
const submitError = ref("");

const contactMethodOptions = computed(() => [
  { value: "call", label: text("apply_form_contact_method_call", "Позвоните мне") },
  {
    value: "telegram",
    label: text("apply_form_contact_method_telegram", "Напишите мне в Telegram"),
  },
  {
    value: "whatsapp",
    label: text("apply_form_contact_method_whatsapp", "Напишите мне в Whatsapp"),
  },
  { value: "max", label: text("apply_form_contact_method_max", "Напишите мне в Max") },
]);

const sourceOptions = computed(() => [
  { value: "referral", label: text("apply_form_source_referral", "По рекомендации") },
  { value: "ads", label: text("apply_form_source_ads", "Реклама") },
  { value: "internet", label: text("apply_form_source_internet", "В интернете") },
  { value: "social", label: text("apply_form_source_social", "В социальных сетях") },
  { value: "maps", label: text("apply_form_source_maps", "В картах") },
]);

const onSubmit = handleSubmit(async (values) => {
  status.value = "idle";
  submitError.value = "";
  try {
    await api.submitApplication({
      name: values.name,
      phone: values.phone,
      email: values.email || undefined,
      // `consent` is deliberately left out here — it's the frontend's own gate
      // (already enforced by the schema before handleSubmit even calls this),
      // not something the Lead model has a column for.
      contactMethod: values.contactMethod as ContactMethod,
      source: values.source as LeadSource,
      requestType: props.context,
      relatedId: props.relatedId,
      relatedSlug: props.relatedSlug,
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
  <div :class="bare ? 'w-full' : 'rounded-[30px] bg-white p-24 sm:p-48 lg:p-64'">
    <div
      v-if="title || lead || status === 'success'"
      class="mb-24 flex flex-col gap-12"
      :class="
        status === 'success' && bare
          ? 'rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150'
          : ''
      "
      :role="status === 'success' ? 'status' : undefined"
    >
      <h2 v-if="status === 'success'" class="font-display text-h2 text-ink">
        {{ text("apply_form_success_title", "заявка отправлена") }}
      </h2>
      <h2 v-else-if="title" class="font-display text-h2 text-ink">{{ title }}</h2>
      <p v-if="status === 'success'" class="font-body text-body text-ink">
        {{
          text(
            "apply_form_success_message",
            "Мы свяжемся с вами в ближайшее время удобным для вас способом.",
          )
        }}
      </p>
      <p v-else-if="lead" class="font-body text-body text-ink">{{ lead }}</p>
    </div>

    <!-- `bare` (the trial-lesson embed) gives every field its own glass
    panel, same as the rest of the site's white-on-glass containers —
    vertical padding only, so the panels' edges stay flush with the form's
    own (no side inset, per the caller's layout). Non-bare callers already
    sit on a solid white card, so the wrapper is a plain unstyled div there. -->
    <form v-if="status !== 'success'" class="flex flex-col gap-16" novalidate @submit="onSubmit">
      <div
        :class="
          bare ? 'rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150' : ''
        "
      >
        <UiInput
          name="name"
          :label="text('apply_form_name_label', 'Имя')"
          required
          autocomplete="name"
          :placeholder="text('apply_form_name_placeholder', 'Как вас зовут')"
        />
      </div>
      <div
        :class="
          bare ? 'rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150' : ''
        "
      >
        <UiPhoneInput
          name="phone"
          :label="text('apply_form_phone_label', 'Номер телефона')"
          required
        />
      </div>
      <div
        :class="
          bare ? 'rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150' : ''
        "
      >
        <UiInput
          name="email"
          :label="text('apply_form_email_label', 'Почта')"
          type="email"
          autocomplete="email"
          :placeholder="text('apply_form_email_placeholder', 'you@example.com')"
        />
      </div>
      <!-- has-[...]:z-30 — `backdrop-blur` makes every one of these wrapper
      divs its own stacking context, so the open dropdown (z-20, absolute,
      positioned relative to UiSelect's own root INSIDE this wrapper) is
      confined to it: a later sibling wrapper (its own stacking context)
      then paints over the whole thing regardless of that internal z-20.
      Bumping the wrapper itself above its siblings while open fixes it. -->
      <div
        :class="
          bare
            ? 'rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150 has-[[aria-expanded=true]]:relative has-[[aria-expanded=true]]:z-30'
            : ''
        "
      >
        <UiSelect
          name="contactMethod"
          :label="text('apply_form_contact_method_label', 'Как с вами связаться?')"
          required
          :options="contactMethodOptions"
        />
      </div>
      <div
        :class="
          bare
            ? 'rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150 has-[[aria-expanded=true]]:relative has-[[aria-expanded=true]]:z-30'
            : ''
        "
      >
        <UiSelect
          name="source"
          :label="text('apply_form_source_label', 'Как вы о нас узнали?')"
          required
          :options="sourceOptions"
        />
      </div>

      <div class="rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150">
        <UiCheckbox name="consent">
          {{ text("apply_form_consent_prefix", "Отправляя форму, вы соглашаетесь с") }}
          <NuxtLink to="/privacy" class="text-primary underline">{{
            text("apply_form_consent_link_text", "политикой обработки персональных данных")
          }}</NuxtLink>
          {{
            text(
              "apply_form_consent_suffix",
              "и даёте согласие на обработку указанных персональных данных.",
            )
          }}
        </UiCheckbox>
      </div>

      <p v-if="status === 'error'" class="font-body text-body text-primary" role="alert">
        {{ submitError }}
      </p>

      <UiButton type="submit" block :loading="isSubmitting" :disabled="isSubmitting">
        {{
          context === "trial_lesson"
            ? text("apply_form_submit_trial", "Записаться на занятие")
            : text("apply_form_submit_default", "Отправить заявку")
        }}
      </UiButton>
    </form>
  </div>
</template>
