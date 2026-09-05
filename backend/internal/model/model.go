package model

import "time"

// Visible lets an admin hide a section from the public site (ListSections)
// without deleting it — its courses stay editable in the admin panel.
type CourseSection struct {
	ID          int64     `db:"id" json:"id"`
	Heading     string    `db:"heading" json:"heading"`
	Description string    `db:"description" json:"description"`
	Visible     bool      `db:"visible" json:"visible"`
	SortOrder   int       `db:"sort_order" json:"sortOrder"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// Visible lets an admin hide a course from the public site — both the
// homepage listing and its own /courses/{slug} page (see CourseCatalogService).
//
// CoverImage/LessonCount/TimeLength/Price/DisplayStyle are the course's OWN
// card — splitting into CourseBlocks is optional (migration 00020), not
// mandatory the way it used to be. A course with zero CourseBlock rows
// renders as exactly one homepage card built from these fields
// (CourseCatalogService synthesizes a CourseBlock-shaped entry from them);
// a course with one or more CourseBlock rows renders one card per block
// instead, and these course-level fields go unused. Adding a second track to
// an already-simple course is then an explicit, visible act (creating a
// CourseBlock) instead of something every course was forced into from the
// start.
type Course struct {
	ID           int64                   `db:"id" json:"id"`
	SectionID    int64                   `db:"section_id" json:"sectionId"`
	Slug         string                  `db:"slug" json:"slug"`
	Name         string                  `db:"name" json:"name"`
	Description  string                  `db:"description" json:"description"`
	CoverImage   string                  `db:"cover_image" json:"coverImage"`
	LessonCount  string                  `db:"lesson_count" json:"lessonCount"`
	TimeLength   string                  `db:"time_length" json:"timeLength"`
	Price        string                  `db:"price" json:"price"`
	DisplayStyle CourseBlockDisplayStyle `db:"display_style" json:"displayStyle"`
	Visible      bool                    `db:"visible" json:"visible"`
	SortOrder    int                     `db:"sort_order" json:"sortOrder"`
	// SingleCard collapses a multi-block course into one homepage card built
	// from this course's own fields instead of one card per block — see
	// syntheticBlock() in course_catalog_service.go, which this flag reuses.
	SingleCard bool `db:"single_card" json:"singleCard"`
	// FAQTitle/FAQDescription are the heading and intro text shown above this
	// course's FAQ block (CourseFAQItem below) on its course page, right
	// after the apply form. FAQVisible toggles the whole block — title,
	// description and items — without deleting anything, same convention as
	// Visible above. Not to be confused with the global faq_items table
	// (model.FAQItem) rendered once on the homepage.
	FAQTitle       string    `db:"faq_title" json:"faqTitle"`
	FAQDescription string    `db:"faq_description" json:"faqDescription"`
	FAQVisible     bool      `db:"faq_visible" json:"faqVisible"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}

type CourseBlockDisplayStyle string

const (
	DisplayStyleBlueBeige  CourseBlockDisplayStyle = "blue-beige"
	DisplayStyleBrownBeige CourseBlockDisplayStyle = "brown-beige"
	DisplayStyleBeigeBlue  CourseBlockDisplayStyle = "beige-blue"
	DisplayStyleBeigeBrown CourseBlockDisplayStyle = "beige-brown"
)

