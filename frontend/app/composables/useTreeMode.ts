/**
 * Whether the admin has the ambient tree background turned on (the
 * `site_tree_enabled` page_content toggle — see admin/page-content's
 * generic list, and AmbientTreeBackground.vue for the 3D scene itself).
 *
 * `glassClass` is the one thing every "glass over the tree" panel across
 * the site needs to switch on this: with the tree on, the panel stays
 * translucent+blurred so the branch reads through it; with it off there's
 * nothing behind the panel to blur and no point pretending there is, so the
 * panel drops its background/blur entirely and its content just sits as
 * plain text/inputs on the page's own white background (see main.css's
 * `body` rule for that white-vs-transparent swap). The header and the
 * cookie banner don't use this — both are fixed over content that scrolls
 * underneath them, so they keep their glass treatment regardless of the
 * tree. Sections that were already beige to begin with (`bg-surface/55`,
 * not `bg-white/*`) aren't part of this either — they don't change.
 *
 * @example
 * const { glassClass } = await useTreeMode();
 * // in template: :class="glassClass" (replaces a literal
 * // 'bg-white/55 backdrop-blur backdrop-saturate-150')
 */
export async function useTreeMode() {
  const { text } = await usePageContent();
  const treeEnabled = computed(() => text("site_tree_enabled", "false") === "true");
  const glassClass = computed(() =>
    treeEnabled.value ? "bg-white/55 backdrop-blur backdrop-saturate-150" : "",
  );
  return { treeEnabled, glassClass };
}
