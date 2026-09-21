/**
 * After an admin upload, fires the same /_ipx/ requests the site's actual
 * <picture> components would make for that image — so IPX's result cache
 * (see server/middleware/ipx-cache.ts) is already warm by the time a real
 * visitor loads the page, instead of the first visitor paying the ~4s
 * re-encode.
 *
 * AdminImageUpload.vue is shared by every admin form (course covers,
 * teacher photos, blog covers, gallery/gift-certificate photos, ...), so it
 * has no idea which component/slot a given upload will end up rendered in.
 * Rather than threading that mapping through every call site, this just
 * warms the union of every *fixed* sizes profile used across the site
 * (UiHeroPicture's internal HERO_SIZES, plus every sizes= passed to
 * UiContentImage) by calling the exact same img.getSizes() those
 * components call — so it can't drift out of sync with them the way a
 * separately-maintained profile table could.
 *
 * PhotoCarousel's thumbnail width is viewport-height-driven and measured
 * off the DOM at runtime rather than a fixed sizes= string (see its own
 * thumbSizesPx comment) — but it now rounds that measurement up to a fixed
 * bucket (see photoCarouselSizes.ts), so the space of widths it can ever
 * request is enumerable too. photoCarouselWarmWidths() is that enumerable
 * range, warmed the same way as everything else here.
 */
const HERO_SIZES = "400:100vw lg:576px";

// One entry per distinct sizes= currently passed to UiContentImage
// (CourseCard, MasterclassCard, index.vue x2, blog/[slug], blog/index).
const CONTENT_SIZES = [
  "400:100vw sm:50vw lg:400px",
  "400:100vw lg:38vw",
  "400:100vw md:50vw",
  "400:100vw md:33vw",
  "400:100vw lg:720px",
  "400:100vw md:50vw lg:33vw",
];

const WARMUP_CONCURRENCY = 3;

async function fetchAllLimited(urls: string[], limit: number) {
  const queue = [...urls];
  async function worker() {
    let url = queue.shift();
    while (url) {
      try {
        await fetch(url);
      } catch {
        // Best-effort — a failed warm-up just means the first real visitor
        // eats the slow path, same as before this existed.
      }
      url = queue.shift();
    }
  }
  await Promise.all(Array.from({ length: limit }, worker));
}

export function useImageCacheWarmup() {
  const img = useImage();
  const quality = useImageQuality();

  function collect(urls: Set<string>, sizes: string, format: string, path: string, q: number) {
    const { srcset } = img.getSizes(path, { sizes, modifiers: { format, quality: q } });
    for (const entry of srcset.split(",")) {
      const url = entry.trim().split(" ")[0];
      if (url) urls.add(url);
    }
  }

  /**
   * `path` must be the same resolveOptimizedMediaUrl() result the display
   * components receive as `src` — otherwise the generated /_ipx/ paths
   * (and their cache keys) won't match what a real page actually requests.
   */
  function warm(path: string) {
    const urls = new Set<string>();
    const q = quality.value;

    collect(urls, HERO_SIZES, "avif", path, q.avif);
    collect(urls, HERO_SIZES, "webp", path, q.mobile);
    collect(urls, HERO_SIZES, "webp", path, q.desktop);
    collect(urls, HERO_SIZES, "jpeg", path, q.mobile);
    collect(urls, HERO_SIZES, "jpeg", path, q.desktop);

    for (const sizes of CONTENT_SIZES) {
      collect(urls, sizes, "avif", path, q.avif);
      collect(urls, sizes, "webp", path, q.mobile);
      collect(urls, sizes, "webp", path, q.desktop);
    }

    for (const width of photoCarouselWarmWidths()) {
      const sizes = `${width}px`;
      collect(urls, sizes, "avif", path, q.avif);
      collect(urls, sizes, "webp", path, q.mobile);
      collect(urls, sizes, "webp", path, q.desktop);
    }

    void fetchAllLimited([...urls], WARMUP_CONCURRENCY);
  }

  return { warm };
}
