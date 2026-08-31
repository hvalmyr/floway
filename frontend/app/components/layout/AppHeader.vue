<script setup lang="ts">
// useFocusTrap actually lives in @vueuse/integrations (wraps the `focus-trap`
// package), not @vueuse/core — the brief named the wrong package; installed
// the correct one instead of hand-rolling a trap.
import { useFocusTrap } from "@vueuse/integrations/useFocusTrap";
import { Menu, X } from "lucide-vue-next";
import { nextTick, ref, watch } from "vue";

const route = useRoute();
const isOpen = ref(false);
const panelRef = ref<HTMLElement | null>(null);

const { activate, deactivate } = useFocusTrap(panelRef, {
  immediate: false,
  escapeDeactivates: true,
  // Без этого клик по крестику (он вне panelRef — живёт в шапке, не в самой
  // выпадающей панели) перехватывался трапом: фокус возвращался внутрь
  // панели, и закрыть меню кликом было невозможно.
  allowOutsideClick: true,
  // Пока панель уезжает transition'ом, focus-trap's MutationObserver может
  // застать контейнер без единого tabbable-узла в процессе анмаунта и
  // выбросить "must have at least one container with at least one tabbable
  // node" — fallbackFocus даёт ему всегда валидную цель и убирает throw.
  fallbackFocus: "body",
  onDeactivate: () => {
    isOpen.value = false;
  },
});

// The background page is deliberately left scrollable while the menu is
// open (no body scroll lock) — the panel/backdrop are `fixed`, so they stay
// pinned over the viewport regardless of scroll position anyway, and
// locking scroll also hid the browser's own scrollbar, which read as a bug
// ("the scroll button disappears") rather than a feature.
watch(isOpen, async (open) => {
  if (open) {
    await nextTick();
    activate();
  } else {
    deactivate();
  }
});

const leftLinks = [
  { to: "/#courses", label: "Курсы" },
  { to: "/masterclasses", label: "Мастер-классы" },
];
// TODO: вернуть "Блог" в навигацию, когда раздел будет доделан — сайт
// запускают раньше блога, ссылка на недоделанный раздел временно убрана.
const rightLinks = [{ to: "/contacts", label: "Контакты" }];

// "Курсы" ведёт на якорь на главной ("/#courses") — на самой главной странице
// (без явного скролла к якорю) её всё равно считаем открытой страницей,
// иначе на "/" ни одна ссылка никогда не подсвечивалась бы активной.
function isActive(to: string) {
  if (to === "/#courses") return route.path === "/";
  return !to.includes("#") && to !== "/" && route.path.startsWith(to);
}
</script>

<template>
  <header
    class="sticky top-0 z-40 h-[64px] bg-white/70 backdrop-blur backdrop-saturate-150 lg:h-[88px]"
  >
    <div class="container flex h-full items-center justify-between lg:grid lg:grid-cols-3">
      <NuxtLink to="/" class="font-display text-h2 text-primary lg:hidden">фловей</NuxtLink>

      <nav
        class="hidden items-center gap-32 lg:flex lg:justify-self-start"
        aria-label="Основная навигация"
      >
        <NuxtLink
          v-for="link in leftLinks"
          :key="link.to"
          :to="link.to"
          class="font-display text-button text-ink hover:text-primary"
          :class="isActive(link.to) ? 'text-primary' : ''"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <NuxtLink
        to="/"
        class="hidden font-display text-h2 text-primary lg:block lg:justify-self-center"
      >
        фловей
      </NuxtLink>

      <nav
        class="hidden items-center gap-32 lg:flex lg:justify-self-end"
        aria-label="Дополнительная навигация"
      >
        <NuxtLink
          v-for="link in rightLinks"
          :key="link.to"
          :to="link.to"
          class="font-display text-button text-ink hover:text-primary"
          :class="isActive(link.to) ? 'text-primary' : ''"
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
        <X v-if="isOpen" class="size-24" aria-hidden="true" />
        <Menu v-else class="size-24" aria-hidden="true" />
        <span class="sr-only">{{ isOpen ? "Закрыть меню" : "Открыть меню" }}</span>
      </button>
    </div>
  </header>

  <!-- Teleported out of <header> rather than left as its children: a
  `position:sticky` header establishes its own stacking context, and WITHIN
  that context any positioned descendant with an explicit z-index — no
  matter the number — paints above the header's own background/content
  (CSS stacking-order rule, not a z-index magnitude issue). z-30 here was
  never actually being compared against the header's z-40; it was being
  compared against the header's own painted box, which it always wins.
  Teleporting to <body> makes it a true sibling of <header> in the root
  stacking context, where 30 < 40 finally means what it looks like it
  means — the header stays visually on top and the menu reads as sliding
  out from under it. -->
  <Teleport to="body">
    <Transition name="header-fade">
      <div
        v-if="isOpen"
        class="fixed inset-x-0 bottom-0 top-[64px] z-30 bg-ink/40 lg:hidden"
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
        class="container fixed inset-x-0 top-[64px] z-30 bg-white py-24 shadow-lg lg:hidden"
        @keydown.esc="isOpen = false"
      >
        <!-- -mx-12 компенсирует px-12 ссылки, чтобы текст выравнивался по
        тому же левому краю, что и логотип в шапке (тот же .container),
        а не проваливался глубже из-за собственного паддинга ссылки. -->
        <nav class="-mx-12 flex flex-col gap-8" aria-label="Мобильная навигация">
          <NuxtLink
            v-for="link in [...leftLinks, ...rightLinks]"
            :key="link.to"
            :to="link.to"
            class="rounded-sm px-12 py-12 font-display text-h4 text-ink hover:bg-surface hover:text-primary"
            :class="isActive(link.to) ? 'bg-surface text-primary' : ''"
            @click="isOpen = false"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>
      </div>
    </Transition>
  </Teleport>
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
