/**
 * Counterpart to server/plugins/strip-duplicate-css-link.ts: that plugin
 * turns Nuxt's own stylesheet <link> into a non-blocking
 * rel="preload" as="style" so it doesn't delay first paint (the page's
 * inline critical CSS already covers that — see the `inlineStyles`
 * comment in nuxt.config.ts). A preload is only ever fetched, never
 * applied as CSS, so this flips it to a real stylesheet as soon as the
 * client runs. By then the browser's initial-HTML parse is long past, so
 * doing it here is non-blocking exactly like any other script-inserted
 * stylesheet — no CLS risk, since the page is already laid out correctly
 * from the inlined critical CSS and these rules are identical, just
 * arriving a moment later.
 *
 * Without this, the full site CSS — which covers components the current
 * page's SSR pass never rendered, and so never inlined (e.g.
 * RichTextContent.vue, only ever mounted from /blog/[slug]) — never loads
 * at all, since Nuxt Router navigations don't reload the document and so
 * never get another SSR pass either. Those components then render
 * unstyled (blog post headings/quotes/images at browser-default size)
 * after any client-side navigation, correcting itself only on a hard
 * reload.
 */
export default defineNuxtPlugin(() => {
  document.querySelectorAll<HTMLLinkElement>('link[rel="preload"][as="style"]').forEach((link) => {
    link.rel = "stylesheet";
  });
});
