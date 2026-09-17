<script setup lang="ts">
/**
 * Masterclasses-page strip promoting the catalog below it, linking down to
 * the apply form. Content is admin-editable (see
 * /admin/page-content/masterclasses-marquee); rendering itself lives in
 * MarqueeStrip.
 */
const { text } = await usePageContent();
const message = computed(() =>
  text("masterclasses_marquee_text", "Мастер-классы по флористике — запишитесь прямо сейчас"),
);
const icon = computed(() => text("masterclasses_marquee_icon", "tulips"));
</script>

<template>
  <MarqueeStrip :message="message" :icon="icon" to="#apply" />
</template>
