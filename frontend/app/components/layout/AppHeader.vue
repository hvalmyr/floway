<script setup lang="ts">
// useFocusTrap actually lives in @vueuse/integrations (wraps the `focus-trap`
// package), not @vueuse/core — the brief named the wrong package; installed
// the correct one instead of hand-rolling a trap.
import { useFocusTrap } from "@vueuse/integrations/useFocusTrap";
import { Menu } from "lucide-vue-next";
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
  <header class="sticky top-0 z-40 h-[64px] border-b-2 border-primary/15 bg-white lg:h-[88px]">
    <div class="container flex h-full items-center justify-between lg:grid lg:grid-cols-3">
      <NuxtLink to="/" class="font-display text-h4 text-primary lg:hidden">ФлоВей</NuxtLink>

      <nav
        class="hidden items-center gap-32 lg:flex lg:justify-self-start"
        aria-label="Основная навигация"
      >
        <NuxtLink
          v-for="link in leftLinks"
          :key="link.to"
          :to="link.to"
          class="border-b-2 border-transparent font-body text-body text-ink hover:text-primary"
          :class="isActive(link.to) ? 'border-primary text-primary' : ''"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <NuxtLink
        to="/"
        class="hidden font-display text-h4 text-primary lg:block lg:justify-self-center"
      >
        ФлоВей
      </NuxtLink>

      <nav
        class="hidden items-center gap-32 lg:flex lg:justify-self-end"
        aria-label="Дополнительная навигация"
      >
        <NuxtLink
          v-for="link in rightLinks"
          :key="link.to"
          :to="link.to"
          class="border-b-2 border-transparent font-body text-body text-ink hover:text-primary"
          :class="isActive(link.to) ? 'border-primary text-primary' : ''"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <button
        type="button"
        class="grid size-[44px] place-items-center rounded-sm text-ink lg:hidden"
        aria-haspopup="true"
        :aria-expanded="isOpen"
        aria-controls="mobile-menu-panel"
        @click="isOpen = !isOpen"
      >
        <IconClose v-if="isOpen" class="size-24" aria-hidden="true" />
        <Menu v-else class="size-24" aria-hidden="true" />
        <span class="sr-only">{{ isOpen ? "Закрыть меню" : "Открыть меню" }}</span>
      </button>
    </div>

    <!-- Мобильное меню — выпадает сверху (dropdown), а не выезжает сбоку. -->
    <Transition name="header-fade">
      <div
        v-if="isOpen"
        class="fixed inset-x-0 bottom-0 top-[64px] z-50 bg-ink/40 lg:hidden"
        @click="isOpen = false"
      />
    </Transition>

    <Transition name="header-drop">
      <div
        v-if="isOpen"
        id="mobile-menu-panel"
        ref="panelRef"
        role="dialog"
        aria-modal="true"
        aria-label="Меню"
        class="fixed inset-x-0 top-[64px] z-50 border-b-2 border-primary/15 bg-white p-24 shadow-lg lg:hidden"
        @keydown.esc="isOpen = false"
      >
        <nav class="flex flex-col gap-8" aria-label="Мобильная навигация">
          <NuxtLink
            v-for="link in [...leftLinks, ...rightLinks]"
            :key="link.to"
            :to="link.to"
            class="rounded-sm px-12 py-12 font-display text-h4 text-ink hover:bg-surface hover:text-primary"
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

.header-drop-enter-active,
.header-drop-leave-active {
  transition: transform 0.25s ease;
}
.header-drop-enter-from,
.header-drop-leave-to {
  transform: translateY(-100%);
}

@media (prefers-reduced-motion: reduce) {
  .header-fade-enter-active,
  .header-fade-leave-active,
  .header-drop-enter-active,
  .header-drop-leave-active {
    transition: none;
  }
}
</style>
