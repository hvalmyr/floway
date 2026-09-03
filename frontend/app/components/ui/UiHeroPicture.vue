<script setup lang="ts">
/**
 * The 3 page heroes (home, course, masterclass) are each that page's LCP
 * element — eager-loaded, `fetchpriority="high"`, real avif→webp→jpeg
 * `<picture>` fallback chain. Built by hand with `useImage()` instead of
 * `<NuxtPicture>` because avif needs a much lower `quality` than webp/jpeg
 * for an equivalent-looking result (~50-60 vs ~80) — <NuxtPicture> applies
 * one shared quality to every format in its `format` list, which measured
 * out to avif files *larger* than the jpeg fallback on real photos here.
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

const img = useImage();
const avif = computed(() =>
  img.getSizes(props.src, { sizes: HERO_SIZES, modifiers: { format: "avif", quality: 55 } }),
);
const webp = computed(() =>
  img.getSizes(props.src, { sizes: HERO_SIZES, modifiers: { format: "webp", quality: 80 } }),
);
const fallback = computed(() =>
  img.getSizes(props.src, { sizes: HERO_SIZES, modifiers: { format: "jpeg", quality: 80 } }),
);
</script>

<template>
  <picture>
    <source type="image/avif" :srcset="avif.srcset" :sizes="avif.sizes" />
    <source type="image/webp" :srcset="webp.srcset" :sizes="webp.sizes" />
    <img
      :src="fallback.src"
      :srcset="fallback.srcset"
      :sizes="fallback.sizes"
      :alt="alt"
      class="aspect-square w-full rounded-lg object-cover"
      loading="eager"
      fetchpriority="high"
    />
  </picture>
</template>
