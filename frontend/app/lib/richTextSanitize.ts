/**
 * Whitelist sanitizer for AdminRichTextEditor.vue's contenteditable HTML —
 * runs in the browser only (uses a <template>, whose content is an inert
 * document: nothing renders, no image loads, no script runs while we walk
 * it). Anything outside the allowed tag/attribute set is stripped; disallowed
 * elements are unwrapped (their children survive, the wrapping tag doesn't)
 * rather than deleted outright, so pasted rich text degrades to plain
 * formatting instead of losing content.
 */
const ALLOWED_TAGS = new Set([
  "P",
  "H3",
  "BLOCKQUOTE",
  "A",
  "STRONG",
  "B",
  "EM",
  "I",
  "IMG",
  "FIGURE",
  "BR",
]);

const ALLOWED_URL_SCHEMES = new Set(["http:", "https:", "mailto:"]);

function isSafeUrl(value: string, base: string): boolean {
  try {
    const url = new URL(value, base);
    return ALLOWED_URL_SCHEMES.has(url.protocol) || value.startsWith("/");
  } catch {
    return false;
  }
}

function sanitizeElement(el: Element, baseUrl: string) {
  for (const attr of Array.from(el.attributes)) {
    let keep = false;
    if (el.tagName === "A" && attr.name === "href") {
      keep = isSafeUrl(attr.value, baseUrl);
    } else if (el.tagName === "IMG" && (attr.name === "src" || attr.name === "alt")) {
      keep = attr.name === "alt" || isSafeUrl(attr.value, baseUrl);
    }
    if (!keep) el.removeAttribute(attr.name);
  }
  if (el.tagName === "A" && el.hasAttribute("href")) {
    el.setAttribute("target", "_blank");
    el.setAttribute("rel", "noopener noreferrer");
  }
}

function walk(node: Node, baseUrl: string) {
  for (const child of Array.from(node.childNodes)) {
    if (child.nodeType === Node.TEXT_NODE) continue;
    if (child.nodeType !== Node.ELEMENT_NODE) {
      node.removeChild(child);
      continue;
    }
    const el = child as Element;
    walk(el, baseUrl);
    if (!ALLOWED_TAGS.has(el.tagName)) {
      while (el.firstChild) node.insertBefore(el.firstChild, el);
      node.removeChild(el);
      continue;
    }
    sanitizeElement(el, baseUrl);
  }
}

export function sanitizeRichTextHtml(html: string): string {
  if (typeof document === "undefined") return html;
  const template = document.createElement("template");
  template.innerHTML = html;
  walk(template.content, window.location.origin);
  return template.innerHTML;
}
