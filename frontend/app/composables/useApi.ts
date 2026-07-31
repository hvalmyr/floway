import { FetchError } from "ofetch";
import { mockGetCourse, mockGetCourses } from "~/mocks/courses";
import { mockGetMasterClass, mockGetMasterClasses } from "~/mocks/masterclasses";
import type { ApiError, ApplicationPayload, Course, CourseDetail, Lead, Masterclass } from "~/types/api";

function toApiError(err: unknown): ApiError {
  if (err instanceof FetchError) {
    return {
      message: (err.data as { message?: string } | undefined)?.message ?? "Не удалось выполнить запрос к серверу",
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
 * While the backend doesn't yet expose the routes these methods need (see
 * the TODOs on getCourse/getMasterClass below), runtimeConfig.public.useMocks
 * makes every method fall back to app/mocks/* — same shape, no component
 * changes needed once the real routes exist.
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

  /**
   * TODO: бэкенд пока не отдаёт курс по slug со вложенными блоками/занятиями
   * (см. TODO на CourseDetail в ~/types/api.ts) — нужен отдельный публичный
   * эндпоинт, например GET /api/v1/courses/{slug}/full. Пока используется мок.
   */
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

  /**
   * TODO: бэкенд пока не отдаёт публичный GET одного мастер-класса по slug
   * (только список и CRUD по числовому id для админки, см.
   * backend/internal/httpserver/masterclass_handler.go). Пока используется мок.
   */
  async function getMasterClass(slug: string): Promise<Masterclass | null> {
    if (useMocks) return mockGetMasterClass(slug);
    try {
      return await client<Masterclass>(`/api/v1/masterclasses/${slug}`);
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

  return { getCourses, getCourse, getMasterClasses, getMasterClass, submitApplication };
}
