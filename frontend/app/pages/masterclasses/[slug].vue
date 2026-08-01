<script setup lang="ts">
definePageMeta({
  validate: async (route) => {
    const api = useApi();
    const mc = await api.getMasterClass(route.params.slug as string);
    return !!mc;
  },
});

const route = useRoute();
const slug = route.params.slug as string;

const api = useApi();
const { data: masterclass } = await useAsyncData(`masterclass-${slug}`, () =>
  api.getMasterClass(slug),
);

if (!masterclass.value) {
  throw createError({ statusCode: 404, statusMessage: "Мастер-класс не найден", fatal: true });
}

useSeoMeta({
  title: () => `${masterclass.value?.title} — мастер-класс ФлоВей`,
  description: () => masterclass.value?.shortDescription,
});

const infoItems = computed(() => {
  const mc = masterclass.value!;
  const duration = mc.duration ?? "уточняется";
  if (mc.priceGroup == null) return [duration, "цена уточняется"];
  if (mc.priceIndividual != null) {
    return [
      duration,
      `${mc.priceGroup.toLocaleString("ru-RU")}₽ или ${mc.priceIndividual.toLocaleString("ru-RU")}₽ [*]`,
    ];
  }
  return [duration, `${mc.priceGroup.toLocaleString("ru-RU")}₽`];
});
</script>

<template>
  <div v-if="masterclass">
    <Hero>
      <template #title>Мастер-класс “{{ masterclass.title }}”</template>
      <template #lead>{{ masterclass.shortDescription }}</template>
      <template #actions>
        <UiButton variant="primary" size="lg" to="#apply"
          >Оставить заявку</UiButton
        >
      </template>
      <template v-if="masterclass.coverImage" #media>
        <img
          :src="resolveMediaUrl(masterclass.coverImage)"
          :alt="masterclass.title"
          class="aspect-[4/3] w-full rounded-lg object-cover"
        />
      </template>
    </Hero>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-24">
        <p v-if="masterclass.fullDescription" class="font-body text-body-l text-ink">
          {{ masterclass.fullDescription }}
        </p>
        <p v-if="masterclass.endingText" class="font-body text-body-l text-ink">
          {{ masterclass.endingText }}
        </p>
        <UiInfoRow :items="infoItems" />
        <p v-if="masterclass.priceDescription" class="font-body text-small text-ink/70">
          {{ masterclass.priceDescription }}
        </p>
      </div>
    </section>

    <section id="apply" class="scroll-mt-64 bg-surface py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm
            context="masterclass"
            :related-id="masterclass.id"
            :title="`Записаться на “${masterclass.title}”`"
          />
        </div>
      </div>
    </section>
  </div>
</template>
