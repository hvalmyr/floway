/**
 * Shapes mirror the Go backend's JSON output (backend/internal/model/model.go),
 * which uses explicit camelCase `json` tags — confirmed by reading the backend
 * source, not guessed. Keep these two files in sync until the contract is
 * generated/shared automatically.
 */

export type ImportMode = "merge" | "replace";

/**
 * Response from POST /api/v1/admin/content/import. `counts` keys match
 * SiteContent's own field names (courseSections, courses, ...) — the number
 * of rows written for that entity. `pageContentSkipped` lists page-content
 * keys from the import that don't exist on this server (that entity has no
 * create path — see backend/internal/service/content_export_service.go).
 */
export interface ImportResult {
  mode: ImportMode;
  counts: Record<string, number>;
  pageContentSkipped?: string[];
}

export type MasterclassStatus = "active" | "archived";

/** Background/text color pair for a course block's homepage card — see CourseCard.vue. */
export type CourseBlockDisplayStyle = "blue-beige" | "brown-beige" | "beige-blue" | "beige-brown";

export interface CourseSection {
  id: number;
  heading: string;
  description: string;
  sortOrder: number;
  /** Hides the whole section (and its courses) from the public site without deleting it. */
  visible: boolean;
}

/**
 * coverImage/lessonCount/timeLength/price/displayStyle are this course's OWN
 * card — splitting into CourseBlocks is optional. A course with zero blocks
 * renders as one homepage card built from these fields; a course with one
 * or more blocks renders one card per block instead, and these go unused.
 */
export interface Course {
  id: number;
  sectionId: number;
  slug: string;
  name: string;
  description: string;
  coverImage: string;
  lessonCount: string;
  timeLength: string;
  price: string;
  displayStyle: CourseBlockDisplayStyle;
  sortOrder: number;
  /** Hides this course from the public site (its section listing and its own page) without deleting it. */
  visible: boolean;
  /** Collapses a multi-block course into one homepage card (this course's own fields) instead of one card per block. */
  singleCard: boolean;
  /** Heading shown above this course's FAQ block (CourseFAQItem below), on the course page after the apply form. */
  faqTitle: string;
  /** Intro text shown between the FAQ heading and its items. */
  faqDescription: string;
  /** Shows/hides the whole FAQ block (title, description and items) without deleting anything. */
  faqVisible: boolean;
}

/**
 * One Q&A pair in a single course's FAQ block — distinct from the global,
 * homepage-only FAQItem (no course concept) fetched via getFAQItems().
 */
export interface CourseFAQItem {
  id: number;
  courseId: number;
  question: string;
  answer: string;
  sortOrder: number;
}

/** Pages with their own FAQ block, managed under Редактирование → Компоненты → FAQ. */
export type FaqPage = "masterclasses" | "gift_certificate";

/** One Q&A pair in a page-scoped FAQ block — the non-course counterpart of CourseFAQItem. */
export interface PageFaqItem {
  id: number;
  page: FaqPage;
  question: string;
  answer: string;
  sortOrder: number;
}

/** Shape returned by the public GET /api/v1/page-faq/{page}. */
export interface PageFaq {
  page: FaqPage;
  title: string;
  description: string;
  visible: boolean;
  items: PageFaqItem[];
}

/**
 * A course's blocks carry their own lesson list directly (via
 * CourseBlockWithLessons below) — there's no separate "curriculum" level.
 * lessonCount/timeLength/price are free text ("7 занятий", "30 часов",
 * "38 500 ₽", or even a "было/стало" discount narrative) rather than
 * numbers, so the admin can phrase them however the site copy needs.
 */
export interface CourseBlock {
  id: number;
  courseId: number;
  blockName: string;
  /** Intro text shown above this block's lesson list. */
  description: string;
  blockCover: string;
  lessonCount: string;
  timeLength: string;
  price: string;
  displayStyle: CourseBlockDisplayStyle;
  sortOrder: number;
  /** Hides this block's homepage card (and drops it from the course page) without deleting it. */
  visible: boolean;
}

// A lesson belongs to exactly one parent — courseBlockId (a course split
// into blocks) or courseId (a course with no blocks, editing its lessons
// directly) — never both.
export interface Lesson {
  id: number;
  courseBlockId?: number;
  courseId?: number;
  name: string;
  description: string;
  sortOrder: number;
}

export interface CourseBlockWithLessons extends CourseBlock {
  lessons: Lesson[];
}

/** Shape returned by the public GET /api/v1/courses/{slug}/full. */
export interface CourseWithBlocks extends Course {
  blockCount: number;
  blocks: CourseBlockWithLessons[];
  /** Always populated regardless of faqVisible — gate rendering on that flag yourself. */
  faqItems: CourseFAQItem[];
}

