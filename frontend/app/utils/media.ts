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
