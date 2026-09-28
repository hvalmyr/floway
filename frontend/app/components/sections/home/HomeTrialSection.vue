<script setup lang="ts">
/** Homepage "Пробное занятие" block — see HomeFeaturesSection.vue's doc comment for the pattern. */
defineProps<{
  text: (key: string, fallback?: string) => string;
  glassClass: string;
}>();
</script>

<template>
  <section id="trial" class="scroll-mt-64 py-48 sm:py-64 lg:scroll-mt-96 lg:py-80">
    <div class="container flex flex-col gap-48">
      <SectionHeading color="primary">
        {{ text("trial_section_heading", "Попробуйте флористику на практике") }}
        <template #lead>
          {{
            text(
              "trial_section_description",
              "Вы научитесь собирать круглый букет в спиральной технике: поймёте принцип работы со спиралью, научитесь уверенно удерживать букет в руках во время сборки и правильно подвязывать букет. Кроме практики, вы сможете познакомиться с нашим педагогом, узнать, как проходят занятия в школе, задать все интересующие вопросы и понять, подходит ли вам обучение.",
            )
          }}
        </template>
      </SectionHeading>

      <div class="grid grid-cols-1 gap-32 md:grid-cols-2 md:gap-64">
        <div class="order-2 flex w-full flex-col gap-24 md:order-1">
          <div class="flex flex-col gap-16 rounded-md px-16 py-24" :class="glassClass">
            <h3 class="font-body text-h4 text-ink">
              {{ text("trial_heading", "Пробное занятие") }}
            </h3>
            <!-- Single \n (not \n\n) so duration+price render as one tight
            paragraph with a <br> between them, not two separately-spaced
            ones — otherwise the gap between them (markdown paragraph
            margin) was bigger than the gap to the heading above, so
            duration read as grouped with "Пробное занятие" instead of
            with the price right under it. -->
            <MarkdownContent
              :source="
                text('trial_description', 'Продолжительность: 2,5 часа.\nСтоимость: 3 000 ₽.')
              "
            />
          </div>
          <LazyApplyForm context="trial_lesson" title="" bare class="w-full" hydrate-on-visible />
        </div>
        <!-- TODO: заменить на видео с пробным уроком, когда оно будет готово (пока фото). -->
        <!-- Портретное 9:16, растянуто до ширины колонки — но не выше 80%
        экрана: max-h ограничивает высоту, а aspect-ratio при этом сжимает
        и ширину пропорционально, так что соотношение сторон не ломается
        (стандартное поведение aspect-ratio + max-height у браузера, без
        JS). mx-auto центрирует, когда из-за max-h ширина не дотягивает до
        полной колонки. min-h-0 — эта колонка ещё и grid-item, а
        min-height:auto по умолчанию у grid/flex-item позволяет
        собственному соотношению сторон загруженного фото (если оно не
        9:16) перебить aspect-ratio и растянуть блок; см. коммит с
        разбором в CourseCard.vue. -->
        <UiContentImage
          v-if="text('home_trial_image')"
          :src="resolveOptimizedMediaUrl(text('home_trial_image'))"
          alt=""
          class="order-1 mx-auto aspect-[9/16] max-h-[80vh] max-w-full min-h-0 rounded-lg object-cover md:sticky md:top-96 md:order-2"
          sizes="400:100vw md:50vw"
        />
        <div
          v-else
          class="order-1 mx-auto aspect-[9/16] max-h-[80vh] max-w-full rounded-lg bg-primary md:sticky md:top-96 md:order-2"
        />
      </div>
    </div>
  </section>
</template>
