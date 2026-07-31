<script setup lang="ts">
// useFocusTrap actually lives in @vueuse/integrations (wraps the `focus-trap`
// package), not @vueuse/core — the brief named the wrong package; installed
// the correct one instead of hand-rolling a trap.
import { useFocusTrap } from "@vueuse/integrations/useFocusTrap";
import { Menu, X } from "lucide-vue-next";
import { nextTick, onUnmounted, ref, watch } from "vue";

const route = useRoute();
const isOpen = ref(false);
const panelRef = ref<HTMLElement | null>(null);

const { activate, deactivate } = useFocusTrap(panelRef, {
  immediate: false,
  escapeDeactivates: true,
  onDeactivate: () => {
    isOpen.value = false;
  },
});

watch(isOpen, async (open) => {
  if (open) {
    document.body.style.overflow = "hidden";
    await nextTick();
    activate();
  } else {
    deactivate();
    document.body.style.overflow = "";
  }
});

onUnmounted(() => {
  document.body.style.overflow = "";
});

const leftLinks = [
  { to: "/#courses", label: "Курсы" },
  { to: "/masterclasses", label: "Мастер-классы" },
];
const rightLinks = [
  { to: "/blog", label: "Блог" },
  { to: "/contacts", label: "Контакты" },
];

function isActive(to: string) {
  return !to.includes("#") && to !== "/" && route.path.startsWith(to);
}
</script>

<template>
  <header class="sticky top-0 z-40 h-[64px] border-b border-line bg-white lg:h-[88px]">
    <div class="container flex h-full items-center justify-between lg:grid lg:grid-cols-3">
      <NuxtLink to="/" class="font-display text-h4 text-ink-900 lg:hidden">ФлоВей</NuxtLink>

      <nav class="hidden items-center gap-32 lg:flex lg:justify-self-start" aria-label="Основная навигация">
        <NuxtLink
          v-for="link in leftLinks"
          :key="link.to"
          :to="link.to"
          class="border-b-2 border-transparent text-body text-ink-700 hover:text-primary-600"
          :class="isActive(link.to) ? 'border-primary-600 text-primary-600' : ''"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <NuxtLink to="/" class="hidden font-display text-h4 text-ink-900 lg:block lg:justify-self-center">
        ФлоВей
      </NuxtLink>

      <nav class="hidden items-center gap-32 lg:flex lg:justify-self-end" aria-label="Дополнительная навигация">
        <NuxtLink
          v-for="link in rightLinks"
          :key="link.to"
          :to="link.to"
          class="border-b-2 border-transparent text-body text-ink-700 hover:text-primary-600"
          :class="isActive(link.to) ? 'border-primary-600 text-primary-600' : ''"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <button
        type="button"
        class="grid size-[44px] place-items-center rounded-sm text-ink-900 lg:hidden"
        aria-haspopup="true"
        :aria-expanded="isOpen"
        aria-controls="mobile-menu-panel"
        @click="isOpen = true"
      >
        <Menu class="size-24" aria-hidden="true" />
        <span class="sr-only">Открыть меню</span>
      </button>
    </div>

    <Transition name="header-fade">
      <div v-if="isOpen" class="fixed inset-0 z-50 bg-ink-900/40 lg:hidden" @click="isOpen = false" />
    </Transition>

    <Transition name="header-slide">
      <div
        v-if="isOpen"
        id="mobile-menu-panel"
        ref="panelRef"
        role="dialog"
        aria-modal="true"
        aria-label="Меню"
        class="fixed inset-y-0 right-0 z-50 flex w-full max-w-[360px] flex-col gap-32 bg-white p-32 lg:hidden"
        @keydown.esc="isOpen = false"
      >
        <div class="flex items-center justify-between">
          <span class="font-display text-h4 text-ink-900">Меню</span>
          <button type="button" class="grid size-[44px] place-items-center rounded-sm text-ink-900" @click="isOpen = false">
            <X class="size-24" aria-hidden="true" />
            <span class="sr-only">Закрыть меню</span>
          </button>
        </div>
        <nav class="flex flex-col gap-24" aria-label="Мобильная навигация">
          <NuxtLink
            v-for="link in [...leftLinks, ...rightLinks]"
            :key="link.to"
            :to="link.to"
            class="font-display text-h4 text-ink-900"
            @click="isOpen = false"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>
      </div>
    </Transition>
  </header>
</template>

<style scoped>
.header-fade-enter-active,
.header-fade-leave-active {
  transition: opacity 0.2s ease;
}
.header-fade-enter-from,
.header-fade-leave-to {
  opacity: 0;
}

.header-slide-enter-active,
.header-slide-leave-active {
  transition: transform 0.25s ease;
}
.header-slide-enter-from,
.header-slide-leave-to {
  transform: translateX(100%);
}

@media (prefers-reduced-motion: reduce) {
  .header-fade-enter-active,
  .header-fade-leave-active,
  .header-slide-enter-active,
  .header-slide-leave-active {
    transition: none;
  }
}
</style>
