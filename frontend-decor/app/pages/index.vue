<script setup lang="ts">
// Placeholder homepage — proves the design system/layout pipeline builds
// and renders end to end. Real content (offer, formats, portfolio,
// advantages, FAQ, apply form — п. 6 ТЗ) lands in Phase 4 once the actual
// object-type/landing-page content exists in the admin.
</script>

<template>
  <div>
    <Hero>
      <template #title>Новогоднее оформление домов, офисов и мероприятий</template>
      <template #lead>
        Оформляем на сезон или на один праздник — рассчитываем стоимость после бесплатного выезда и
        замеров.
      </template>
      <template #actions>
        <UiButton variant="primary" to="/contacts">Оставить заявку</UiButton>
      </template>
    </Hero>
  </div>
</template>
