<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { useDebounceFn } from "@vueuse/core";
import { useForm } from "vee-validate";
import { applyFormSchema, applyFormSchemaWithObjectType } from "~/lib/validation/applyForm";
import type { ContactMethod, LeadFormat, LeadSource, ObjectType } from "~/types/api";

/**
 * Decor site's lead-capture form (single request type — see
 * model.LeadRequestType.decor on the backend) — embedded on the home page
 * and every landing page. Mirrors frontend/app/components/sections/
 * ApplyForm.vue's structure (draft persistence, thank-you redirect,
 * page_content-editable option labels) but replaces the school's course/
 * masterclass context selection with the decor-specific fields the TZ
 * requires (п. 8.1): object type and format, plus UTM/yclid ad-tracking
 * params captured from the URL (п. 9).
 *
 * `objectTypeId` presets and hides the object-type field when embedded on
 * that type's own landing page ("на посадочной странице нужный тип выбран
 * автоматически") — omit it (e.g. on the home page) to show a select
 * populated from the shared object-types dictionary instead.
 *
 * @example
 * <ApplyForm title="Оставить заявку" />
 * <ApplyForm :object-type-id="landingPage.objectTypeId" title="Оставить заявку на оформление дома" bare />
 */
const props = withDefaults(
  defineProps<{
    objectTypeId?: number;
    title?: string;
    lead?: string;
    bare?: boolean;
  }>(),
  { title: "Оставить заявку", lead: "", objectTypeId: undefined, bare: false },
);

const api = useApi();
const route = useRoute();
const { text } = await usePageContent();

const objectTypes = ref<ObjectType[]>([]);
if (!props.objectTypeId) {
  const { data } = await useAsyncData("apply-form-object-types", () => api.getObjectTypes());
  objectTypes.value = (data.value ?? []).filter((t) => t.visible);
}

const { handleSubmit, isSubmitting, values, setValues } = useForm({
  validationSchema: toTypedSchema(
    props.objectTypeId ? applyFormSchema : applyFormSchemaWithObjectType,
  ),
  initialValues: { name: "", phone: "", email: "", consent: false },
});

// Persist in-progress input (not `consent`) so a reload doesn't throw away
// what the visitor already typed — same reasoning as the school's ApplyForm.
// Keyed by objectTypeId: a visitor filling in the homepage's generic form
// (no preset type) gets a separate draft from one already on a specific
// landing page.
const draftStorageKey = `apply-form-draft:${props.objectTypeId ?? "any"}`;

function readDraft() {
  if (!import.meta.client) return null;
  try {
    const raw = localStorage.getItem(draftStorageKey);
    return raw ? (JSON.parse(raw) as Partial<Record<string, string>>) : null;
  } catch {
    return null;
  }
}

function clearDraft() {
  if (!import.meta.client) return;
  localStorage.removeItem(draftStorageKey);
}

onMounted(() => {
  const draft = readDraft();
  if (draft) setValues(draft, false);
});

const saveDraft = useDebounceFn(() => {
  if (!import.meta.client) return;
  const { name, phone, email, contactMethod, source, format, objectTypeId } = values;
  localStorage.setItem(
    draftStorageKey,
    JSON.stringify({ name, phone, email, contactMethod, source, format, objectTypeId }),
  );
}, 300);

watch(values, saveDraft, { deep: true });

const status = ref<"idle" | "error">("idle");
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

const formatOptions = computed(() => [
  { value: "season", label: text("apply_form_format_season", "На сезон") },
  { value: "event", label: text("apply_form_format_event", "На праздник") },
]);

const objectTypeOptions = computed(() =>
  objectTypes.value.map((t) => ({ value: String(t.id), label: t.name })),
);

// UTM/yclid — carried straight from the landing page's query string into the
// lead (п. 9 ТЗ), never shown as form fields, just captured invisibly at
// submit time.
const UTM_KEYS = ["utm_source", "utm_medium", "utm_campaign", "utm_content", "utm_term"] as const;
function queryParam(key: string): string | undefined {
  const value = route.query[key];
  return typeof value === "string" ? value : undefined;
}

const onSubmit = handleSubmit(async (values) => {
  status.value = "idle";
  submitError.value = "";
  try {
    await api.submitApplication({
      name: values.name,
      phone: values.phone,
      email: values.email || undefined,
      contactMethod: values.contactMethod as ContactMethod,
      source: values.source as LeadSource,
      requestType: "decor",
      objectTypeId: props.objectTypeId ?? Number(values.objectTypeId),
      format: values.format as LeadFormat,
      utmSource: queryParam("utm_source"),
      utmMedium: queryParam("utm_medium"),
      utmCampaign: queryParam("utm_campaign"),
      utmContent: queryParam("utm_content"),
      utmTerm: queryParam("utm_term"),
      yclid: queryParam("yclid"),
    });
    clearDraft();
    await navigateTo("/thank-you/decor");
  } catch (err) {
    status.value = "error";
    submitError.value =
      (err as { message?: string } | undefined)?.message ?? "Не удалось отправить заявку";
  }
});
</script>

<template>
  <div :class="bare ? 'w-full' : 'rounded-[30px] bg-white p-24 sm:p-48 lg:p-64'">
    <div v-if="title || lead" class="mb-24 flex flex-col gap-12">
      <h2 v-if="title" class="font-display text-h2 text-ink">{{ title }}</h2>
      <p v-if="lead" class="font-body text-body text-ink">{{ lead }}</p>
    </div>

    <form class="flex flex-col gap-16" novalidate @submit="onSubmit">
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
      <div
        v-if="!objectTypeId"
        :class="
          bare
            ? 'rounded-md bg-white/55 px-16 py-16 backdrop-blur backdrop-saturate-150 has-[[aria-expanded=true]]:relative has-[[aria-expanded=true]]:z-30'
            : ''
        "
      >
        <UiSelect
          name="objectTypeId"
          :label="text('apply_form_object_type_label', 'Тип объекта')"
          required
          :options="objectTypeOptions"
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
          name="format"
          :label="text('apply_form_format_label', 'Формат')"
          required
          :options="formatOptions"
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
          {{ text("apply_form_consent_middle", "и") }}
          <NuxtLink to="/pd-consent" class="text-primary underline">{{
            text("apply_form_consent_link2_text", "даёте согласие на обработку персональных данных")
          }}</NuxtLink
          >{{ text("apply_form_consent_suffix", ".") }}
        </UiCheckbox>
      </div>

      <p v-if="status === 'error'" class="font-body text-body text-primary" role="alert">
        {{ submitError }}
      </p>

      <UiButton type="submit" block :loading="isSubmitting" :disabled="isSubmitting">
        {{ text("apply_form_submit_default", "Отправить заявку") }}
      </UiButton>
    </form>
  </div>
</template>