// CourseBlock carries its own lesson list directly (via CourseBlockWithLessons
// below) — there's no separate "curriculum" level. lesson_count/time_length/
// price are free text ("7 занятий", "30 часов", "38 500 ₽" or even a
// "было/стало" discount narrative) rather than numbers, so the admin can
// phrase them however the site copy needs without any formatting logic.
// DisplayStyle picks the background/text color pair for this block's
// homepage card — chosen per block, not per course, since a multi-block
// course renders one card per block (see frontend CourseCard.vue). Visible
// lets an admin hide just this one block (its card on the homepage, and its
// badge/curriculum on the course page) without touching the rest of the course.
type CourseBlock struct {
	ID           int64                   `db:"id" json:"id"`
	CourseID     int64                   `db:"course_id" json:"courseId"`
	BlockName    string                  `db:"block_name" json:"blockName"`
	Description  string                  `db:"description" json:"description"`
	BlockCover   string                  `db:"block_cover" json:"blockCover"`
	LessonCount  string                  `db:"lesson_count" json:"lessonCount"`
	TimeLength   string                  `db:"time_length" json:"timeLength"`
	Price        string                  `db:"price" json:"price"`
	DisplayStyle CourseBlockDisplayStyle `db:"display_style" json:"displayStyle"`
	Visible      bool                    `db:"visible" json:"visible"`
	SortOrder    int                     `db:"sort_order" json:"sortOrder"`
	CreatedAt    time.Time               `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time               `db:"updated_at" json:"updatedAt"`
}

// Lesson belongs to exactly one parent — CourseBlockID (a course split into
// blocks) or CourseID (a course with no blocks, editing its lessons
// directly) — never both, never neither (enforced by the
// lessons_exactly_one_parent CHECK constraint, migration 00023).
type Lesson struct {
	ID            int64     `db:"id" json:"id"`
	CourseBlockID *int64    `db:"course_block_id" json:"courseBlockId,omitempty"`
	CourseID      *int64    `db:"course_id" json:"courseId,omitempty"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	SortOrder     int       `db:"sort_order" json:"sortOrder"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}

// The types below are response-only aggregates assembled by
// CourseCatalogService (not stored directly) — BlockCount/Length are
// computed from len(...) rather than persisted columns, so they can never
// drift out of sync with the actual child rows.

type CourseBlockWithLessons struct {
	CourseBlock
	Lessons []Lesson `json:"lessons"`
}

// CourseWithBlocks is the public course-page shape (GET /api/v1/courses/{slug}/full).
type CourseWithBlocks struct {
	Course
	BlockCount int                      `json:"blockCount"`
	Blocks     []CourseBlockWithLessons `json:"blocks"`
	// FAQItems is this course's own Q&A list (see Course.FAQVisible/FAQTitle/
	// FAQDescription) — always populated regardless of FAQVisible, so an
	// admin previewing a hidden block doesn't need a separate endpoint; the
	// frontend is the one that gates rendering on FAQVisible.
	FAQItems []CourseFAQItem `json:"faqItems"`
}

// CourseFAQItem is one Q&A pair in a single course's FAQ block — scoped by
// CourseID, distinct from the global, homepage-only FAQItem above.
type CourseFAQItem struct {
	ID        int64     `db:"id" json:"id"`
	CourseID  int64     `db:"course_id" json:"courseId"`
	Question  string    `db:"question" json:"question"`
	Answer    string    `db:"answer" json:"answer"`
	SortOrder int       `db:"sort_order" json:"sortOrder"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// CourseSummary is the homepage listing shape — blocks without lesson text,
// since the homepage only needs each course's first block's cover/lessonCount/
// timeLength (and blockName when there's more than one block).
type CourseSummary struct {
	Course
	BlockCount int           `json:"blockCount"`
	Blocks     []CourseBlock `json:"blocks"`
}

// CourseSectionWithCourses is the public homepage shape (GET /api/v1/course-sections/full).
type CourseSectionWithCourses struct {
	CourseSection
	Length  int             `json:"length"`
	Courses []CourseSummary `json:"courses"`
}

type MasterclassStatus string

const (
	MasterclassStatusActive   MasterclassStatus = "active"
	MasterclassStatusArchived MasterclassStatus = "archived"
)

type Masterclass struct {
	ID          int64  `db:"id" json:"id"`
	Slug        string `db:"slug" json:"slug"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	// Description2 is an optional second paragraph (e.g. a closing
	// call-to-action) — not every masterclass needs one.
	Description2 string `db:"description2" json:"description2"`
	EndingText   string `db:"ending_text" json:"endingText"`
	Duration     string `db:"duration" json:"duration"`
	// Price is free text (like CourseBlock.Price) — not every masterclass
	// has exactly one price ("3000₽", "3000₽ или 4500₽", "от 2500₽"), so a
	// fixed group/individual pair of numbers doesn't fit every case.
	Price      string            `db:"price" json:"price"`
	CoverImage string            `db:"cover_image" json:"coverImage"`
	Status     MasterclassStatus `db:"status" json:"status"`
	CreatedAt  time.Time         `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time         `db:"updated_at" json:"updatedAt"`
}

type Teacher struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Photo       string    `db:"photo" json:"photo"`
	Description string    `db:"description" json:"description"`
	SortOrder   int       `db:"sort_order" json:"sortOrder"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// GalleryPhoto is one slide of the homepage's vertical-photo carousel
// (between "Преимущества" and "О школе").
type GalleryPhoto struct {
	ID        int64     `db:"id" json:"id"`
	Image     string    `db:"image" json:"image"`
	SortOrder int       `db:"sort_order" json:"sortOrder"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type BlogPostStatus string

const (
	BlogPostStatusDraft     BlogPostStatus = "draft"
	BlogPostStatusPublished BlogPostStatus = "published"
)

type BlogPost struct {
	ID          int64          `db:"id" json:"id"`
	Slug        string         `db:"slug" json:"slug"`
	Title       string         `db:"title" json:"title"`
	CoverImage  string         `db:"cover_image" json:"coverImage"`
	Category    string         `db:"category" json:"category"`
	Tags        []string       `db:"tags" json:"tags"`
	Author      string         `db:"author" json:"author"`
	PublishedAt *time.Time     `db:"published_at" json:"publishedAt,omitempty"`
	Content     string         `db:"content" json:"content"`
	Status      BlogPostStatus `db:"status" json:"status"`
	CreatedAt   time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updatedAt"`
}

type ContactMethod string

const (
	ContactMethodCall     ContactMethod = "call"
	ContactMethodTelegram ContactMethod = "telegram"
	ContactMethodWhatsapp ContactMethod = "whatsapp"
	ContactMethodMax      ContactMethod = "max"
)

type LeadSource string

const (
	LeadSourceReferral LeadSource = "referral"
	LeadSourceAds      LeadSource = "ads"
	LeadSourceInternet LeadSource = "internet"
	LeadSourceSocial   LeadSource = "social"
	LeadSourceMaps     LeadSource = "maps"
)

type LeadRequestType string

const (
	LeadRequestTypeCourse      LeadRequestType = "course"
	LeadRequestTypeMasterclass LeadRequestType = "masterclass"
	LeadRequestTypeTrialLesson LeadRequestType = "trial_lesson"
)

type LeadStatus string

const (
	LeadStatusNew           LeadStatus = "new"
	LeadStatusInProgress    LeadStatus = "in_progress"
	LeadStatusWaitingClient LeadStatus = "waiting_client"
	LeadStatusBooked        LeadStatus = "booked"
	LeadStatusPostponed     LeadStatus = "postponed"
	LeadStatusClosedWon     LeadStatus = "closed_won"
	LeadStatusClosedLost    LeadStatus = "closed_lost"
)

type Lead struct {
	ID            int64           `db:"id" json:"id"`
	Name          string          `db:"name" json:"name"`
	Phone         string          `db:"phone" json:"phone"`
	Email         string          `db:"email" json:"email"`
	ContactMethod ContactMethod   `db:"contact_method" json:"contactMethod"`
	Source        LeadSource      `db:"source" json:"source"`
	RequestType   LeadRequestType `db:"request_type" json:"requestType"`
	RelatedID     *int64          `db:"related_id" json:"relatedId,omitempty"`
	// The course/masterclass slug the visitor was looking at when they
	// submitted — sent straight from the frontend page, not looked up via
	// RelatedID, so the admin panel can show it without a join. Blank for
	// trial_lesson leads (no specific entity to name).
	RelatedSlug string `db:"related_slug" json:"relatedSlug"`
	// RelatedName is the resolved course/masterclass title for RelatedSlug
	// — not a real column (db:"-"), populated only by the enriched list
	// queries (LeadRepository.ListWithClient/ListByClientID) via a join,
	// since that's the only place a human-readable name is actually
	// needed. Empty when RelatedSlug is blank or the course/masterclass
	// was since deleted/renamed — callers should fall back to RelatedSlug.
	RelatedName string     `db:"-" json:"relatedName,omitempty"`
	Status      LeadStatus `db:"status" json:"status"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	// ClientID points at the deduped client profile this submission belongs
	// to (see ClientRepository.FindByPhoneOrEmail) — Name/Phone/Email above
	// stay an untouched historical snapshot of what was submitted this one
	// time, distinct from the client's current-best-known contact info.
	ClientID int64 `db:"client_id" json:"clientId"`
	// NeedsStatusReview is set on leads auto-migrated from the old 3-value
	// status enum (old "closed" -> closed_won, see migration 00031) and
	// cleared the moment someone explicitly picks a status for the lead.
	NeedsStatusReview bool `db:"needs_status_review" json:"needsStatusReview"`
}

// Client is the deduped customer profile a Lead attaches to — see
// LeadService.Create for the phone/email matching that decides whether a
// new submission joins an existing Client or creates a new one.
type Client struct {
	ID    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Phone string `db:"phone" json:"phone"`
	// PhoneNormalized is digits-only (RU "8" dialing prefix folded to "7")
	// — the matching key used by FindByPhoneOrEmail. Not exposed over JSON;
	// Phone above is what the UI shows.
	PhoneNormalized string    `db:"phone_normalized" json:"-"`
	Email           string    `db:"email" json:"email"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}

// TagType selects which of the two independent tag tables a Tag belongs
// to. The two are never mixed — see migration 00032.
type TagType string

const (
	TagTypeProduct    TagType = "product"
	TagTypeClientType TagType = "client_type"
)

type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	// Color is a "#rrggbb" hex string — defaults to a neutral gray
	// (migration 00035) until explicitly set via TagService.SetColor.
	Color string `db:"color" json:"color"`
}

