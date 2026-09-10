<script setup lang="ts">
import { displayStyleColorClasses } from "~/constants/display-style-colors";

// Card grid mirrors the homepage's course section exactly (same
// flex-wrap + w-[calc(...)] widths, same container/py rhythm, no glass-card
// wrapper around the whole page) rather than UiGlassPage's own container —
// nesting the grid inside that card's extra p-24/48/64 padding shrank every
// card noticeably narrower than a course card at the same breakpoint. The
// article page (blog/[slug].vue) is still a UiGlassPage — running text
// benefits from that reading surface in a way a card grid doesn't.
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
  <section class="py-48 sm:py-64 lg:py-80">
    <div class="container flex flex-col gap-48">
      <h1 class="font-display text-h1 text-ink">Блог</h1>

      <p
        v-if="pending"
        class="rounded-md bg-white/55 px-24 py-16 font-body text-body text-ink backdrop-blur backdrop-saturate-150"
      >
        Загрузка…
      </p>
      <p
        v-else-if="error"
        class="rounded-md bg-white/55 px-24 py-16 font-body text-body text-ink backdrop-blur backdrop-saturate-150"
      >
        Не удалось загрузить статьи. Попробуйте позже.
      </p>
      <p
        v-else-if="!posts?.length"
        class="rounded-md bg-white/55 px-24 py-16 font-body text-body text-ink backdrop-blur backdrop-saturate-150"
      >
        Пока нет опубликованных статей.
      </p>

      <div v-else class="flex flex-wrap justify-center gap-24 lg:gap-32">
        <NuxtLink
          v-for="post in posts"
          :key="post.id"
          :to="`/blog/${post.slug}`"
          class="block w-full sm:w-[calc(50%-12px)] lg:w-[calc(33.333%-22px)]"
        >
          <UiCard variant="custom" :class="displayStyleColorClasses[post.displayStyle]">
            <template v-if="post.coverImage" #media>
              <UiContentImage
                :src="resolveOptimizedMediaUrl(post.coverImage)"
                :alt="post.title"
                class="aspect-square w-full rounded-sm object-cover"
                sizes="400:100vw sm:50vw lg:400px"
              />
            </template>
            <template v-if="post.category" #title
              ><span class="font-body text-body">{{ post.category }}</span></template
            >
            <p class="mb-8 font-display text-h4">{{ post.title }}</p>
            <p v-if="post.publishedAt" class="text-body">
              {{ formatDate(post.publishedAt) }}
            </p>
          </UiCard>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>
