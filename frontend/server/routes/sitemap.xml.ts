// Dynamic sitemap: lists static pages plus every published blog post and
// visible course, so new content shows up without a redeploy. Backend
// fetches use apiBaseInternal (server-only docker-network address, same
// reasoning as useApiClient.ts) — this route only ever runs server-side.
// Absolute URLs are built from the incoming request's own Host header
// (getRequestURL) rather than config.public.apiBase — the two only happen
// to be the same origin in prod (Caddy serves both the site and /api/* off
// FLOWAY_DOMAIN); in local dev and docker-compose dev the frontend and
// backend sit on different ports, so apiBase would point a search engine
// at the API's own origin instead of the site's.
const STATIC_PATHS = [
  "/",
  "/blog",
  "/masterclasses",
  "/sertifikaty",
  "/contacts",
  "/privacy",
  "/terms",
  "/cookie-policy",
  "/legal-info",
  "/pd-consent",
];

interface BlogPostSummary {
  slug: string;
  status: string;
  updatedAt: string;
}

interface CourseSectionSummary {
  courses: { slug: string; visible: boolean }[];
}

function xmlEscape(value: string): string {
  return value.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const siteOrigin = getRequestURL(event).origin;

  let blogPosts: BlogPostSummary[] = [];
  try {
    const posts = await $fetch<BlogPostSummary[]>("/api/v1/blog-posts", {
      baseURL: config.apiBaseInternal,
      query: { status: "published" },
    });
    blogPosts = posts.filter((post) => post.status === "published");
  } catch {
    // Backend unreachable — сайтмап всё равно отдаём со статическими страницами.
  }

  let courseSlugs: string[] = [];
  try {
    const sections = await $fetch<CourseSectionSummary[]>("/api/v1/course-sections/full", {
      baseURL: config.apiBaseInternal,
    });
    courseSlugs = sections.flatMap((section) =>
      section.courses.filter((course) => course.visible).map((course) => course.slug),
    );
  } catch {
    // Same fallback as above.
  }

  const urls: { loc: string; lastmod?: string }[] = [
    ...STATIC_PATHS.map((path) => ({ loc: `${siteOrigin}${path}` })),
    ...courseSlugs.map((slug) => ({ loc: `${siteOrigin}/courses/${slug}` })),
    ...blogPosts.map((post) => ({
      loc: `${siteOrigin}/blog/${post.slug}`,
      lastmod: post.updatedAt?.slice(0, 10),
    })),
  ];

  const body = urls
    .map(
      (u) =>
        `  <url>\n    <loc>${xmlEscape(u.loc)}</loc>${
          u.lastmod ? `\n    <lastmod>${u.lastmod}</lastmod>` : ""
        }\n  </url>`,
    )
    .join("\n");

  setHeader(event, "Content-Type", "application/xml; charset=UTF-8");
  return `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${body}\n</urlset>\n`;
});
