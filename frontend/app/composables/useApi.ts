import { FetchError } from "ofetch";
import { mockGetMasterClasses } from "~/mocks/masterclasses";
import type {
  AboutItem,
  ApiError,
  ApplicationPayload,
  BlogPost,
  CourseSectionWithCourses,
  CourseWithBlocks,
  FAQItem,
  Feature,
  FeaturePage,
  GalleryPhoto,
  Icon,
  Lead,
  Masterclass,
  PageContent,
  Teacher,
} from "~/types/api";

function toApiError(err: unknown): ApiError {
  if (err instanceof FetchError) {
    return {
      message:
        (err.data as { message?: string } | undefined)?.message ??
        "Не удалось выполнить запрос к серверу",
      status: err.statusCode,
    };
  }
  return { message: "Не удалось выполнить запрос к серверу" };
}

/**
 * Single access point for public-facing (non-admin) API calls. Wraps
 * useApiClient() with named methods and normalized errors, so components
 * never touch $fetch or raw error shapes directly.
 *
 * All methods hit the real backend by default. Set
 * NUXT_PUBLIC_USE_MOCKS=true to fall back to app/mocks/* instead — same
 * shape, no component changes needed — useful for working on layout/markup
 * without a backend running.
 *
 * @example
 * const api = useApi();
 * const { data: sections } = await useAsyncData('course-sections', () => api.getCourseSections());
 */
export function useApi() {
  const config = useRuntimeConfig();
  const client = useApiClient();
  const useMocks = config.public.useMocks;

  /**
   * GET /api/v1/course-sections/full — public, no auth. Every course
   * section with its courses and each course's blocks (no lesson text) —
   * the homepage courses block. No mocks fallback: no design brief existed
   * for this shape before it was built.
   */
  async function getCourseSections(): Promise<CourseSectionWithCourses[]> {
    try {
      return await client<CourseSectionWithCourses[]>("/api/v1/course-sections/full");
    } catch (err) {
      throw toApiError(err);
    }
  }

  async function getCourse(slug: string): Promise<CourseWithBlocks | null> {
    try {
      return await client<CourseWithBlocks>(`/api/v1/courses/${slug}/full`);
    } catch (err) {
      throw toApiError(err);
    }
  }

  async function getMasterClasses(): Promise<Masterclass[]> {
    if (useMocks) return mockGetMasterClasses();
    try {
      return await client<Masterclass[]>("/api/v1/masterclasses");
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/teachers — public, no auth (teacher_handler.go's list route
   * has no admin middleware). No mocks fallback: there's no design brief for
   * this data shape (the homepage section was hardcoded name+accent-color
   * placeholders), same reasoning as blog posts below.
   */
  async function getTeachers(): Promise<Teacher[]> {
    try {
      return await client<Teacher[]>("/api/v1/teachers");
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/gallery-photos — public, no auth. Slides for the homepage
   * carousel between "Преимущества" and "О школе", pre-sorted by sortOrder.
   */
  async function getGalleryPhotos(): Promise<GalleryPhoto[]> {
    try {
      return await client<GalleryPhoto[]>("/api/v1/gallery-photos");
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/faq — public, no auth (faq_handler.go's list route has no
   * admin middleware). Already sorted by sortOrder on the backend. No mocks
   * fallback, same reasoning as teachers above.
   */
  async function getFAQItems(): Promise<FAQItem[]> {
    try {
      return await client<FAQItem[]>("/api/v1/faq");
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/page-content — public, no auth. Generic freeform site copy
   * (Hero text, legal pages, etc.) — see usePageContent() for the
   * key-lookup helper components actually use.
   */
  async function getPageContent(): Promise<PageContent[]> {
    try {
      return await client<PageContent[]>("/api/v1/page-content");
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/features?page=home|masterclasses — public, no auth.
   * Icon+title+description cards (feature_handler.go's list route has no
   * admin middleware). Already sorted by sortOrder on the backend.
   */
  async function getFeatures(page: FeaturePage): Promise<Feature[]> {
    try {
      return await client<Feature[]>("/api/v1/features", { query: { page } });
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/about-items — public, no auth. Badge+description cards for
   * the homepage "О школе" section. Already sorted by sortOrder.
   */
  async function getAboutItems(): Promise<AboutItem[]> {
    try {
      return await client<AboutItem[]>("/api/v1/about-items");
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/icons — public, no auth (icon_handler.go's list route has
   * no admin middleware, same as features/gallery-photos). The uploaded
   * icon library — see the Icon type doc comment for how a Feature/
   * PageContent icon value references one of these.
   */
  async function getIcons(): Promise<Icon[]> {
    try {
      return await client<Icon[]>("/api/v1/icons");
    } catch (err) {
      throw toApiError(err);
    }
  }

  /**
   * GET /api/v1/blog-posts?status=published — draft posts are never
   * returned (see blog_post_handler.go). No mocks fallback: there's no
   * design brief for the blog yet, so this always hits the real backend.
   */
  async function getBlogPosts(): Promise<BlogPost[]> {
    try {
      return await client<BlogPost[]>("/api/v1/blog-posts", { query: { status: "published" } });
    } catch (err) {
      throw toApiError(err);
    }
  }

  async function getBlogPost(slug: string): Promise<BlogPost | null> {
    try {
      return await client<BlogPost>(`/api/v1/blog-posts/${slug}`);
    } catch (err) {
      throw toApiError(err);
    }
  }

  /** POST /api/v1/leads — this route is real and public (see lead_handler.go). */
  async function submitApplication(payload: ApplicationPayload): Promise<Lead> {
    try {
      return await client<Lead>("/api/v1/leads", { method: "POST", body: payload });
    } catch (err) {
      throw toApiError(err);
    }
  }

  return {
    getCourseSections,
    getCourse,
    getMasterClasses,
    getTeachers,
    getGalleryPhotos,
    getFAQItems,
    getFeatures,
    getAboutItems,
    getIcons,
    getPageContent,
    getBlogPosts,
    getBlogPost,
    submitApplication,
  };
}
