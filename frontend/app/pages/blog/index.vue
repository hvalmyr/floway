<script setup lang="ts">
// Дизайна для блога пока нет (нет в макете) — вёрстка минимальная, на
// существующих токенах/компонентах дизайн-системы. Стилизацию поправит
// дизайнер отдельно, когда дойдёт очередь; данные уже реальные из БД.
useSeoMeta({
  title: "Блог — Фловей",
  description:
    "Статьи школы флористики «Фловей»: советы по уходу за цветами, разбор техник, новости школы.",
});

const api = useApi();
const { data: posts, pending, error } = await useAsyncData("blog-posts", () => api.getBlogPosts());

function formatDate(dateString: string | null) {
  if (!dateString) return "";
  return new Date(dateString).toLocaleDateString("ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}
</script>

<template>
  <div class="container flex flex-col gap-48 py-48 sm:py-64 lg:py-80">
    <h1 class="font-display text-h1 text-ink">Блог</h1>

    <p v-if="pending" class="font-body text-body text-ink">Загрузка…</p>
    <p v-else-if="error" class="font-body text-body text-ink">
      Не удалось загрузить статьи. Попробуйте позже.
    </p>
    <p v-else-if="!posts?.length" class="font-body text-body text-ink">
      Пока нет опубликованных статей.
    </p>

    <div v-else class="grid grid-cols-1 gap-24 md:grid-cols-2 lg:grid-cols-3">
      <NuxtLink v-for="post in posts" :key="post.id" :to="`/blog/${post.slug}`" class="block">
        <UiCard>
          <template v-if="post.category" #title>{{ post.category }}</template>
          <p class="mb-8 font-display text-h4 text-ink">{{ post.title }}</p>
          <p v-if="post.publishedAt" class="text-body text-ink">
            {{ formatDate(post.publishedAt) }}
          </p>
        </UiCard>
      </NuxtLink>
    </div>
  </div>
</template>