// ClientComment has no author field by design — single-admin system, no
// one to attribute a comment to.
type ClientComment struct {
	ID        int64     `db:"id" json:"id"`
	ClientID  int64     `db:"client_id" json:"clientId"`
	Text      string    `db:"text" json:"text"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// Reminder is a "contact this client again by RemindAt" flag. RemindAt is
// date-granularity (not a timestamp) — "N days from now" and a "due today"
// badge don't need time-of-day precision, and this sidesteps timezone
// complexity entirely.
type Reminder struct {
	ID          int64      `db:"id" json:"id"`
	ClientID    int64      `db:"client_id" json:"clientId"`
	RemindAt    time.Time  `db:"remind_at" json:"remindAt"`
	Note        string     `db:"note" json:"note"`
	CompletedAt *time.Time `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
}

// LeadListItem is what GET /leads returns: a Lead enriched with everything
// its list-page card needs in one response, so the frontend never has to
// N+1 across /clients, /tags, /comments per row.
type LeadListItem struct {
	Lead
	Client            Client     `json:"client"`
	ProductTags       []Tag      `json:"productTags"`
	ClientTypeTags    []Tag      `json:"clientTypeTags"`
	LatestCommentText string     `json:"latestCommentText,omitempty"`
	LatestCommentAt   *time.Time `json:"latestCommentAt,omitempty"`
	NextReminderAt    *time.Time `json:"nextReminderAt,omitempty"`
}

