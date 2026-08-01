/**
 * Shapes mirror the Go backend's JSON output (backend/internal/model/model.go),
 * which uses explicit camelCase `json` tags — confirmed by reading the backend
 * source, not guessed. Keep these two files in sync until the contract is
 * generated/shared automatically.
 */

export type CourseStatus = "active" | "archived";
export type MasterclassStatus = "active" | "archived";

export interface Course {
  id: number;
  slug: string;
  title: string;
  shortDescription: string;
  fullDescription: string;
  status: CourseStatus;
  coverImage: string;
  gallery: string[];
  sortOrder: number;
}

export interface CourseBlock {
  id: number;
  courseId: number;
  title: string;
  lessonsCount: number;
  hours: number;
  price: number;
  /**
   * TODO: согласовать с бэкенд-командой — поля "старая цена" (для скидок) в
   * backend/internal/model/model.go сейчас нет. Пока это чисто фронтовое
   * опциональное поле, приходит только из моков.
   */
  oldPrice?: number;
  sortOrder: number;
}

export interface Lesson {
  id: number;
  courseBlockId: number;
  number: number;
  title: string;
  topics: string;
  outcomes: string;
  durationHours: number;
}

export interface CourseModule extends CourseBlock {
  /** Intro text shown above the "Учебный план" accordion for this block. */
  description?: string;
  lessons: Lesson[];
}

/** Shape returned by the public GET /api/v1/courses/{slug}/full. */
export interface CourseDetail extends Course {
  modules: CourseModule[];
}

export interface Masterclass {
  id: number;
  slug: string;
  title: string;
  shortDescription: string;
  fullDescription: string;
  endingText: string;
  /**
   * Optional here (unlike the Go model, where these are required columns)
   * because one masterclass in the brief ("Интерьерная композиция") has no
   * confirmed duration/price yet — see the TODO in mocks/masterclasses.ts.
   * A real backend record will always have these populated.
   */
  duration?: string;
  priceGroup?: number;
  priceIndividual?: number;
  priceDescription?: string;
  coverImage: string;
  status: MasterclassStatus;
}

/**
 * Shape returned by the public GET /api/v1/page-content (list, no auth).
 * Generic freeform site copy (Hero text, legal pages, etc.) — see
 * usePageContent(). `value` may contain markdown; render it with
 * <MarkdownContent> where the surrounding markup allows block content.
 */
export type PageContentType = "text" | "image";

export interface PageContent {
  key: string;
  label: string;
  value: string;
  type: PageContentType;
  updatedAt: string;
}

export type FeaturePage = "home" | "masterclasses";

/**
 * Shape returned by the public GET /api/v1/features?page=home|masterclasses
 * (list, no auth, pre-sorted by sortOrder). `icon` is a key into
 * FEATURE_ICONS (constants/feature-icons.ts), not a rendered value.
 */
export interface Feature {
  id: number;
  page: FeaturePage;
  icon: string;
  title: string;
  description: string;
  sortOrder: number;
}

/** Shape returned by the public GET /api/v1/about-items (list, no auth, pre-sorted by sortOrder). */
export interface AboutItem {
  id: number;
  badge: string;
  description: string;
  sortOrder: number;
}

/** Shape returned by the public GET /api/v1/faq (list, no auth, pre-sorted by sortOrder). */
export interface FAQItem {
  id: number;
  question: string;
  answer: string;
  sortOrder: number;
}

/** Shape returned by the public GET /api/v1/teachers (list, no auth). */
export interface Teacher {
  id: number;
  name: string;
  photo: string;
  description: string;
  sortOrder: number;
}

export type BlogPostStatus = "draft" | "published";

/**
 * Shape returned by GET /api/v1/blog-posts?status=published (list) and
 * GET /api/v1/blog-posts/{slug} (single post) — both public, published-only.
 */
export interface BlogPost {
  id: number;
  slug: string;
  title: string;
  coverImage: string;
  category: string;
  tags: string[];
  author: string;
  publishedAt: string | null;
  content: string;
  status: BlogPostStatus;
}

export type ContactMethod = "call" | "telegram" | "whatsapp" | "max";
export type LeadSource = "referral" | "ads" | "internet" | "social" | "maps";
export type LeadRequestType = "course" | "masterclass" | "trial_lesson";

/** POST body for /api/v1/leads (public route — see lead_handler.go). */
export interface ApplicationPayload {
  name: string;
  phone: string;
  email?: string;
  contactMethod: ContactMethod;
  source: LeadSource;
  requestType: LeadRequestType;
  relatedId?: number;
}

export interface Lead extends ApplicationPayload {
  id: number;
  status: "new" | "in_progress" | "closed";
  createdAt: string;
}

/** Normalized error shape used by useApi() so components never touch raw $fetch errors. */
export interface ApiError {
  message: string;
  status?: number;
}
