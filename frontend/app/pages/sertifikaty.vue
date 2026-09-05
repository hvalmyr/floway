<script setup lang="ts">
useSeoMeta({
  title: "Подарочный мастер-класс — Фловей",
  description:
    "Подарите мастер-класс по флористике: сертификат школы «Фловей» на любую сумму, курс или мастер-класс.",
});

const api = useApi();
const { text } = await usePageContent();

const { data: featuresData } = await useAsyncData("gift-certificate-features", () =>
  api.getFeatures("gift_certificate"),
);

const advantages = computed(
  () =>
    featuresData.value
      ?.slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((f) => ({ icon: f.icon, title: f.title, description: f.description })) ?? [],
);
</script>

<template>
  <div>
    <Hero>
      <template #title>{{
        text("gift_certificate_hero_title", "Подарите мастер-класс по флористике")
      }}</template>
      <template #lead>
        {{
          text(
            "gift_certificate_hero_lead",
            "Сертификат «Фловей» — оригинальный подарок для тех, кто любит цветы и творчество. Получатель сам выберет мастер-класс и удобную дату визита в школу.",
          )
        }}
      </template>
      <template #actions>
        <UiButton variant="primary" to="#apply">Оставить заявку</UiButton>
        <UiButton variant="outline" to="/masterclasses">Мастер-классы</UiButton>
      </template>
      <template v-if="text('gift_certificate_hero_image')" #media>
        <UiHeroPicture
          :src="resolveOptimizedMediaUrl(text('gift_certificate_hero_image'))"
          alt=""
        />
      </template>
    </Hero>

    <section class="bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary" on-glass>
          {{ text("gift_certificate_features_heading", "Почему это отличный подарок") }}
          <template #lead>{{
            text(
              "gift_certificate_features_lead",
              "Дарите не вещь, а впечатление — тёплый опыт создания своими руками.",
            )
          }}</template>
        </SectionHeading>
        <FeatureGrid :items="advantages" />
      </div>
    </section>

    <section
      id="apply"
      class="scroll-mt-64 bg-surface/55 py-48 backdrop-blur backdrop-saturate-150 sm:py-64 lg:scroll-mt-96 lg:py-80"
    >
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm
            context="masterclass"
            title="Оставить заявку на подарочный сертификат"
            lead="Свяжемся с вами, поможем выбрать мастер-класс и оформим сертификат."
          />
        </div>
      </div>
    </section>
  </div>
</template>
