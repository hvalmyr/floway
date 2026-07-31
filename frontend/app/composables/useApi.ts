/**
 * Configured $fetch instance for the Go backend. `credentials: 'include'` is
 * required so the browser sends the httpOnly admin session cookie on
 * cross-origin requests (frontend and backend run on different ports).
 */
export function useApi() {
  const config = useRuntimeConfig();

  return $fetch.create({
    baseURL: config.public.apiBase,
    credentials: "include",
  });
}
