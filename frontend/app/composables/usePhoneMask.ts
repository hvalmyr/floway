import { AsYouType, getExampleNumber } from "libphonenumber-js/min";
import type { CountryCode } from "libphonenumber-js";
// @ts-expect-error -- untyped JSON metadata, shaped for getExampleNumber's `examples` param
import examples from "libphonenumber-js/examples.mobile.json";
import { ref } from "vue";

/**
 * Formats raw input into a live international phone mask. Students come
 * from Russia, Belarus, and further afield (Europe, South America), so this
 * isn't Russia-only: type a bare number (no "+") and it's assumed Russian —
 * the "8"/"7" trunk prefix is stripped and +7 is applied, matching the old
 * behavior and the common local habit of dialing with a leading 8 — but
 * typing an explicit "+" followed by any country's calling code (+375, +34,
 * +55, ...) reformats to that country's own pattern instead.
 */
export function formatPhoneMask(rawValue: string): string {
  let digits = rawValue.replace(/[^\d+]/g, "");

  if (digits.startsWith("+")) {
    return new AsYouType().input(digits);
  }

  if (digits.startsWith("8") || digits.startsWith("7")) {
    digits = digits.slice(1);
  }
  return new AsYouType().input(`+7${digits}`);
}

/**
 * Same detection rule formatPhoneMask() uses internally (bare digits assume
 * Russia; a leading "+" is parsed for its actual calling code) — exposed
 * separately so UiPhoneInput can build a matching placeholder template
 * without re-deriving the logic. Returns undefined while a "+"-prefixed
 * calling code is still ambiguous (e.g. "+3" alone could be several
 * countries) — there's nothing sensible to show yet at that point. The one
 * exception is +7 itself: Russia and Kazakhstan share that calling code, so
 * getCountry() stays undefined until enough of the subscriber number is
 * typed to tell them apart — default that window to Russia (the school's
 * actual market) instead of falling back to the generic template.
 */
export function detectPhoneCountry(rawValue: string): CountryCode | undefined {
  const digits = rawValue.replace(/[^\d+]/g, "");
  if (!digits) return undefined;
  if (!digits.startsWith("+")) return "RU";
  const formatter = new AsYouType();
  formatter.input(digits);
  const country = formatter.getCountry();
  if (country) return country;
  return formatter.getCallingCode() === "7" ? "RU" : undefined;
}

const GENERIC_PHONE_TEMPLATE = "+000 000 000 000";

/**
 * A same-punctuation "+X XXX XXX XX XX"-shaped template for `country`, built
 * from a real example number so the grouping always matches what
 * formatPhoneMask() actually produces for that country — the calling code
 * digits stay real, everything after is zeroed. Falls back to a generic
 * placeholder before a country is known.
 */
export function phoneTemplateFor(country: CountryCode | undefined): string {
  if (!country) return GENERIC_PHONE_TEMPLATE;
  const example = getExampleNumber(country, examples);
  if (!example) return GENERIC_PHONE_TEMPLATE;
  const formatted = example.formatInternational();
  const callingCodeLength = String(example.countryCallingCode).length;
  let digitsSeen = 0;
  return formatted.replace(/\d/g, (digit) => (++digitsSeen <= callingCodeLength ? digit : "0"));
}

export function usePhoneMask() {
  const phone = ref("+7 ");

  function onInput(event: Event) {
    const target = event.target as HTMLInputElement;
    phone.value = formatPhoneMask(target.value);
  }

  return { phone, onInput };
}
