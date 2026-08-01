<script setup lang="ts">
// Дизайна для статьи блога пока нет — вёрстка минимальная, на существующих
// токенах дизайн-системы. Контент рендерится как обычный текст с
// сохранением переносов строк (формат контента — plain text, не HTML/markdown).
const route = useRoute();
const slug = route.params.slug as string;

const api = useApi();
const { data: post } = await useAsyncData(`blog-post-${slug}`, () => api.getBlogPost(slug));

if (!post.value) {
  throw createError({ statusCode: 404, statusMessage: "Статья не найдена", fatal: true });
}

useSeoMeta({
  title: () => `${post.value?.title} — блог ФлоВей`,
  description: () => post.value?.title,
});

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
  <div v-if="post" class="container flex flex-col gap-24 py-64 sm:py-96 lg:py-120">
    <NuxtLink to="/blog" class="font-body text-body text-primary hover:underline"
      >← Ко всем статьям</NuxtLink
    >

    <div class="flex flex-col gap-8">
      <p v-if="post.category" class="font-body text-small text-primary">{{ post.category }}</p>
      <h1 class="font-display text-h1 text-ink">{{ post.title }}</h1>
      <p class="font-body text-small text-ink">
        <span v-if="post.author">{{ post.author }}</span>
        <span v-if="post.author && post.publishedAt"> · </span>
        <span v-if="post.publishedAt">{{ formatDate(post.publishedAt) }}</span>
      </p>
    </div>

    <div class="max-w-[720px] whitespace-pre-wrap font-body text-body-l text-ink">
      {{ post.content }}
    </div>

    <div v-if="post.tags.length" class="flex flex-wrap gap-8">
      <UiBadge v-for="tag in post.tags" :key="tag">{{ tag }}</UiBadge>
    </div>
  </div>
</template>
