// Dynamic sitemap — decor-site equivalent of frontend/server/routes/sitemap.xml.ts.
// Lists static pages plus every visible landing page, so a new object-type
// page (created from the admin without a deploy, п. 7.5 ТЗ) shows up here
// too. Backend fetch uses apiBaseInternal (server-only docker-network
// address, same reasoning as useApiClient.ts) — this route only ever runs
// server-side. Absolute URLs are built from the incoming request's own Host
// header (getRequestURL), not config.public.apiBase — see the school
// sitemap's own comment for why (different origins in local/dev-docker).
const STATIC_PATHS = [
  "/",
  "/portfolio",
  "/how-it-works",
  "/contacts",
  "/privacy",
  "/terms",
  "/cookie-policy",
  "/pd-consent",
  "/legal-info",
];

interface LandingPageSummary {
  slug: string;
}

function xmlEscape(value: string): string {
  return value.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const siteOrigin = getRequestURL(event).origin;

  let landingSlugs: string[] = [];
  try {
    const pages = await $fetch<LandingPageSummary[]>("/api/v1/landing-pages/visible", {
      baseURL: config.apiBaseInternal,
    });
    landingSlugs = pages.map((page) => page.slug);
  } catch {
    // Backend unreachable — сайтмап всё равно отдаём со статическими страницами.
  }

  const urls: { loc: string }[] = [
    ...STATIC_PATHS.map((path) => ({ loc: `${siteOrigin}${path}` })),
    ...landingSlugs.map((slug) => ({ loc: `${siteOrigin}/${slug}` })),
  ];

  const body = urls.map((u) => `  <url>\n    <loc>${xmlEscape(u.loc)}</loc>\n  </url>`).join("\n");

  setHeader(event, "Content-Type", "application/xml; charset=UTF-8");
  return `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${body}\n</urlset>\n`;
});
