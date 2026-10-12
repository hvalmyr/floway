/**
 * The site's whole CSS is one ~9KB file (see vite.build.cssCodeSplit in
 * nuxt.config.ts). `features.inlineStyles` (see its comment in
 * nuxt.config.ts) inlines each page's CSS as a <style> tag so the very
 * first paint never blocks on a stylesheet request — but it only inlines
 * the rules for components actually rendered in *that page's* SSR tree,
 * not the whole bundle. A component that first mounts after a client-side
 * route change (e.g. RichTextContent.vue, only ever used on
 * /blog/[slug], never on /blog itself) has no inlined rules on the page
 * it's mounted from, and since Vue Router never reloads the document, no
 * later SSR pass ever inlines them either — its headings/quotes/etc stay
 * unstyled until a full page reload.
 *
 * This plugin used to delete Nuxt's own external `<link rel="stylesheet">`
 * outright as "redundant" (reasoning: every route's styles are already
 * inlined, and a client-side nav never needs fresh CSS) — both false, as
 * above, and exactly what caused that bug: deleting it removed the one
 * copy of the CSS that actually covers every component, not just the
 * current page's.
 *
 * So the link now stays, but as a non-blocking `rel="preload"` instead of
 * a normal blocking `rel="stylesheet"` — the inlined critical CSS already
 * makes first paint correct, so this fetch doesn't need to block
 * anything. app/plugins/css-preload-swap.client.ts flips it to a real
 * stylesheet once the client runs, which is what actually makes it apply.
 * A `<noscript>` fallback covers requests with JS disabled, where that
 * plugin never runs.
 */
const STYLESHEET_RE = /<link\b[^>]*rel="stylesheet"[^>]*href="(\/_nuxt\/[^"]+\.css)"[^>]*>/g;

export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook("render:html", (html) => {
    html.head = html.head.map((chunk) =>
      typeof chunk === "string"
        ? chunk.replace(
            STYLESHEET_RE,
            (_match, href) =>
              `<link rel="preload" as="style" href="${href}">` +
              `<noscript><link rel="stylesheet" href="${href}"></noscript>`,
          )
        : chunk,
    );
  });
});
