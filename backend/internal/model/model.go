package model

import "time"

type CourseStatus string

const (
	CourseStatusActive   CourseStatus = "active"
	CourseStatusArchived CourseStatus = "archived"
)

type Course struct {
	ID         int64        `db:"id" json:"id"`
	Slug       string       `db:"slug" json:"slug"`
	Title      string       `db:"title" json:"title"`
	ShortDesc  string       `db:"short_description" json:"shortDescription"`
	FullDesc   string       `db:"full_description" json:"fullDescription"`
	Status     CourseStatus `db:"status" json:"status"`
	CoverImage string       `db:"cover_image" json:"coverImage"`
	Gallery    []string     `db:"gallery" json:"gallery"`
	SortOrder  int          `db:"sort_order" json:"sortOrder"`
	CreatedAt  time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time    `db:"updated_at" json:"updatedAt"`
}

type CourseBlock struct {
	ID           int64  `db:"id" json:"id"`
	CourseID     int64  `db:"course_id" json:"courseId"`
	Title        string `db:"title" json:"title"`
	LessonsCount int    `db:"lessons_count" json:"lessonsCount"`
	Hours        int    `db:"hours" json:"hours"`
	Price        int    `db:"price" json:"price"`
	// nil — обычная цена, без скидки "было/стало".
	OldPrice  *int      `db:"old_price" json:"oldPrice,omitempty"`
	SortOrder int       `db:"sort_order" json:"sortOrder"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type Lesson struct {
	ID            int64     `db:"id" json:"id"`
	CourseBlockID int64     `db:"course_block_id" json:"courseBlockId"`
	Number        int       `db:"number" json:"number"`
	Title         string    `db:"title" json:"title"`
	Topics        string    `db:"topics" json:"topics"`
	Outcomes      string    `db:"outcomes" json:"outcomes"`
	DurationHours int       `db:"duration_hours" json:"durationHours"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}

// CourseModule and CourseDetail are the shape of the public "course page"
// aggregation (GET /api/v1/courses/{slug}/full) — Course/CourseBlock embed
// so their fields stay inline in the JSON object, matching the frontend's
// CourseModule/CourseDetail types (see frontend/app/types/api.ts). Assembled
// by CourseDetailService, not stored directly.
type CourseModule struct {
	CourseBlock
	Lessons []Lesson `json:"lessons"`
}

type CourseDetail struct {
	Course
	Modules []CourseModule `json:"modules"`
}

type MasterclassStatus string

const (
	MasterclassStatusActive   MasterclassStatus = "active"
	MasterclassStatusArchived MasterclassStatus = "archived"
)

type Masterclass struct {
	ID               int64             `db:"id" json:"id"`
	Slug             string            `db:"slug" json:"slug"`
	Title            string            `db:"title" json:"title"`
	ShortDesc        string            `db:"short_description" json:"shortDescription"`
	FullDesc         string            `db:"full_description" json:"fullDescription"`
	EndingText       string            `db:"ending_text" json:"endingText"`
	Duration         string            `db:"duration" json:"duration"`
	PriceGroup       int               `db:"price_group" json:"priceGroup"`
	PriceIndividual  int               `db:"price_individual" json:"priceIndividual"`
	PriceDescription string            `db:"price_description" json:"priceDescription"`
	CoverImage       string            `db:"cover_image" json:"coverImage"`
	Status           MasterclassStatus `db:"status" json:"status"`
	CreatedAt        time.Time         `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time         `db:"updated_at" json:"updatedAt"`
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
	LeadStatusNew        LeadStatus = "new"
	LeadStatusInProgress LeadStatus = "in_progress"
	LeadStatusClosed     LeadStatus = "closed"
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
	Status        LeadStatus      `db:"status" json:"status"`
	CreatedAt     time.Time       `db:"created_at" json:"createdAt"`
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

type Feature struct {
	ID          int64     `db:"id" json:"id"`
	Page        string    `db:"page" json:"page"`
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

type AdminUser struct {
	ID           int64     `db:"id" json:"id"`
	Login        string    `db:"login" json:"login"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}
