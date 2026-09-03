import { createHash } from "node:crypto";
import type { ServerResponse } from "node:http";

/**
 * @nuxt/image's IPX route has no caching of its own — it re-fetches the
 * original and re-encodes from scratch on every single request. Measured
 * live on the actual VPS: avif alone takes ~4s per request, every request,
 * forever (confirmed by requesting the exact same URL three times in a row
 * and getting ~4s every time) — that's what was actually behind slow image
 * loads and a real chunk of the TBT/LCP numbers, not just the 3D
 * background.
 *
 * Tried Nitro's built-in routeRules `cache` option first; that broke
 * createIPXNodeHandler's own URL parsing outright — confirmed against the
 * real built image, "Invalid URL: http://:80/_ipx/..." on literally the
 * first request, because Nitro's cache wrapper re-dispatches the request
 * through a path that doesn't carry a proper Host.
 *
 * Tried capturing the body via a `beforeResponse` hook next — also
 * doesn't work: confirmed live that `response.body` is plain `null` at
 * that point (not merely falsy — actually null), because
 * createIPXNodeHandler's fromNodeMiddleware wraps a Node-style
 * (req, res) => {} handler that calls res.end(buffer) *directly* on the
 * raw stream, bypassing h3's own body-capture mechanism entirely —
 * there's nothing in a beforeResponse hook's `response` argument to read.
 *
 * So this wraps the raw Node response stream itself, at the lowest level
 * that's guaranteed to see the real bytes regardless of how the handler
 * upstream chooses to write them: patches res.write/res.end for the
 * duration of one request, accumulates every chunk, and stores the result
 * once the response actually finishes. Content-Type is parsed straight
 * from the URL's `f_<format>` modifier (same regex as fix-ipx-headers.ts)
 * rather than read back off the response headers — IPX's own Content-Type
 * header is unreliable (see that file), and this way caching doesn't
 * depend on hook execution order between the two plugins either.
 */
export const IPX_CACHE_STORAGE = "cache:ipx-manual";
export const IPX_CACHE_HIT_HEADER = "x-ipx-cache";

export function ipxCacheKey(path: string): string {
  return createHash("sha256").update(path).digest("hex");
}

export function contentTypeFromIpxPath(path: string): string | undefined {
  const match = path.match(/(?:^|[&/])f(?:ormat)?_([a-z0-9]+)/i);
  if (!match) return undefined;
  const format = match[1]!.toLowerCase() === "jpg" ? "jpeg" : match[1]!.toLowerCase();
  return `image/${format}`;
}

export default defineEventHandler(async (event) => {
  if (!event.path.startsWith("/_ipx/")) return;

  const contentType = contentTypeFromIpxPath(event.path);
  const key = ipxCacheKey(event.path);
  const storage = useStorage(IPX_CACHE_STORAGE);
  const cached = await storage.getItem<string>(key);
  if (cached) {
    if (contentType) setResponseHeader(event, "content-type", contentType);
    setResponseHeader(event, "cache-control", "public, max-age=31536000, immutable");
    setResponseHeader(event, IPX_CACHE_HIT_HEADER, "HIT");
    return Buffer.from(cached, "base64");
  }

  // Cache miss — let the real (slow) handler run, but capture every byte
  // it writes so the *next* request can skip straight to the block above.
  if (!contentType) return; // nothing to key the cache's format on
  const res = event.node.res;
  const chunks: Buffer[] = [];
  const originalWrite = res.write.bind(res);
  const originalEnd = res.end.bind(res);

  res.write = ((chunk: unknown, ...args: unknown[]) => {
    if (chunk) chunks.push(Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk as string));
    // @ts-expect-error — passing through whatever encoding/callback args
    // the caller used; we only care about observing the bytes, not
    // re-typing Node's overloaded write() signature.
    return originalWrite(chunk, ...args);
  }) as ServerResponse["write"];

  res.end = ((chunk?: unknown, ...args: unknown[]) => {
    if (chunk) chunks.push(Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk as string));
    if (res.statusCode === 200 && chunks.length > 0) {
      void storage.setItem(key, Buffer.concat(chunks).toString("base64"));
    }
    // @ts-expect-error — same as write() above.
    return originalEnd(chunk, ...args);
  }) as ServerResponse["end"];
});
