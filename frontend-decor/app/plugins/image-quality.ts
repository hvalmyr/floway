import {
  DEFAULT_IMAGE_QUALITY,
  useImageQuality,
  type ImageQualitySettings,
} from "~/composables/useImageQuality";

// Same useAsyncData key ("page-content") as usePageContent() — Nuxt
// dedupes calls sharing a key, so pages that also call usePageContent()
// reuse this one fetch instead of hitting the API twice. Runs as an async
// plugin, which Nuxt awaits (on both server and client) before the app
// renders, so every component's first paint already sees the real
// admin-configured numbers — no flash of the defaults, no hydration
// mismatch between SSR output and the client.
const KEY_BY_FIELD: Record<keyof ImageQualitySettings, string> = {
  desktop: "image_quality_desktop",
  mobile: "image_quality_mobile",
  avif: "image_quality_avif",
};

export default defineNuxtPlugin({
  name: "image-quality",
  async setup() {
    const quality = useImageQuality();
    const api = useApi();

    const { data } = await useAsyncData("page-content", () => api.getPageContent());
    const byKey = Object.fromEntries((data.value ?? []).map((item) => [item.key, item.value]));

    const resolved = { ...DEFAULT_IMAGE_QUALITY };
    for (const field of Object.keys(KEY_BY_FIELD) as (keyof ImageQualitySettings)[]) {
      const raw = Number(byKey[KEY_BY_FIELD[field]]);
      if (Number.isFinite(raw) && raw >= 1 && raw <= 100) resolved[field] = raw;
    }
    quality.value = resolved;
  },
});
