/**
 * The site's whole CSS is one ~9KB file (see vite.build.cssCodeSplit in
 * nuxt.config.ts), but Nuxt still ships it as a plain
 * `<link rel="stylesheet">` — render-blocking, confirmed on PageSpeed
 * Insights to cost ~950ms (a full extra round-trip before first paint, on
 * every page). Nuxt's own `features.inlineStyles` doesn't fix this: it
 * inlines the CSS as an extra `<style>` block but keeps this exact `<link>`
 * too (confirmed live in the built output — `getRequestDependencies` always
 * emits it regardless of inlining), so it only adds bytes without removing
 * the block.
 *
 * This rewrites that `<link>` into the classic non-blocking-CSS pattern:
 * `rel="preload"` (fetches without blocking render), swapped to
 * `rel="stylesheet"` once loaded. The swap runs from a <script> rather than
 * an inline `onload="..."` attribute — this app's CSP sets
 * `script-src-attr: 'none'` (see nuxt.config.ts), which blocks inline event
 * handler attributes outright, nonce or not. `<noscript>` keeps the page
 * styled with JS disabled.
 *
 * The injected `<script>` carries no `nonce` — it's allowed via a fixed
 * `'sha256-...'` hash in script-src instead (see nuxt.config.ts). Two
 * reasons: nuxt-security's own `render:html` hook unconditionally stamps
 * *every* `<script>` tag it finds with the request's nonce (see its
 * `40-cspSsrNonce.js`, `SCRIPT_RE` branch — no "already has one" guard,
 * unlike its `<link>` handling), so a nonce added here too produced a
 * literal `nonce="x" nonce="x"` duplicate that Chrome's CSP check silently
 * treated as *no* nonce at all, blocking the script outright even though
 * `script.nonce` still read the right value (confirmed live). And even
 * without that clash, a per-request nonce is the wrong tool for a script
 * whose text never changes — the hash is checked independently of
 * whatever nonce ends up on the tag, so it isn't sensitive to nuxt-
 * security's behavior or hook ordering at all. The script text must stay
 * byte-for-byte identical to the hash comment in nuxt.config.ts or the
 * hash stops matching and the page silently loses its CSS.
 */
const STYLESHEET_RE =
  /<link\b([^>]*?)rel="stylesheet"([^>]*?)href="(\/_nuxt\/[^"]+\.css)"([^>]*?)>/g;

export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook("render:html", (html) => {
    html.head = html.head.map((chunk) =>
      typeof chunk === "string"
        ? chunk.replace(STYLESHEET_RE, (_match, pre, mid, href, post) => {
            const attrs = `${pre}${mid}${post}`.trim();
            const attrsStr = attrs ? ` ${attrs}` : "";
            return (
              `<link rel="preload" as="style" href="${href}"${attrsStr}>` +
              `<script>document.currentScript.previousElementSibling.addEventListener("load",function(){this.rel="stylesheet"})</script>` +
              `<noscript><link rel="stylesheet" href="${href}"${attrsStr}></noscript>`
            );
          })
        : chunk,
    );
  });
});
