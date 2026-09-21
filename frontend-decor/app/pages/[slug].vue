<script setup lang="ts">
const route = useRoute();
const slug = route.params.slug as string;

const api = useApi();
const { data: page } = await useAsyncData(`landing-page-${slug}`, () => api.getLandingPage(slug));

if (!page.value) {
  throw createError({ statusCode: 404, statusMessage: "Страница не найдена", fatal: true });
}

useSeoMeta({
  title: () => page.value?.metaTitle || `${page.value?.h1} — flo-way`,
  description: () => page.value?.metaDescription || undefined,
});

const blocks = computed(
  () => page.value?.blocks.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);
const faqItems = computed(
  () => page.value?.faqItems.slice().sort((a, b) => a.sortOrder - b.sortOrder) ?? [],
);
</script>

<template>
  <div v-if="page">
    <Hero>
      <template #title>{{ page.h1 }}</template>
      <template #lead
        >Оформляем на сезон или на один праздник — рассчитываем стоимость индивидуально после
        бесплатного выезда и замеров.</template
      >
      <template #actions>
        <UiButton variant="primary" to="#apply">Оставить заявку</UiButton>
      </template>
      <template v-if="blocks[0]?.image" #media>
        <UiHeroPicture :src="resolveOptimizedMediaUrl(blocks[0].image)" :alt="page.h1" />
      </template>
    </Hero>

    <section v-if="blocks.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-64">
        <div
          v-for="(block, i) in blocks"
          :key="block.id"
          class="flex flex-col gap-24 lg:flex-row lg:items-center lg:gap-64"
        >
          <div class="lg:w-1/2" :class="i % 2 === 1 ? 'lg:order-2' : ''">
            <UiContentImage
              v-if="block.image"
              :src="resolveOptimizedMediaUrl(block.image)"
              :alt="page.h1"
              class="aspect-[4/3] w-full rounded-lg object-cover"
              sizes="400:100vw lg:50vw"
            />
            <UiMediaPlaceholder v-else aspect="4/3" />
          </div>
          <div class="lg:w-1/2" :class="i % 2 === 1 ? 'lg:order-1' : ''">
            <MarkdownContent :source="block.text" class="font-body text-body text-ink" />
          </div>
        </div>
      </div>
    </section>

    <section
      id="apply"
      class="scroll-mt-64 bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:scroll-mt-96 lg:py-80"
    >
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm
            :object-type-id="page.objectTypeId"
            :title="`Оставить заявку на оформление: ${page.objectType.name.toLowerCase()}`"
          />
        </div>
      </div>
    </section>

    <section v-if="page.faqVisible && faqItems.length" class="py-48 sm:py-64 lg:py-80">
      <div class="container">
        <FaqSection
          :title="page.faqTitle || 'Вопросы и ответы'"
          :description="page.faqDescription"
          :items="faqItems"
        />
      </div>
    </section>
  </div>
</template>
