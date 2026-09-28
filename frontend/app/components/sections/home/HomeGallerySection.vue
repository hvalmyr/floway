<script setup lang="ts">
import type { GalleryPhoto } from "~/types/api";

/** Homepage photo carousel block — see HomeFeaturesSection.vue's doc comment for the pattern. */
defineProps<{
  photos: GalleryPhoto[];
}>();
</script>

<template>
  <section v-if="photos.length" class="py-48 sm:py-64 lg:py-80">
    <div class="container">
      <PhotoCarousel :photos="photos" />
    </div>
  </section>
</template>
