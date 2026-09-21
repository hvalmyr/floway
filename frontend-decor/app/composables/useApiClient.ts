/**
 * Configured $fetch instance for the Go backend. `credentials: 'include'` is
 * required so the browser sends the httpOnly admin session cookie on
 * cross-origin requests (frontend and backend run on different ports).
 *
 * Uses a different base URL on the server than in the browser: SSR runs
 * inside the frontend container, so it must reach the backend over the
 * internal docker network (runtimeConfig.apiBaseInternal, e.g.
 * http://backend:8080) — the public URL's hostname/domain isn't reachable
 * (or means something else entirely) from inside that container. The
 * browser always uses the public URL (runtimeConfig.public.apiBase).
 *
 * This is the low-level client used by the admin panel, which needs generic
 * verb/body control (POST/PUT/PATCH/DELETE) that a fixed set of named
 * methods can't express. Public-facing pages should use useApi() instead —
 * it wraps this same client with named methods, error normalization, and a
 * mocks fallback.
 */
export function useApiClient() {
  const config = useRuntimeConfig();
  const baseURL = import.meta.server ? config.apiBaseInternal : config.public.apiBase;

  return $fetch.create({
    baseURL,
    credentials: "include",
  });
}
