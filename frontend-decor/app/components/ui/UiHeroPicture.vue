<script setup lang="ts">
/**
 * The 3 page heroes (home, course, masterclass) are each that page's LCP
 * element — eager-loaded, `fetchpriority="high"`, real avif→webp→jpeg
 * `<picture>` fallback chain. Built by hand with `useImage()` instead of
 * `<NuxtPicture>` because avif needs a much lower `quality` than webp/jpeg
 * for an equivalent-looking result — <NuxtPicture> applies one shared
 * quality to every format in its `format` list, which measured out to
 * avif files *larger* than the jpeg fallback on real photos here.
 *
 * webp/jpeg also split mobile vs desktop quality (avif doesn't need to —
 * see UiContentImage's doc comment for why). Quality numbers come from
 * useImageQuality() (admin-editable, see /admin/page-content) rather than
 * being hardcoded.
 *
 * `sizes` is fixed (matches Hero.vue's `lg:w-1/2` of a max-1280px
 * container) rather than a prop — all 3 call sites share the exact same
 * layout, so there's nothing to parameterize yet.
 *
 * @example
 * <UiHeroPicture :src="resolveOptimizedMediaUrl(course.blocks[0].blockCover)" :alt="course.name" />
 */
const props = defineProps<{
  src: string;
  alt: string;
}>();

const HERO_SIZES = "400:100vw lg:576px";
const MOBILE_MEDIA = "(max-width: 767px)";
const DESKTOP_MEDIA = "(min-width: 768px)";

const img = useImage();
const quality = useImageQuality();

const avif = computed(() =>
  img.getSizes(props.src, {
    sizes: HERO_SIZES,
    modifiers: { format: "avif", quality: quality.value.avif },
  }),
);
const webpMobile = computed(() =>
  img.getSizes(props.src, {
    sizes: HERO_SIZES,
    modifiers: { format: "webp", quality: quality.value.mobile },
  }),
);
const webpDesktop = computed(() =>
  img.getSizes(props.src, {
    sizes: HERO_SIZES,
    modifiers: { format: "webp", quality: quality.value.desktop },
  }),
);
const jpegMobile = computed(() =>
  img.getSizes(props.src, {
    sizes: HERO_SIZES,
    modifiers: { format: "jpeg", quality: quality.value.mobile },
  }),
);
const jpegDesktop = computed(() =>
  img.getSizes(props.src, {
    sizes: HERO_SIZES,
    modifiers: { format: "jpeg", quality: quality.value.desktop },
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
    <source
      :media="MOBILE_MEDIA"
      type="image/jpeg"
      :srcset="jpegMobile.srcset"
      :sizes="jpegMobile.sizes"
    />
    <source
      :media="DESKTOP_MEDIA"
      type="image/jpeg"
      :srcset="jpegDesktop.srcset"
      :sizes="jpegDesktop.sizes"
    />
    <img
      :src="jpegDesktop.src"
      :srcset="jpegDesktop.srcset"
      :sizes="jpegDesktop.sizes"
      :alt="alt"
      class="aspect-square w-full rounded-lg object-cover"
      loading="eager"
      fetchpriority="high"
    />
  </picture>
</template>
