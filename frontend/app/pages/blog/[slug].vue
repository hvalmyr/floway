<script setup lang="ts">
// Hero block matches the site's shared Hero.vue layout (category/title/
// author/CTA left, 1:1 cover photo right, media-first on mobile) — reimplemented
// here rather than reusing that component since it's a standalone <section>
// with its own container/py padding meant to sit directly on the page
// background, while this page's content lives inside UiGlassPage's own
// container+padding. The content column below is centered visually (an
// mx-auto max-w wrapper) while its text itself stays left-aligned. Контент —
// HTML, написанный через AdminRichTextEditor.vue (см. RichTextContent.vue про
// доверие к источнику).
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
  title: () => `${post.value?.metaTitle || post.value?.title} — блог Фловей`,
  description: () => post.value?.metaDescription || post.value?.title,
  ogTitle: () => post.value?.metaTitle || post.value?.title,
  ogDescription: () => post.value?.metaDescription || post.value?.title,
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
  <UiGlassPage v-if="post">
    <div class="flex flex-col gap-24 lg:gap-64">
      <div class="flex flex-col gap-24 lg:flex-row lg:items-stretch lg:gap-64">
        <div class="order-2 flex flex-col items-start gap-24 lg:order-1 lg:w-1/2">
          <div class="flex flex-col gap-8">
            <p v-if="post.category" class="font-body text-body text-primary">
              {{ post.category }}
            </p>
            <h1 class="font-display text-h1 text-ink">{{ post.title }}</h1>
            <p class="font-body text-body text-ink">
              <span v-if="post.author">{{ post.author }}</span>
              <span v-if="post.author && post.publishedAt"> · </span>
              <span v-if="post.publishedAt">{{ formatDate(post.publishedAt) }}</span>
            </p>
          </div>
          <UiButton variant="outline" to="/blog" class="mt-auto">← Ко всем статьям</UiButton>
        </div>
        <div class="order-1 lg:order-2 lg:w-1/2">
          <UiContentImage
            v-if="post.coverImage"
            :src="resolveOptimizedMediaUrl(post.coverImage)"
            :alt="post.title"
            class="aspect-square w-full rounded-lg object-cover"
            sizes="400:100vw lg:480px"
          />
          <UiMediaPlaceholder v-else aspect="1/1" />
        </div>
      </div>

      <div class="mx-auto flex w-full max-w-[720px] flex-col gap-24">
        <RichTextContent :source="post.content" />

        <div v-if="post.tags.length" class="flex flex-wrap gap-8">
          <UiBadge v-for="tag in post.tags" :key="tag">{{ tag }}</UiBadge>
        </div>
      </div>
    </div>
  </UiGlassPage>
</template>
