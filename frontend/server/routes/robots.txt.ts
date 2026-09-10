// Replaces the old static public/robots.txt so the Sitemap: line can carry
// an absolute, environment-correct URL — built from the incoming request's
// own Host header (see sitemap.xml.ts's comment on why that's used instead
// of config.public.apiBase) instead of a value baked in at build time.
export default defineEventHandler((event) => {
  const siteOrigin = getRequestURL(event).origin;

  setHeader(event, "Content-Type", "text/plain; charset=UTF-8");
  return `User-Agent: *\nDisallow:\n\nSitemap: ${siteOrigin}/sitemap.xml\n`;
});
