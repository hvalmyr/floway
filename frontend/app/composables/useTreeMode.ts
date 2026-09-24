/**
 * Whether the admin has the ambient tree background turned on (the
 * `site_tree_enabled` page_content toggle — see admin/page-content's
 * generic list, and AmbientTreeBackground.vue for the 3D scene itself).
 *
 * `glassClass` is the one thing every "glass over the tree" panel across
 * the site needs to switch on this: with the tree on, the panel stays
 * translucent+blurred so the branch reads through it; with it off there's
 * nothing behind the panel to blur, so it becomes a plain solid beige
 * (bg-surface) card instead — see main.css's `body` rule for the matching
 * page-background swap (white without the tree, transparent with it).
 * Sections that were already beige to begin with (`bg-surface/55`, not
 * `bg-white/*`) aren't part of this — they don't change either way.
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
    treeEnabled.value ? "bg-white/55 backdrop-blur backdrop-saturate-150" : "bg-surface",
  );
  return { treeEnabled, glassClass };
}
