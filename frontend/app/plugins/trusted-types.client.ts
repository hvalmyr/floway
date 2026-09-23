import { sanitizeRichTextHtml } from "~/lib/richTextSanitize";

/**
 * Registers the "default" Trusted Types policy — the browser's reserved
 * catch-all name (also allow-listed in the CSP `trusted-types` directive,
 * see nuxt.config.ts): any raw string assigned to a restricted DOM sink
 * (.innerHTML, etc.) that isn't already going through an explicit policy
 * gets piped through this one automatically, with no call-site changes
 * needed. Vue's own v-html bindings go through their own built-in "vue"
 * policy instead (a passthrough — see @vue/runtime-dom's
 * `unsafeToTrustedHTML`), so in practice this only catches imperative
 * innerHTML writes like AdminRichTextEditor.vue's editor-sync assignment.
 *
 * No-op if the browser doesn't support Trusted Types (current Firefox/
 * Safari) — require-trusted-types-for is spec'd to simply not apply there.
 */
export default defineNuxtPlugin(() => {
  if (typeof window === "undefined" || !window.trustedTypes) return;
  try {
    window.trustedTypes.createPolicy("default", {
      createHTML: (input: string) => sanitizeRichTextHtml(input),
      // useYandexMetrika.ts sets a <script>'s textContent to a hardcoded,
      // first-party snippet (only the numeric counter ID is interpolated —
      // no admin/DB content reaches this sink), so a passthrough is enough;
      // createScriptURL is defined for the same reason no other sink here
      // needs anything stricter, not because anything currently uses it.
      createScript: (input: string) => input,
      createScriptURL: (input: string) => input,
    });
  } catch {
    // A given policy name can only be created once — don't crash the app
    // if something (HMR in dev, a future third-party script) beat us to it.
  }
});
