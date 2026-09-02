export type CookieConsentStatus = "accepted";

const STORAGE_KEY = "floway-cookie-consent";

/**
 * Cookie/analytics consent, persisted in a cookie (not localStorage) so the
 * server can read it on the very first render — the banner's `v-if` is
 * correct in the SSR-rendered HTML itself, with no client-side "hide it
 * after mount" step that would otherwise flash the banner on every reload
 * for a visitor who already answered. Accept-only: not accepting simply
 * leaves the cookie unset, which keeps analytics off and the banner
 * showing — there's no separate "declined" state to track.
 */
export function useCookieConsent() {
  const consent = useCookie<CookieConsentStatus | null>(STORAGE_KEY, {
    default: () => null,
    maxAge: 60 * 60 * 24 * 365,
    sameSite: "lax",
  });

  function accept() {
    consent.value = "accepted";
  }

  return { status: consent, accept };
}
