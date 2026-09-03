<script setup lang="ts">
import { normalizeClass } from "vue";
import { featureIconComponent } from "~/constants/feature-icons";

/**
 * Resolves a Feature/PageContent icon value to something renderable —
 * either a built-in Vue icon component (a bare FEATURE_ICONS key, e.g.
 * "gift") or an admin-uploaded SVG ("icon:<id>", looked up via
 * useIconLibrary). One component so every call site (FeatureGrid,
 * GiftCertificateMarquee, ...) renders either kind the same way instead of
 * branching on the value shape itself.
 *
 * `class`/other attrs are meant to land on the actual rendered <svg> either
 * way — for the built-in branch that's automatic (Vue's normal attrs
 * fallthrough onto the icon component's own root element), but v-html
 * content has no component boundary for Vue to attach fallthrough attrs to,
 * so the uploaded-icon branch copies them onto the injected <svg> by hand
 * once it's actually in the DOM (see the watcher below).
 *
 * @example
 * <AppIcon :icon="feature.icon" class="h-40 w-auto text-primary" />
 */
defineOptions({ inheritAttrs: false });

const props = defineProps<{ icon: string | undefined }>();
const attrs = useAttrs();

const { findIcon } = useIconLibrary();

const builtinComponent = computed(() => {
  if (!props.icon || props.icon.startsWith("icon:")) return undefined;
  return featureIconComponent(props.icon);
});

const uploadedSvg = computed(() => {
  if (!props.icon?.startsWith("icon:")) return undefined;
  const id = Number(props.icon.slice("icon:".length));
  return findIcon(id)?.svg;
});

const hostEl = ref<HTMLElement | null>(null);

watch(
  [uploadedSvg, () => attrs.class],
  async () => {
    if (!uploadedSvg.value) return;
    await nextTick();
    const svg = hostEl.value?.querySelector("svg");
    if (!svg) return;
    svg.removeAttribute("class");
    const className = normalizeClass(attrs.class);
    if (className) svg.setAttribute("class", className);
  },
  { immediate: true },
);
</script>

<template>
  <component :is="builtinComponent" v-if="builtinComponent" v-bind="attrs" aria-hidden="true" />
  <span v-else-if="uploadedSvg" ref="hostEl" aria-hidden="true" v-html="uploadedSvg" />
</template>
