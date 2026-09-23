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
  "H2",
  // H3 stays whitelisted for backward compat with posts saved before the
  // editor's single heading level was renamed H3 -> H2 (matches the H1 page
  // title semantically, and the toolbar's Heading2 icon) — without this,
  // re-editing an old post would silently unwrap its existing headings.
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

/**
 * Resolves a URL typed or pasted by an admin into something that actually
 * points where they meant. A bare domain like "vk.com/floway" (no scheme)
 * has no special meaning to the URL parser — resolved against the site's
 * own origin it turns into a same-site relative path, which parses as
 * "safe" but silently points at the wrong place once rendered. Site-relative
 * links ("/blog", "#section") are the one case where that resolution is
 * actually what's wanted, so those pass through unchanged; anything else
 * without an explicit scheme is assumed to be an external link and gets
 * "https://" prepended.
 */
function normalizeUrl(value: string, base: string): string | null {
  const trimmed = value.trim();
  if (trimmed === "") return null;
  if (trimmed.startsWith("/") || trimmed.startsWith("#")) {
    try {
      new URL(trimmed, base);
      return trimmed;
    } catch {
      return null;
    }
  }
  const hasScheme = /^[a-z][a-z0-9+.-]*:/i.test(trimmed);
  const candidate = hasScheme ? trimmed : `https://${trimmed}`;
  try {
    const url = new URL(candidate, base);
    return ALLOWED_URL_SCHEMES.has(url.protocol) ? candidate : null;
  } catch {
    return null;
  }
}

function sanitizeElement(el: Element, baseUrl: string) {
  for (const attr of Array.from(el.attributes)) {
    let newValue: string | null = null;
    if (el.tagName === "A" && attr.name === "href") {
      newValue = normalizeUrl(attr.value, baseUrl);
    } else if (el.tagName === "IMG" && attr.name === "src") {
      newValue = normalizeUrl(attr.value, baseUrl);
    } else if (el.tagName === "IMG" && attr.name === "alt") {
      newValue = attr.value;
    }
    if (newValue === null) {
      el.removeAttribute(attr.name);
    } else if (newValue !== attr.value) {
      el.setAttribute(attr.name, newValue);
    }
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

export function normalizeLinkUrl(value: string): string | null {
  if (typeof window === "undefined") return value;
  return normalizeUrl(value, window.location.origin);
}

export function sanitizeRichTextHtml(html: string): string {
  if (typeof DOMParser === "undefined") return html;
  // DOMParser, not <template>.innerHTML — this function backs the "default"
  // Trusted Types policy (see trusted-types.client.ts), so it can't itself
  // assign a raw string to .innerHTML without recursing into that same
  // policy. DOMParser.parseFromString isn't a Trusted Types sink (its
  // output document is inert, same as <template>.content), so it sidesteps
  // that without changing behavior for this tag set (no table/list-context
  // parsing quirks to worry about — ALLOWED_TAGS are all valid direct
  // children of <body>).
  const parsed = new DOMParser().parseFromString(html, "text/html");
  walk(parsed.body, window.location.origin);
  return parsed.body.innerHTML;
}
