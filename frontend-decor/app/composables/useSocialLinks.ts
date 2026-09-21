import { socialLinks as fallbackSocialLinks, type SocialLink } from "~/constants/contact-info";

/**
 * Social link icons (footer + /contacts) — admin-editable via
 * /admin/social-links, same list backing both. Falls back to the
 * hardcoded contact-info.ts values if the API call fails or the table is
 * empty, same pattern as usePageContent()'s text() fallback.
 */
export async function useSocialLinks() {
  const api = useApi();
  const { data } = await useAsyncData("social-links", () => api.getSocialLinks());

  const socialLinks = computed<SocialLink[]>(() => {
    if (!data.value || data.value.length === 0) return fallbackSocialLinks;
    return data.value
      .slice()
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((item) => ({
        label: item.label,
        href: item.href,
        disclaimer: item.disclaimer || undefined,
      }));
  });

  return { socialLinks };
}