/** Homepage listing shape — blocks without lesson text. */
export interface CourseSummary extends Course {
  blockCount: number;
  blocks: CourseBlock[];
}

/** Shape returned by the public GET /api/v1/course-sections/full. */
export interface CourseSectionWithCourses extends CourseSection {
  length: number;
  courses: CourseSummary[];
}

export interface Masterclass {
  id: number;
  slug: string;
  title: string;
  description: string;
  /** Optional second paragraph — not every masterclass needs one. */
  description2: string;
  /** Optional closing note. */
  endingText: string;
  /**
   * Optional here (unlike the Go model, where these are required columns)
   * because one masterclass in the brief ("Интерьерная композиция") has no
   * confirmed duration/price yet — see the TODO in mocks/masterclasses.ts.
   * A real backend record will always have these populated.
   */
  duration?: string;
  /** Free text ("3000₽", "3000₽ или 4500₽", "от 2500₽") — not always one number. */
  price?: string;
  coverImage: string;
  status: MasterclassStatus;
}

/**
 * Shape returned by the public GET /api/v1/page-content (list, no auth).
 * Generic freeform site copy (Hero text, legal pages, etc.) — see
 * usePageContent(). `value` may contain markdown; render it with
 * <MarkdownContent> where the surrounding markup allows block content.
 */
export type PageContentType = "text" | "image" | "icon" | "number";

export interface PageContent {
  key: string;
  label: string;
  value: string;
  type: PageContentType;
  updatedAt: string;
}

export type FeaturePage = "home" | "masterclasses" | "gift_certificate";

/**
 * An icon value (Feature.icon, a page_content row of type "icon") is either
 * a bare key into FEATURE_ICONS (constants/feature-icons.ts, e.g. "gift")
 * or "icon:<id>" referencing an uploaded Icon below — see AppIcon.vue,
 * which resolves either form to something renderable.
 */
export interface Icon {
  id: number;
  name: string;
  svg: string;
  createdAt: string;
}

