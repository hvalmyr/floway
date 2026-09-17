<script setup lang="ts">
/**
 * Homepage-only strip promoting gift certificates, placed right after the
 * courses section (see pages/index.vue) — same spot conceptually as a
 * "free shipping" ticker. Content is admin-editable (see
 * /admin/page-content/marquee); rendering itself lives in MarqueeStrip.
 */
const { text } = await usePageContent();
const message = computed(() =>
  text("gift_certificate_marquee_text", "Подарочные сертификаты — порадуйте близких цветами"),
);
const icon = computed(() => text("gift_certificate_marquee_icon", "gift"));
</script>

<template>
  <MarqueeStrip :message="message" :icon="icon" to="/sertifikaty" />
</template>
