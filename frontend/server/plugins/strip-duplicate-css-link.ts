/**
 * The site's whole CSS is one ~9KB file (see vite.build.cssCodeSplit in
 * nuxt.config.ts), fully inlined into every page's <head> as a <style> tag
 * by the `features.inlineStyles` config (see its comment there for why
 * inlining, not deferring, is the right fix — a deferred stylesheet caused
 * a severe CLS regression on this Tailwind-only-sizing site). Nuxt inlines
 * the styles but does NOT drop the matching external
 * `<link rel="stylesheet">` it would otherwise emit — confirmed live in the
 * built output (`getRequestDependencies` always returns it, inlining or
 * not) — so every page still shipped a real render-blocking network
 * request for content that was already on the page twice over.
 *
 * With `cssCodeSplit: false`, that single global stylesheet is genuinely
 * dead weight once inlined: every route's styles are already in the one
 * inlined chunk from the very first response, and Nuxt Router's
 * client-side navigations never reload the document (so nothing later in
 * the session needs the external file either). This just deletes the
 * link outright instead of letting the browser fetch a second, unneeded
 * copy of bytes already sitting in the page.
 */
const STYLESHEET_RE = /<link\b[^>]*rel="stylesheet"[^>]*href="\/_nuxt\/[^"]+\.css"[^>]*>/g;

export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook("render:html", (html) => {
    html.head = html.head.map((chunk) =>
      typeof chunk === "string" ? chunk.replace(STYLESHEET_RE, "") : chunk,
    );
  });
});
