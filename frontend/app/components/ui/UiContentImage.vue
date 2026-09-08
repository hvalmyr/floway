<script setup lang="ts">
/**
 * Lazy-loaded content photo (course/masterclass/blog covers, gallery
 * thumbnails, teacher photos) — same idea as UiHeroPicture but for
 * below-the-fold images: not eager, no fetchpriority, and split only by
 * webp/avif (no jpeg fallback layer, matching what every one of these call
 * sites already shipped before this component existed).
 *
 * Quality comes from useImageQuality() (admin-editable, see
 * /admin/page-content) instead of being hardcoded per call site. webp gets
 * a real mobile/desktop split (screens are smaller and viewed further
 * away, so mobile tolerates more compression) — avif doesn't, since its
 * quality number is about codec artifacts, not viewport size, and it
 * already needs a much lower value than webp for an equivalent look.
 *
 * Every call site's `class` was written for a bare <img>/<NuxtImg> that WAS
 * the direct flex or grid item (e.g. MasterclassCard's `lg:w-[38%]
 * lg:shrink-0`, index.vue's `order-1`/`md:sticky`/`aspect-[9/16]`) — so the
 * *whole* class, including layout properties that only mean something on
 * the actual flex/grid item, goes on <picture> (inheritAttrs off, bound
 * explicitly below) rather than <img>. Two things this is NOT done, each
 * broken for a different real reason (both confirmed live on index.vue's
 * trial photo, a sticky/ordered grid item):
 *  - `<picture class="contents">` (to make <img> the item instead) reads
 *    right for flexbox, but CSS Grid can straight up fail to place a grid
 *    item behind a `display: contents` ancestor — the photo collapsed to a
 *    0×0 box positioned below the entire grid, not just in the wrong spot.
 *  - Putting the *same* class on <img> too (so it also gets object-cover
 *    etc.) reads right, but <img>'s own `aspect-[9/16]` + `max-h-[80vh]`
 *    then has to resolve against a parent (<picture>) whose width in turn
 *    depends on that same img — real browsers resolve that circular replaced-
 *    element sizing to 0×0 rather than the intended box.
 * <picture>'s box is already exactly right once it alone carries the full
 * class (confirmed: 100% of caller intent — order, sticky, aspect-ratio,
 * max-height — applies correctly to a plain block box with no img inside
 * competing for the same computation). <img> only needs to fill that
 * already-correctly-sized box, no aspect-ratio math of its own: `size-full
 * object-cover`, plus non-`class` attrs (draggable, event handlers) via
 * `restAttrs` so e.g. the carousel's `draggable="false"` still suppresses
 * the native drag ghost on the actual <img>. `overflow-hidden` is forced
 * onto <picture> alongside the caller's class so a caller's `rounded-*`
 * still visually clips the (now full-bleed) <img> the way it clipped a
 * bare <img> directly before this component existed.
 *
 * @example
 * <UiContentImage
 *   :src="resolveOptimizedMediaUrl(coverImage)"
 *   :alt="name"
 *   sizes="400:100vw sm:50vw lg:400px"
 *   class="size-full rounded-sm object-cover"
 * />
 */
defineOptions({ inheritAttrs: false });

const props = defineProps<{
  src: string;
  alt: string;
  sizes: string;
}>();

const MOBILE_MEDIA = "(max-width: 767px)";
const DESKTOP_MEDIA = "(min-width: 768px)";

const img = useImage();
const quality = useImageQuality();

// $attrs.class goes on <picture> only (see the doc comment above) — the
// rest (draggable, event handlers, ...) still belongs on the actual <img>
// too, so it's forwarded there separately with `class` stripped out.
const attrs = useAttrs();
const restAttrs = computed(() => {
  const { class: _callerClass, ...rest } = attrs;
  return rest;
});

const avif = computed(() =>
  img.getSizes(props.src, {
    sizes: props.sizes,
    modifiers: { format: "avif", quality: quality.value.avif },
  }),
);
const webpMobile = computed(() =>
  img.getSizes(props.src, {
    sizes: props.sizes,
    modifiers: { format: "webp", quality: quality.value.mobile },
  }),
);
const webpDesktop = computed(() =>
  img.getSizes(props.src, {
    sizes: props.sizes,
    modifiers: { format: "webp", quality: quality.value.desktop },
  }),
);
</script>

<template>
  <picture v-bind="$attrs" class="overflow-hidden">
    <source type="image/avif" :srcset="avif.srcset" :sizes="avif.sizes" />
    <source
      :media="MOBILE_MEDIA"
      type="image/webp"
      :srcset="webpMobile.srcset"
      :sizes="webpMobile.sizes"
    />
    <source
      :media="DESKTOP_MEDIA"
      type="image/webp"
      :srcset="webpDesktop.srcset"
      :sizes="webpDesktop.sizes"
    />
    <img
      v-bind="restAttrs"
      class="block size-full object-cover"
      :src="webpDesktop.src"
      :srcset="webpDesktop.srcset"
      :sizes="webpDesktop.sizes"
      :alt="alt"
      loading="lazy"
    />
  </picture>
</template>