/**
 * Shape returned by the public GET /api/v1/features?page=home|masterclasses
 * (list, no auth, pre-sorted by sortOrder). `icon` is not a rendered value —
 * see the Icon doc comment above for the two forms it can take.
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

/** Shape returned by the public GET /api/v1/social-links (list, no auth, pre-sorted by sortOrder). */
export interface SocialLink {
  id: number;
  label: string;
  href: string;
  /** Legally required disclaimer text for some platforms in Russia (e.g. Meta-owned apps) — empty when none applies. */
  disclaimer: string;
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

/**
 * One slide of the homepage's vertical-photo carousel (between
 * "Преимущества" and "О школе"). Shape returned by the public
 * GET /api/v1/gallery-photos (list, no auth, pre-sorted by sortOrder).
 */
export interface GalleryPhoto {
  id: number;
  image: string;
  sortOrder: number;
}

/**
 * One slide of the gift-certificates page's photo carousel, shown right
 * after the hero block. Shape returned by the public
 * GET /api/v1/gift-certificate-carousel-photos (list, no auth, pre-sorted by
 * sortOrder).
 */
export interface GiftCertificateCarouselPhoto {
  id: number;
  image: string;
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
  metaTitle: string;
  metaDescription: string;
  coverImage: string;
  /** Same 4-value background/text color pair as CourseBlockDisplayStyle,
   * reused for the card on /blog instead of a separate blog-only palette. */
  displayStyle: CourseBlockDisplayStyle;
  category: string;
  tags: string[];
  author: string;
  publishedAt: string | null;
  content: string;
  status: BlogPostStatus;
}

/**
 * One slide of a thank-you page variant's optional photo carousel. Shape
 * returned as part of the public GET /api/v1/thank-you-pages/{variant}
 * (ThankYouPage.photos), pre-sorted by sortOrder.
 */
export interface ThankYouPagePhoto {
  id: number;
  variant: ThankYouPageVariant;
  image: string;
  sortOrder: number;
}

/**
 * One Q&A pair in a thank-you page variant's optional mini-FAQ. Shape
 * returned as part of the public GET /api/v1/thank-you-pages/{variant}
 * (ThankYouPage.faqItems), pre-sorted by sortOrder.
 */
export interface ThankYouPageFaqItem {
  id: number;
  variant: ThankYouPageVariant;
  question: string;
  answer: string;
  sortOrder: number;
}

/**
 * Shape returned by the public GET /api/v1/thank-you-pages/{variant} — the
 * content ApplyForm.vue navigates to after a successful submission. The
 * messenger icons and social-network links shown alongside this content are
 * NOT part of this shape — they're the site's existing global contact
 * channels (usePageContent's contact_*_url keys, same ones AppFooter.vue
 * renders) and useSocialLinks(); showMessengers/showSocialLinks only toggle
 * whether that shared, already-editable content appears for this variant.
 */
export interface ThankYouPage {
  variant: ThankYouPageVariant;
  title: string;
  subtitle: string;
  description: string;
  /** Square photo for the page's hero block (title/description/CTAs left,
   * 1:1 photo right) — falls back to a placeholder box when empty. */
  heroImage: string;
  showMessengers: boolean;
  showSocialLinks: boolean;
  showBlogLink: boolean;
  blogLinkText: string;
  blogLinkUrl: string;
  showCarousel: boolean;
  showFaq: boolean;
  showCommunity: boolean;
  communityText: string;
  communityUrl: string;
  photos: ThankYouPagePhoto[];
  faqItems: ThankYouPageFaqItem[];
}

export type ContactMethod = "call" | "telegram" | "whatsapp" | "max";
export type LeadSource = "referral" | "ads" | "internet" | "social" | "maps";
export type LeadRequestType = "course" | "masterclass" | "trial_lesson";

/** Thank-you page variants — identical to LeadRequestType, since the page a
 * lead is sent to after ApplyForm.vue submits is keyed by the same
 * "context" the lead itself records. */
export type ThankYouPageVariant = LeadRequestType;

/** POST body for /api/v1/leads (public route — see lead_handler.go). */
export interface ApplicationPayload {
  name: string;
  phone: string;
  email?: string;
  contactMethod: ContactMethod;
  source: LeadSource;
  requestType: LeadRequestType;
  relatedId?: number;
  /** The course/masterclass slug the visitor was looking at — lets the
   * admin panel show which specific one without cross-referencing
   * relatedId. Blank for trial_lesson (no specific entity). */
  relatedSlug?: string;
}

/** Mirrors backend/internal/model/model.go's LeadStatus — see app/lib/leadStatus.ts
 * for the single source of truth on ordering/labels. */
export type LeadStatus =
  | "new"
  | "in_progress"
  | "waiting_client"
  | "booked"
  | "postponed"
  | "closed_won"
  | "closed_lost";

export interface Lead extends ApplicationPayload {
  id: number;
  status: LeadStatus;
  createdAt: string;
  /** Points at the deduped Client profile this submission belongs to. */
  clientId: number;
  /** Set on leads auto-migrated from the old 3-value status enum, cleared
   * the moment someone explicitly picks a status for the lead. */
  needsStatusReview: boolean;
  /** Resolved course/masterclass title for relatedSlug — only present on
   * enriched list responses (GET /leads, client detail's requests). Empty
   * when relatedSlug is blank or the course/masterclass was since deleted;
   * fall back to relatedSlug in that case. */
  relatedName?: string;
}

/** The deduped customer profile a Lead attaches to — see ClientRepository
 * .FindByPhoneOrEmail on the backend for the phone/email matching rule. */
export interface Client {
  id: number;
  name: string;
  phone: string;
  email: string;
  createdAt: string;
  updatedAt: string;
}

/** Selects which of the two independent tag tables a Tag belongs to — the
 * two are never mixed (see migration 00032 on the backend). */
export type TagType = "product" | "client_type";

export interface Tag {
  id: number;
  name: string;
  /** "#rrggbb" background color — defaults to a neutral gray until
   * explicitly set (see AdminTagCombobox.vue). */
  color: string;
}

/** No author field by design — single-admin system, nothing to attribute a
 * comment to. */
export interface ClientComment {
  id: number;
  clientId: number;
  text: string;
  createdAt: string;
}

export interface Reminder {
  id: number;
  clientId: number;
  remindAt: string;
  note: string;
  completedAt?: string;
  createdAt: string;
}

/** What GET /api/v1/leads returns: a Lead enriched with everything its
 * list-page card needs in one response. */
export interface LeadListItem extends Lead {
  client: Client;
  productTags: Tag[];
  clientTypeTags: Tag[];
  latestCommentText?: string;
  latestCommentAt?: string;
  nextReminderAt?: string;
}

/** What GET /api/v1/clients/{id} returns — backs the client detail page. */
export interface ClientDetail extends Client {
  requests: Lead[];
  comments: ClientComment[];
  productTags: Tag[];
  clientTypeTags: Tag[];
  reminders: Reminder[];
}

/** What GET /api/v1/clients returns — backs the client-list page's cards. */
export interface ClientListItem extends Client {
  productTags: Tag[];
  clientTypeTags: Tag[];
  requestCount: number;
  latestStatus?: LeadStatus;
  latestRequestAt?: string;
  latestCommentAt?: string;
}

/** Normalized error shape used by useApi() so components never touch raw $fetch errors. */
export interface ApiError {
  message: string;
  status?: number;
}
