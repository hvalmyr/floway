/**
 * Works around a confirmed bug in @nuxt/image's IPX route under this app's
 * production build (Nitro's "bun" preset): the route internally does set
 * Content-Type and Cache-Control on every response (ipx's own
 * createIPXHandler, via h3's fromNodeMiddleware wrapping srvx's Node
 * compat shim), but neither header reaches the client — confirmed with the
 * real built .output image running standalone, not just on the deployed
 * site, so this isn't a deploy/config issue. Every /_ipx/** response comes
 * back as Content-Type: application/octet-stream with no Cache-Control at
 * all, on every request, regardless of IPX_MAX_AGE/IPX_HTTP_MAX_AGE (which
 * only affect what IPX would set if the header write path worked).
 *
 * This hook runs at the outer Nitro app level, right before any response
 * is sent, so it reliably overrides whatever the broken inner path did.
 *
 * Cache-Control is unconditional: every uploaded image is immutable
 * (UUID-keyed, never edited in place — see upload_handler.go), and IPX's
 * own id encodes every relevant modifier (width/quality/format) into the
 * URL itself, so any two requests that could produce different bytes are
 * already different URLs — nothing here needs revalidation.
 *
 * Content-Type is derived from the `f_<format>` modifier IPX encodes into
 * the path (e.g. `/_ipx/f_avif&q_55&w_576/...`) — every image in this app
 * requests an explicit format (see resolveOptimizedMediaUrl call sites:
 * NuxtImg's `format` prop, NuxtPicture's per-source format, UiHeroPicture's
 * modifiers), so there's always one to read.
 */
export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook("beforeResponse", (event) => {
    if (!event.path.startsWith("/_ipx/")) return;

    event.node.res.setHeader("cache-control", "public, max-age=31536000, immutable");

    const match = event.path.match(/(?:^|[&/])f(?:ormat)?_([a-z0-9]+)/i);
    if (match) {
      const format = match[1]!.toLowerCase() === "jpg" ? "jpeg" : match[1]!.toLowerCase();
      event.node.res.setHeader("content-type", `image/${format}`);
    }
  });
});
