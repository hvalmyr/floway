<script setup lang="ts">
useSeoMeta({
  title: "Портфолио — flo-way",
  description: "Примеры новогоднего оформления домов, офисов, магазинов, ресторанов и отелей.",
});

const api = useApi();

const { data: photosData } = await useAsyncData("portfolio-photos", () => api.getGalleryPhotos());
const photos = computed(
  () => photosData.value?.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);

const { data: objectTypesData } = await useAsyncData("portfolio-object-types", () =>
  api.getObjectTypes(),
);
const objectTypes = computed(() => objectTypesData.value?.filter((t) => t.visible) ?? []);

const objectTypeFilter = ref<number | null>(null);
const formatFilter = ref<"" | "season" | "event">("");

const filteredPhotos = computed(() =>
  photos.value.filter((photo) => {
    if (objectTypeFilter.value !== null && photo.objectTypeId !== objectTypeFilter.value) {
      return false;
    }
    if (formatFilter.value && photo.format !== formatFilter.value) return false;
    return true;
  }),
);
</script>

<template>
  <div class="container flex flex-col gap-32 py-48 sm:py-64 lg:py-80">
    <SectionHeading color="ink">
      Портфолио
      <template #lead>Примеры нашего новогоднего оформления.</template>
    </SectionHeading>

    <div class="flex flex-wrap items-center gap-16">
      <div class="flex flex-wrap gap-8">
        <button
          type="button"
          class="rounded-pill border px-16 py-8 font-body text-body"
          :class="
            objectTypeFilter === null
              ? 'border-primary bg-primary text-white'
              : 'border-primary/40 text-ink'
          "
          @click="objectTypeFilter = null"
        >
          Все объекты
        </button>
        <button
          v-for="type in objectTypes"
          :key="type.id"
          type="button"
          class="rounded-pill border px-16 py-8 font-body text-body"
          :class="
            objectTypeFilter === type.id
              ? 'border-primary bg-primary text-white'
              : 'border-primary/40 text-ink'
          "
          @click="objectTypeFilter = type.id"
        >
          {{ type.name }}
        </button>
      </div>
      <div class="flex flex-wrap gap-8">
        <button
          type="button"
          class="rounded-pill border px-16 py-8 font-body text-body"
          :class="formatFilter === '' ? 'border-ink bg-ink text-white' : 'border-ink/40 text-ink'"
          @click="formatFilter = ''"
        >
          Любой формат
        </button>
        <button
          type="button"
          class="rounded-pill border px-16 py-8 font-body text-body"
          :class="
            formatFilter === 'season' ? 'border-ink bg-ink text-white' : 'border-ink/40 text-ink'
          "
          @click="formatFilter = 'season'"
        >
          На сезон
        </button>
        <button
          type="button"
          class="rounded-pill border px-16 py-8 font-body text-body"
          :class="
            formatFilter === 'event' ? 'border-ink bg-ink text-white' : 'border-ink/40 text-ink'
          "
          @click="formatFilter = 'event'"
        >
          На праздник
        </button>
      </div>
    </div>

    <PhotoCarousel v-if="filteredPhotos.length" :photos="filteredPhotos" />
    <p v-else class="rounded-md bg-white/55 p-24 text-center font-body text-body text-ink">
      Пока нет фото для этого фильтра.
    </p>
  </div>
</template>
