<script setup lang="ts">
useSeoMeta({
  title: "Подарочный мастер-класс — Фловей",
  description:
    "Подарите мастер-класс по флористике: сертификат школы «Фловей» на любую сумму, курс или мастер-класс.",
});

const { text } = await usePageContent();

const advantages = [
  {
    icon: "gift",
    title: "Любая сумма или мастер-класс",
    description:
      "Оформим сертификат на конкретный мастер-класс или на сумму — получатель выберет сам.",
  },
  {
    icon: "flex-start",
    title: "Свободная дата записи",
    description: "Получатель сам выберет удобный день и время — без привязки к конкретной дате.",
  },
  {
    icon: "checklist",
    title: "Все материалы включены",
    description: "Всё необходимое для мастер-класса уже включено в стоимость сертификата.",
  },
  {
    icon: "calendar-check",
    title: "Долгий срок действия",
    description: "Достаточно времени, чтобы выбрать удобный момент и записаться.",
  },
];
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
            title="Оставить заявку на подарочный мастер-класс"
            lead="Свяжемся с вами, поможем выбрать мастер-класс и оформим сертификат."
          />
        </div>
      </div>
    </section>
  </div>
</template>
