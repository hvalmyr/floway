<script setup lang="ts">
const props = defineProps<{ error: { statusCode: number; statusMessage?: string } }>();

useSeoMeta({ title: () => `${props.error.statusCode} — flo-way` });
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppHeader />
    <main class="flex flex-1 items-center justify-center px-16 py-64">
      <div class="flex flex-col items-center gap-24 text-center">
        <IconSnowflake class="h-64 w-auto text-primary" />
        <h1 class="font-display text-h1 text-ink">{{ error.statusCode }}</h1>
        <p class="font-body text-body text-ink">
          {{ error.statusCode === 404 ? "Страница не найдена" : "Что-то пошло не так" }}
        </p>
        <UiButton variant="primary" to="/" @click="clearError({ redirect: '/' })">
          На главную
        </UiButton>
      </div>
    </main>
    <AppFooter />
  </div>
</template>
