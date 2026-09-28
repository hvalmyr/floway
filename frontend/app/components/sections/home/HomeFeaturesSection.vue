<script setup lang="ts">
/**
 * Homepage "Почему стоит учиться в школе «Фловей»?" block — one of the
 * order/visibility-managed sections in /admin/page-content/home-sections
 * (see HOME_SECTION_COMPONENTS in pages/index.vue). `text` is the resolved
 * usePageContent() lookup, passed down rather than re-fetched so every
 * section shares the one page-content request the parent already awaited.
 */
defineProps<{
  text: (key: string, fallback?: string) => string;
  features: { icon: string; title: string; description: string }[];
}>();
</script>

<template>
  <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
    <div class="container flex flex-col gap-48">
      <SectionHeading color="primary" on-glass>
        {{ text("home_features_heading", "Почему стоит учиться в школе «Фловей»?") }}
        <template #lead>
          {{
            text(
              "home_features_lead",
              "Мы создали школу, в которой удобно учиться, легко развиваться и получать реальные практические навыки флористики.",
            )
          }}
        </template>
      </SectionHeading>
      <FeatureGrid :items="features" />
    </div>
  </section>
</template>
