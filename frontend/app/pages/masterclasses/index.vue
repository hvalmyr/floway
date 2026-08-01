<script setup lang="ts">
import IconCalendarCheck from "~/components/ui/IconCalendarCheck.vue";
import IconChecklist from "~/components/ui/IconChecklist.vue";
import IconGift from "~/components/ui/IconGift.vue";
import IconMcLocation from "~/components/ui/IconMcLocation.vue";
import IconPeopleTrio from "~/components/ui/IconPeopleTrio.vue";
import IconTulips from "~/components/ui/IconTulips.vue";

useSeoMeta({
  title: "Мастер-классы по флористике в Москве — ФлоВей",
  description:
    "Мастер-классы по флористике в свободном графике: букеты и композиции, все материалы включены.",
});

const api = useApi();
const { data: masterclasses } = await useAsyncData("masterclasses-list", () =>
  api.getMasterClasses(),
);

const features = [
  {
    icon: IconCalendarCheck,
    title: "свободный график",
    description:
      "Все мастер-классы проходят по системе свободного посещения. После покупки вы самостоятельно выбираете удобные дату и время. Записаться можно минимум за один день через сайт, в мессенджерах или по телефону.",
  },
  {
    icon: IconPeopleTrio,
    title: "для одного, двоих или компании",
    description:
      "Мы проводим индивидуальные мастер-классы, занятия для двоих и групповые мастер-классы. Вы можете прийти самостоятельно, устроить творческое свидание, отметить день рождения, провести корпоратив или любое другое мероприятие в атмосфере живых цветов.",
  },
  {
    icon: IconMcLocation,
    title: "в школе или с выездом",
    description:
      "Мастер-классы проходят в школе флористики ФлоВей, а также в выездном формате. При необходимости преподаватель приедет на вашу площадку со всеми необходимыми материалами и проведёт мастер-класс для любого количества участников.",
  },
  {
    icon: IconChecklist,
    title: "индивидуальная программа",
    description:
      "На сайте представлены самые популярные мастер-классы, но ими выбор не ограничивается. Если вы хотите освоить определённую технику, создать конкретный букет или композицию либо провести занятие по собственной программе, мы разработаем мастер-класс специально под ваш запрос.",
  },
  {
    icon: IconGift,
    title: "подарочный сертификат",
    description:
      "Любой мастер-класс можно оформить в виде подарочного сертификата. Вы можете выбрать сертификат на конкретное занятие или универсальный сертификат, чтобы получатель самостоятельно выбрал мастер-класс и удобную дату посещения.",
  },
  {
    icon: IconTulips,
    title: "живые цветы и максимум практики",
    description:
      "Все мастер-классы проходят исключительно на живых цветах. Большую часть занятия занимает практика, а преподаватель сопровождает вас на каждом этапе, помогая освоить технику и получить удовольствие от работы с материалом.",
  },
];
</script>

<template>
  <div>
    <Hero>
      <template #title>Мастер-классы по флористике в свободном графике</template>
      <template #actions>
        <UiButton variant="primary" size="lg" to="#apply" class="w-full sm:w-auto"
          >Оставить заявку</UiButton
        >
      </template>
    </Hero>

    <section class="py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-48">
        <SectionHeading color="primary">Почему стоит выбрать мастер-класс «ФлоВей»</SectionHeading>
        <FeatureGrid :items="features" />
      </div>
    </section>

    <section class="bg-surface py-64 sm:py-96 lg:py-120">
      <div class="container flex flex-col gap-64">
        <MasterclassCard
          v-for="(mc, i) in masterclasses"
          :key="mc.id"
          :masterclass="mc"
          link-to-detail
          :accent="i % 2 === 0"
        />
      </div>
    </section>

    <section id="apply" class="scroll-mt-64 py-64 sm:py-96 lg:scroll-mt-96 lg:py-120">
      <div class="container">
        <div class="mx-auto max-w-[720px]">
          <ApplyForm context="masterclass" title="Оставить заявку на мастер-класс" />
        </div>
      </div>
    </section>
  </div>
</template>
