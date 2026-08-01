import { FetchError } from "ofetch";
import { mockGetCourse, mockGetCourses } from "~/mocks/courses";
import { mockGetMasterClass, mockGetMasterClasses } from "~/mocks/masterclasses";
import type {
  ApiError,
  ApplicationPayload,
  BlogPost,
  Course,
  CourseDetail,
  FAQItem,
  Lead,
  Masterclass,
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
 * const { data: courses } = await useAsyncData('courses', () => api.getCourses());
 */
export function useApi() {
  const config = useRuntimeConfig();
  const client = useApiClient();
  const useMocks = config.public.useMocks;

  async function getCourses(): Promise<Course[]> {
    if (useMocks) return mockGetCourses();
    try {
      return await client<Course[]>("/api/v1/courses");
    } catch (err) {
      throw toApiError(err);
    }
  }

  async function getCourse(slug: string): Promise<CourseDetail | null> {
    if (useMocks) return mockGetCourse(slug);
    try {
      return await client<CourseDetail>(`/api/v1/courses/${slug}/full`);
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

  async function getMasterClass(slug: string): Promise<Masterclass | null> {
    if (useMocks) return mockGetMasterClass(slug);
    try {
      return await client<Masterclass>(`/api/v1/masterclasses/${slug}`);
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
    getCourses,
    getCourse,
    getMasterClasses,
    getMasterClass,
    getTeachers,
    getFAQItems,
    getBlogPosts,
    getBlogPost,
    submitApplication,
  };
}