// ClientDetail backs the client detail page in one round trip.
type ClientDetail struct {
	Client
	Requests       []Lead          `json:"requests"`
	Comments       []ClientComment `json:"comments"`
	ProductTags    []Tag           `json:"productTags"`
	ClientTypeTags []Tag           `json:"clientTypeTags"`
	Reminders      []Reminder      `json:"reminders"`
}

// ClientListItem is what GET /clients returns: a Client enriched with
// everything the client-list-page card needs (both tag lists, and a
// summary of their most recent activity) in one response.
type ClientListItem struct {
	Client
	ProductTags     []Tag      `json:"productTags"`
	ClientTypeTags  []Tag      `json:"clientTypeTags"`
	RequestCount    int        `json:"requestCount"`
	LatestStatus    LeadStatus `json:"latestStatus,omitempty"`
	LatestRequestAt *time.Time `json:"latestRequestAt,omitempty"`
	LatestCommentAt *time.Time `json:"latestCommentAt,omitempty"`
}

type FAQItem struct {
	ID        int64     `db:"id" json:"id"`
	Question  string    `db:"question" json:"question"`
	Answer    string    `db:"answer" json:"answer"`
	SortOrder int       `db:"sort_order" json:"sortOrder"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type PageContent struct {
	Key       string    `db:"key" json:"key"`
	Label     string    `db:"label" json:"label"`
	Value     string    `db:"value" json:"value"`
	Type      string    `db:"type" json:"type"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Icon is an admin-uploaded SVG, usable anywhere a "built-in" icon key
// (FEATURE_ICONS on the frontend) is accepted — see the icons table's
// migration comment for the "icon:<id>" value convention that references
// one of these. `SVG` is already sanitized (see IconService.Create) by the
// time it reaches this struct — every reader may safely render it with
// v-html.
type Icon struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	SVG       string    `db:"svg" json:"svg"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type Feature struct {
	ID   int64  `db:"id" json:"id"`
	Page string `db:"page" json:"page"`
	// Either a bare FEATURE_ICONS key (e.g. "gift") or "icon:<id>"
	// referencing an uploaded Icon — see that type's doc comment.
	Icon        string    `db:"icon" json:"icon"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	SortOrder   int       `db:"sort_order" json:"sortOrder"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type AboutItem struct {
	ID          int64     `db:"id" json:"id"`
	Badge       string    `db:"badge" json:"badge"`
	Description string    `db:"description" json:"description"`
	SortOrder   int       `db:"sort_order" json:"sortOrder"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type SocialLink struct {
	ID    int64  `db:"id" json:"id"`
	Label string `db:"label" json:"label"`
	Href  string `db:"href" json:"href"`
	// Legally required disclaimer text for some platforms in Russia (e.g.
	// Meta-owned apps) — empty when none applies.
	Disclaimer string    `db:"disclaimer" json:"disclaimer"`
	SortOrder  int       `db:"sort_order" json:"sortOrder"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
}

// ExportFile is one object from Garage/S3 storage, bundled whole into a
// content export so "download the export" is a real, self-contained backup
// instead of a DB dump full of /uploads/<key> paths that resolve to nothing
// once restored somewhere else. Data round-trips through JSON as base64
// ([]byte's default encoding) — simple over efficient, since this travels as
// one JSON document already (see ContentExportService).
type ExportFile struct {
	Key         string `json:"key"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

// SiteContent is a full snapshot of every admin-managed content entity —
// everything except accounts (AdminUser) and customer inquiries (Lead),
// which aren't "site content" — plus every file in object storage (Files).
// Used by ContentExportService for backup/restore and for moving content
// between environments.
type SiteContent struct {
	Version        int             `json:"version"`
	ExportedAt     time.Time       `json:"exportedAt"`
	CourseSections []CourseSection `json:"courseSections"`
	Courses        []Course        `json:"courses"`
	CourseBlocks   []CourseBlock   `json:"courseBlocks"`
	Lessons        []Lesson        `json:"lessons"`
	CourseFAQItems []CourseFAQItem `json:"courseFaqItems"`
	Masterclasses  []Masterclass   `json:"masterclasses"`
	Teachers       []Teacher       `json:"teachers"`
	GalleryPhotos  []GalleryPhoto  `json:"galleryPhotos"`
	BlogPosts      []BlogPost      `json:"blogPosts"`
	FAQItems       []FAQItem       `json:"faqItems"`
	Features       []Feature       `json:"features"`
	AboutItems     []AboutItem     `json:"aboutItems"`
	SocialLinks    []SocialLink    `json:"socialLinks"`
	PageContent    []PageContent   `json:"pageContent"`
	Files          []ExportFile    `json:"files"`
}

type AdminUser struct {
	ID           int64     `db:"id" json:"id"`
	Login        string    `db:"login" json:"login"`
	PasswordHash string    `db:"password_hash" json:"-"`
	TokenVersion int       `db:"token_version" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}
