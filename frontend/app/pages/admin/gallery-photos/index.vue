<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Фотокарусели</h1>

    <div class="mt-8 rounded border border-gray-200 bg-white p-4 sm:p-6">
      <AdminPhotoCarouselManager
        api-path="/api/v1/gallery-photos"
        title="Фотогалерея (карусель на главной)"
        description="Вертикальные фото, показанные в слайд-шоу между блоками «Преимущества» и «О школе». Порядок сортировки определяет порядок слайдов."
      />
    </div>

    <div class="mt-8 rounded border border-gray-200 bg-white p-4 sm:p-6">
      <AdminPhotoCarouselManager
        api-path="/api/v1/gift-certificate-carousel-photos"
        title="Карусель (страница сертификатов)"
        description="Вертикальные фото, показанные в слайд-шоу сразу после Hero-блока на странице подарочных сертификатов. Порядок сортировки определяет порядок слайдов."
      />
    </div>
  </div>
</template>
