<script setup lang="ts">
// Дизайна для статьи блога пока нет — вёрстка минимальная, на существующих
// токенах дизайн-системы. Контент — HTML, написанный через
// AdminRichTextEditor.vue (см. RichTextContent.vue про доверие к источнику).
const route = useRoute();
const slug = route.params.slug as string;

const api = useApi();
const { data: post } = await useAsyncData(`blog-post-${slug}`, () => api.getBlogPost(slug));

if (!post.value) {
  throw createError({ statusCode: 404, statusMessage: "Статья не найдена", fatal: true });
}

// resolveMediaUrl() reads useRuntimeConfig(), which needs the active Nuxt
// app context — safe here (still inside <script setup>'s synchronous body),
// but NOT safe called lazily from inside a useSeoMeta getter: unhead invokes
// those later, outside that context, and useRuntimeConfig() throws
// (NUXT_E1001) when called from there. Resolve the URL eagerly instead.
const ogImageUrl = post.value?.coverImage ? resolveMediaUrl(post.value.coverImage) : undefined;

useSeoMeta({
  title: () => `${post.value?.title} — блог Фловей`,
  description: () => post.value?.title,
  ogTitle: () => post.value?.title,
  ogDescription: () => post.value?.title,
  ogImage: ogImageUrl,
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
  <div v-if="post" class="container flex flex-col gap-24 py-48 sm:py-64 lg:py-80">
    <NuxtLink to="/blog" class="font-body text-body text-primary hover:underline"
      >← Ко всем статьям</NuxtLink
    >

    <div class="flex flex-col gap-8">
      <p v-if="post.category" class="font-body text-body text-primary">{{ post.category }}</p>
      <h1 class="font-display text-h1 text-ink">{{ post.title }}</h1>
      <p class="font-body text-body text-ink">
        <span v-if="post.author">{{ post.author }}</span>
        <span v-if="post.author && post.publishedAt"> · </span>
        <span v-if="post.publishedAt">{{ formatDate(post.publishedAt) }}</span>
      </p>
    </div>

    <NuxtImg
      v-if="post.coverImage"
      :src="resolveOptimizedMediaUrl(post.coverImage)"
      format="webp"
      :alt="post.title"
      class="aspect-[16/9] w-full max-w-[720px] rounded-md object-cover"
      sizes="400:100vw lg:720px"
      loading="lazy"
    />

    <RichTextContent :source="post.content" class="max-w-[720px]" />

    <div v-if="post.tags.length" class="flex flex-wrap gap-8">
      <UiBadge v-for="tag in post.tags" :key="tag">{{ tag }}</UiBadge>
    </div>
  </div>
</template>
