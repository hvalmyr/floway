/** Picks readable black or white text for a "#rrggbb" background, using the
 * standard YIQ perceived-brightness approximation (good enough for chip
 * text — not trying to be a full WCAG contrast calculator). Malformed
 * input falls back to black, matching how dark text reads on the default
 * neutral-gray tag color. */
export function readableTextColor(hexColor: string): "#000000" | "#ffffff" {
  const hex = hexColor.replace("#", "");
  if (!/^[0-9a-fA-F]{6}$/.test(hex)) return "#000000";

  const r = parseInt(hex.slice(0, 2), 16);
  const g = parseInt(hex.slice(2, 4), 16);
  const b = parseInt(hex.slice(4, 6), 16);
  const brightness = (r * 299 + g * 587 + b * 114) / 1000;
  return brightness >= 128 ? "#000000" : "#ffffff";
}
