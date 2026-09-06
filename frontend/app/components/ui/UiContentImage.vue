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
 * `class` and any other passthrough attrs (draggable, event handlers, …)
 * land on the rendered <img>, not the wrapping <picture> — inheritAttrs is
 * off so e.g. `object-cover` from a caller's `class` actually sizes the
 * photo instead of doing nothing on a <picture> element.
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
  <picture>
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
      v-bind="$attrs"
      :src="webpDesktop.src"
      :srcset="webpDesktop.srcset"
      :sizes="webpDesktop.sizes"
      :alt="alt"
      loading="lazy"
    />
  </picture>
</template>
