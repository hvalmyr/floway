/**
 * Blocklist sanitizer for admin-uploaded icon SVGs (AppIcon.vue's v-html
 * sink) — parsed as XML so namespaces round-trip correctly, unlike HTML
 * parsing. Strips the handful of constructs that can execute script inside
 * an inline SVG (<script>, <foreignObject>, on* event handlers, javascript:
 * URLs) while leaving every legitimate drawing/structural element alone —
 * a narrow allowlist would need updating for every SVG feature icons
 * legitimately use (gradients, masks, symbols, ...).
 */
const REMOVED_TAGS = new Set(["script", "foreignobject"]);
const DANGEROUS_URL_RE = /^\s*javascript:/i;

function sanitizeSvgElement(el: Element) {
  for (const child of Array.from(el.children)) {
    if (REMOVED_TAGS.has(child.tagName.toLowerCase())) {
      child.remove();
      continue;
    }
    sanitizeSvgElement(child);
  }
  for (const attr of Array.from(el.attributes)) {
    const name = attr.name.toLowerCase();
    if (name.startsWith("on")) {
      el.removeAttribute(attr.name);
      continue;
    }
    if ((name === "href" || name === "xlink:href") && DANGEROUS_URL_RE.test(attr.value)) {
      el.removeAttribute(attr.name);
    }
  }
}

export function sanitizeSvg(svg: string | undefined): string {
  if (!svg) return "";
  if (typeof DOMParser === "undefined") return svg;
  const parsed = new DOMParser().parseFromString(svg, "image/svg+xml");
  if (parsed.querySelector("parsererror")) return "";
  const root = parsed.documentElement;
  if (root.tagName.toLowerCase() !== "svg") return "";
  sanitizeSvgElement(root);
  return new XMLSerializer().serializeToString(root);
}
