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
const { data: masterclass } = await useAsyncData(`masterclass-${slug}`, () => api.getMasterClass(slug));

if (!masterclass.value) {
  throw createError({ statusCode: 404, statusMessage: "Мастер-класс не найден", fatal: true });
}

useSeoMeta({
  title: () => `${masterclass.value?.title} — мастер-класс ФлоВей`,
  description: () => masterclass.value?.shortDescription,
});
</script>

<template>
  <div v-if="masterclass">
    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container">
        <MasterclassCard :masterclass="masterclass" />
      </div>
    </section>

    <section id="apply" class="scroll-mt-64 bg-surface py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm context="masterclass" :related-id="masterclass.id" :title="`Записаться на «${masterclass.title}»`" />
        </div>
      </div>
    </section>
  </div>
</template>
