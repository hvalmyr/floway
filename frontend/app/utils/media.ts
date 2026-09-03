/**
 * Uploaded images are stored/returned as relative paths (`/uploads/<key>`,
 * served by the backend). Legacy manually-pasted URLs may still be full
 * external URLs — passed through unchanged. Only the browser needs this
 * (backend runs on a different origin/port in dev); SSR pages don't render
 * admin-uploaded images.
 */
export function resolveMediaUrl(path: string): string {
  if (!path || !path.startsWith("/")) return path;
  const config = useRuntimeConfig();
  return `${config.public.apiBase}${path}`;
}

/**
 * Same idea as resolveMediaUrl, but for images passed into NuxtImg/
 * NuxtPicture. Those fetch the original themselves, server-side (either
 * during SSR or when the browser later hits the /_ipx/** URL they
 * generate) — so, unlike a plain <img src>, the URL has to stay reachable
 * from the Nuxt server process itself, in every environment. See
 * mediaOptimizeBase's comment in nuxt.config.ts.
 */
export function resolveOptimizedMediaUrl(path: string): string {
  if (!path || !path.startsWith("/")) return path;
  const config = useRuntimeConfig();
  return `${config.public.mediaOptimizeBase}${path}`;
}
