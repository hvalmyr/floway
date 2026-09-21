export interface ImageQualitySettings {
  desktop: number;
  mobile: number;
  avif: number;
}

// Same numbers the hero image used to hardcode before these became
// admin-editable (see page_content keys image_quality_* from migration
// 00041) — kept here as the fallback for as long as the plugin's fetch
// hasn't resolved yet (or the keys are somehow missing).
export const DEFAULT_IMAGE_QUALITY: ImageQualitySettings = {
  desktop: 80,
  mobile: 55,
  avif: 50,
};

/**
 * Shared reactive photo-compression settings, populated once by the
 * `image-quality` plugin from page_content (keys image_quality_desktop/
 * _mobile/_avif) before the app renders. Components read this instead of
 * hardcoding a `quality` number, so an admin edit on /admin/page-content
 * takes effect without a redeploy.
 */
export function useImageQuality() {
  return useState<ImageQualitySettings>("image-quality", () => ({ ...DEFAULT_IMAGE_QUALITY }));
}
