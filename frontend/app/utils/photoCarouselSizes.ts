/**
 * PhotoCarousel measures its thumbnail width directly off the DOM (see its
 * own thumbSizesPx comment) because its box is viewport-*height* driven
 * (h-[40vh] aspect-[3/4]), a unit @nuxt/image's sizes DSL can't express.
 * Left unbucketed, that measurement is effectively continuous — every
 * visitor with a slightly different viewport height requests its own
 * distinct IPX width, so the result cache (server/middleware/ipx-cache.ts)
 * almost never gets reused across visitors, and useImageCacheWarmup has no
 * finite list of widths to pre-warm ahead of time.
 *
 * Rounding the measured width UP to the nearest WIDTH_BUCKET_PX collapses
 * that near-continuous range onto a small fixed set of widths — visually
 * unnoticeable (at most WIDTH_BUCKET_PX of extra size, always rounded up
 * so the served image is never smaller than the rendered box) — which both
 * makes organic cache reuse far more likely and gives useImageCacheWarmup
 * an enumerable range to warm.
 */
export const WIDTH_BUCKET_PX = 20;
// Covers roughly 500px-1600px of real viewport height (box width is ~0.3 *
// viewport height). Taller screens still render correctly — their bucket
// just falls outside the pre-warmed range and relies on organic caching
// after that bucket's first real visitor, same as before this existed.
export const MIN_WARM_WIDTH_PX = 160;
export const MAX_WARM_WIDTH_PX = 480;

export function bucketPhotoCarouselWidth(px: number): number {
  return Math.ceil(px / WIDTH_BUCKET_PX) * WIDTH_BUCKET_PX;
}

export function photoCarouselWarmWidths(): number[] {
  const widths: number[] = [];
  for (let w = MIN_WARM_WIDTH_PX; w <= MAX_WARM_WIDTH_PX; w += WIDTH_BUCKET_PX) widths.push(w);
  return widths;
}
