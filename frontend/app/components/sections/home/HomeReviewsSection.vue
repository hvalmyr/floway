<script setup lang="ts">
 /** Homepage "Отзывы" (Yandex Maps widget) block — fully static, no props. */
</script>

<template>
  <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
    <div class="container flex flex-col gap-48">
      <SectionHeading color="primary" on-glass>
        Отзывы
        <template #lead
          >Смотрите отзывы о школе на Яндекс Картах — заходите оставить свой.</template
        >
      </SectionHeading>
      <!-- Официальный виджет отзывов Яндекс Карт — сам виджет задаёт
      style width:100%/height:100% (заказчик прислал), поэтому обёртка
      ниже даёт ему конкретную высоту, внутри которой он растягивается. -->
      <div class="mx-auto h-[950px] w-full max-w-[400px]">
        <iframe
          src="https://yandex.ru/maps-reviews-widget/83657275642?comments"
          title="Отзывы о школе «Фловей» на Яндекс Картах"
          loading="lazy"
          style="
            width: 100%;
            height: 100%;
            border: 1px solid #e6e6e6;
            border-radius: 8px;
            box-sizing: border-box;
          "
        />
      </div>
    </div>
  </section>
</template>
