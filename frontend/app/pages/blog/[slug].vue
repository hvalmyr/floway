<script setup lang="ts">
import { ArrowLeft } from "lucide-vue-next";

// Hero block matches the site's shared Hero.vue layout exactly (same
// classes/breakpoints — category/title/author/CTA left, 1:1 cover photo
// right, media-first on mobile) but isn't <Hero> itself, which has no
// category/author slots. It sits directly on the page background as its own
// <section>, NOT inside UiGlassPage — nesting it in that card's own
// container+padding made the cover photo render noticeably smaller than
// every other hero photo on the site (same lg:w-1/2 share of a narrower
// box). Blog is a "closeup" route (see layouts/default.vue), so the ambient
// tree renders bigger/nearer right behind this text — illegible on mobile
// without its own glass backing, hence the glass panel around just the text
// column, matching the photo's own rounded-lg (rounded-pill was tried
// first, matching the button instead — on this column's near-square
// stretched-height box that read as a circle, not a card, so it's back to
// the site's normal card radius). Padding is vertical-only (py-*, no
// horizontal) so the text column's actual content width still matches
// every other hero on the site exactly — horizontal padding here would
// narrow it — py stays comfortably above the rounded-lg corner radius so
// nothing gets visually clipped by the curve. The photo column stays
// plain, at full size, with no glass of its own. The
// running text below is its own UiGlassPage instead, since a reading
// surface benefits article text in a way it doesn't a hero photo. The
// content column there is centered visually (an mx-auto max-w wrapper)
// while its text itself stays left-aligned. Контент — HTML, написанный через
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
  <div v-if="post">
    <section class="py-48 sm:py-64 lg:py-80">
      <div class="container flex flex-col gap-32 lg:flex-row lg:items-stretch lg:gap-64">
        <div
          class="order-2 flex flex-col items-start gap-24 rounded-lg bg-white/55 py-32 backdrop-blur backdrop-saturate-150 sm:py-40 lg:order-1 lg:w-1/2 lg:py-48"
        >
          <div class="flex flex-col gap-8">
            <p v-if="post.category" class="font-body text-h4 font-medium text-primary">
              {{ post.category }}
            </p>
            <h1 class="font-display text-h1 text-ink">{{ post.title }}</h1>
            <p class="font-body text-body text-ink">
              <span v-if="post.author">{{ post.author }}</span>
              <span v-if="post.author && post.publishedAt"> · </span>
              <span v-if="post.publishedAt">{{ formatDate(post.publishedAt) }}</span>
            </p>
          </div>
          <UiButton variant="outline" to="/blog" block class="mt-auto">
            <ArrowLeft class="size-24" aria-hidden="true" />
            Ко всем статьям
          </UiButton>
        </div>
        <div class="order-1 lg:order-2 lg:w-1/2">
          <UiContentImage
            v-if="post.coverImage"
            :src="resolveOptimizedMediaUrl(post.coverImage)"
            :alt="post.title"
            class="aspect-square w-full rounded-lg object-cover"
            sizes="400:100vw lg:576px"
          />
          <UiMediaPlaceholder v-else aspect="1/1" />
        </div>
      </div>
    </section>

    <UiGlassPage>
      <div class="mx-auto flex w-full max-w-[720px] flex-col gap-24">
        <RichTextContent :source="post.content" />

        <div v-if="post.tags.length" class="flex flex-wrap gap-8">
          <UiBadge v-for="tag in post.tags" :key="tag">{{ tag }}</UiBadge>
        </div>
      </div>
    </UiGlassPage>
  </div>
</template>
