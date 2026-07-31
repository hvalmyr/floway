/**
 * Configured $fetch instance for the Go backend. `credentials: 'include'` is
 * required so the browser sends the httpOnly admin session cookie on
 * cross-origin requests (frontend and backend run on different ports).
 *
 * This is the low-level client used by the admin panel, which needs generic
 * verb/body control (POST/PUT/PATCH/DELETE) that a fixed set of named
 * methods can't express. Public-facing pages should use useApi() instead —
 * it wraps this same client with named methods, error normalization, and a
 * mocks fallback.
 */
export function useApiClient() {
  const config = useRuntimeConfig();

  return $fetch.create({
    baseURL: config.public.apiBase,
    credentials: "include",
  });
}
