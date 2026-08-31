package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

// ContentExportVersion is bumped whenever the SiteContent shape changes in a
// way an importer needs to know about (a field renamed/removed, not just
// added) — Import rejects anything newer than what this build understands.
const ContentExportVersion = 1

type ImportMode string

const (
	// ImportModeMerge upserts by id (or, for PageContent, by key): existing
	// rows are updated, new ids are inserted, and anything NOT present in the
	// import is left alone. Never deletes.
	ImportModeMerge ImportMode = "merge"
	// ImportModeReplace deletes every row of every entity first (except
	// PageContent, which has no delete concept — see importPageContent) and
	// re-inserts exactly what's in the import. The imported file becomes the
	// single source of truth; anything not in it is gone.
	ImportModeReplace ImportMode = "replace"
)

type ImportResult struct {
	Mode   ImportMode     `json:"mode"`
	Counts map[string]int `json:"counts"`
	// PageContentSkipped lists keys present in the import that don't exist
	// on this server (page_content rows are seeded by migrations, not
	// created through this feature — see importPageContent).
	PageContentSkipped []string `json:"pageContentSkipped,omitempty"`
}

// ContentExportService reads and writes the entire site's editable content —
// every admin-managed entity except accounts (AdminUser) and customer
// inquiries (Lead) — as one JSON document, for backup/restore and for moving
// content between environments (e.g. staging -> prod).
//
// It talks to the pool directly instead of going through the per-entity
// repository/service layers (the same documented exception Services.DB makes
// for /readyz): this is a bulk, transactional, cross-entity operation, and
// none of the existing repositories support transactions or inserting a
// caller-supplied id. Retrofitting all thirteen of them with upsert/tx-aware
// variants just for this one feature would be a far bigger, riskier change
// than keeping it self-contained here.
type ContentExportService struct {
	db *pgxpool.Pool
}

func NewContentExportService(db *pgxpool.Pool) *ContentExportService {
	return &ContentExportService{db: db}
}

// --- Export -----------------------------------------------------------

func (s *ContentExportService) Export(ctx context.Context) (model.SiteContent, error) {
	var out model.SiteContent
	out.Version = ContentExportVersion
	out.ExportedAt = time.Now().UTC()

	var err error
	if out.CourseSections, err = queryAll(ctx, s.db, `SELECT id, heading, description, visible, sort_order, created_at, updated_at FROM course_sections ORDER BY id`, scanCourseSection); err != nil {
		return out, fmt.Errorf("export course sections: %w", err)
	}
	if out.Courses, err = queryAll(ctx, s.db, `SELECT id, section_id, slug, name, description, cover_image, lesson_count, time_length, price, display_style, visible, sort_order, created_at, updated_at FROM courses ORDER BY id`, scanCourse); err != nil {
		return out, fmt.Errorf("export courses: %w", err)
	}
	if out.CourseBlocks, err = queryAll(ctx, s.db, `SELECT id, course_id, block_name, description, block_cover, lesson_count, time_length, price, display_style, visible, sort_order, created_at, updated_at FROM course_blocks ORDER BY id`, scanCourseBlock); err != nil {
		return out, fmt.Errorf("export course blocks: %w", err)
	}
	if out.Lessons, err = queryAll(ctx, s.db, `SELECT id, course_block_id, course_id, name, description, sort_order, created_at, updated_at FROM lessons ORDER BY id`, scanLesson); err != nil {
		return out, fmt.Errorf("export lessons: %w", err)
	}
	if out.Masterclasses, err = queryAll(ctx, s.db, `SELECT `+exportMasterclassColumns+` FROM masterclasses ORDER BY id`, scanExportMasterclass); err != nil {
		return out, fmt.Errorf("export masterclasses: %w", err)
	}
	if out.Teachers, err = queryAll(ctx, s.db, `SELECT id, name, photo, description, sort_order, created_at, updated_at FROM teachers ORDER BY id`, scanTeacher); err != nil {
		return out, fmt.Errorf("export teachers: %w", err)
	}
	if out.GalleryPhotos, err = queryAll(ctx, s.db, `SELECT id, image, sort_order, created_at, updated_at FROM gallery_photos ORDER BY id`, scanGalleryPhoto); err != nil {
		return out, fmt.Errorf("export gallery photos: %w", err)
	}
	if out.BlogPosts, err = queryAll(ctx, s.db, `SELECT `+exportBlogPostColumns+` FROM blog_posts ORDER BY id`, scanExportBlogPost); err != nil {
		return out, fmt.Errorf("export blog posts: %w", err)
	}
	if out.FAQItems, err = queryAll(ctx, s.db, `SELECT id, question, answer, sort_order, created_at, updated_at FROM faq_items ORDER BY id`, scanFAQItem); err != nil {
		return out, fmt.Errorf("export faq items: %w", err)
	}
	if out.Features, err = queryAll(ctx, s.db, `SELECT id, page, icon, title, description, sort_order, created_at, updated_at FROM features ORDER BY id`, scanFeature); err != nil {
		return out, fmt.Errorf("export features: %w", err)
	}
	if out.AboutItems, err = queryAll(ctx, s.db, `SELECT id, badge, description, sort_order, created_at, updated_at FROM about_items ORDER BY id`, scanAboutItem); err != nil {
		return out, fmt.Errorf("export about items: %w", err)
	}
	if out.SocialLinks, err = queryAll(ctx, s.db, `SELECT id, label, href, disclaimer, sort_order, created_at, updated_at FROM social_links ORDER BY id`, scanSocialLink); err != nil {
		return out, fmt.Errorf("export social links: %w", err)
	}
	if out.PageContent, err = queryAll(ctx, s.db, `SELECT key, label, value, type, updated_at FROM page_content ORDER BY key`, scanPageContent); err != nil {
		return out, fmt.Errorf("export page content: %w", err)
	}
	return out, nil
}

// --- Import -------------------------------------------------------------

// Import applies data in one all-or-nothing transaction: any failure (a bad
// foreign key, a CHECK-constraint violation on an enum column, a duplicate
// slug) rolls back the entire import rather than leaving the site half
// written. That's deliberately simpler — and safer for a destructive
// operation like replace — than per-row partial-success reporting.
func (s *ContentExportService) Import(ctx context.Context, data model.SiteContent, mode ImportMode) (ImportResult, error) {
	if mode != ImportModeMerge && mode != ImportModeReplace {
		return ImportResult{}, errors.Join(ErrValidation, fmt.Errorf("unknown import mode %q", mode))
	}
	if data.Version > ContentExportVersion {
		return ImportResult{}, errors.Join(ErrValidation, fmt.Errorf("export version %d is newer than this server understands (%d)", data.Version, ContentExportVersion))
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return ImportResult{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op once Commit succeeds

	result := ImportResult{Mode: mode, Counts: map[string]int{}}

	if mode == ImportModeReplace {
		// course_sections cascades to courses/course_blocks/lessons
		// (ON DELETE CASCADE, migration 00016) — deleting it alone clears the
		// whole tree. page_content is never deleted (see importPageContent).
		for _, table := range []string{"course_sections", "masterclasses", "teachers", "gallery_photos", "blog_posts", "faq_items", "features", "about_items", "social_links"} {
			if _, err := tx.Exec(ctx, "DELETE FROM "+table); err != nil {
				return ImportResult{}, fmt.Errorf("clear %s: %w", table, err)
			}
		}
	}

	conflictCol := "" // plain INSERT after a full delete
	if mode == ImportModeMerge {
		conflictCol = "id"
	}

	n, err := bulkWrite(ctx, tx, "course_sections", []string{"id", "heading", "description", "visible", "sort_order"}, data.CourseSections, courseSectionArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import course sections: %w", err)
	}
	result.Counts["courseSections"] = n

	n, err = bulkWrite(ctx, tx, "courses", []string{"id", "section_id", "slug", "name", "description", "cover_image", "lesson_count", "time_length", "price", "display_style", "visible", "sort_order"}, data.Courses, courseArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import courses: %w", err)
	}
	result.Counts["courses"] = n

	n, err = bulkWrite(ctx, tx, "course_blocks", []string{"id", "course_id", "block_name", "description", "block_cover", "lesson_count", "time_length", "price", "display_style", "visible", "sort_order"}, data.CourseBlocks, courseBlockArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import course blocks: %w", err)
	}
	result.Counts["courseBlocks"] = n

	n, err = bulkWrite(ctx, tx, "lessons", []string{"id", "course_block_id", "course_id", "name", "description", "sort_order"}, data.Lessons, lessonArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import lessons: %w", err)
	}
	result.Counts["lessons"] = n

	n, err = bulkWrite(ctx, tx, "masterclasses", []string{"id", "slug", "title", "description", "description2", "ending_text", "duration", "price", "cover_image", "status"}, data.Masterclasses, masterclassArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import masterclasses: %w", err)
	}
	result.Counts["masterclasses"] = n

	n, err = bulkWrite(ctx, tx, "teachers", []string{"id", "name", "photo", "description", "sort_order"}, data.Teachers, teacherArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import teachers: %w", err)
	}
	result.Counts["teachers"] = n

	n, err = bulkWrite(ctx, tx, "gallery_photos", []string{"id", "image", "sort_order"}, data.GalleryPhotos, galleryPhotoArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import gallery photos: %w", err)
	}
	result.Counts["galleryPhotos"] = n

	n, err = bulkWrite(ctx, tx, "blog_posts", []string{"id", "slug", "title", "cover_image", "category", "tags", "author", "published_at", "content", "status"}, data.BlogPosts, blogPostArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import blog posts: %w", err)
	}
	result.Counts["blogPosts"] = n

	n, err = bulkWrite(ctx, tx, "faq_items", []string{"id", "question", "answer", "sort_order"}, data.FAQItems, faqItemArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import faq items: %w", err)
	}
	result.Counts["faqItems"] = n

	n, err = bulkWrite(ctx, tx, "features", []string{"id", "page", "icon", "title", "description", "sort_order"}, data.Features, featureArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import features: %w", err)
	}
	result.Counts["features"] = n

	n, err = bulkWrite(ctx, tx, "about_items", []string{"id", "badge", "description", "sort_order"}, data.AboutItems, aboutItemArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import about items: %w", err)
	}
	result.Counts["aboutItems"] = n

	n, err = bulkWrite(ctx, tx, "social_links", []string{"id", "label", "href", "disclaimer", "sort_order"}, data.SocialLinks, socialLinkArgs, conflictCol)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import social links: %w", err)
	}
	result.Counts["socialLinks"] = n

	// Every id-bearing table above just got explicit ids inserted — bump each
	// sequence past the highest one, or the next plain admin-panel Create()
	// (which never specifies an id) will collide with an imported row.
	for _, table := range []string{"course_sections", "courses", "course_blocks", "lessons", "masterclasses", "teachers", "gallery_photos", "blog_posts", "faq_items", "features", "about_items", "social_links"} {
		if _, err := tx.Exec(ctx, `SELECT setval(pg_get_serial_sequence($1, 'id'), GREATEST((SELECT COALESCE(MAX(id), 0) FROM `+table+`), 1))`, table); err != nil {
			return ImportResult{}, fmt.Errorf("reset %s id sequence: %w", table, err)
		}
	}

	updated, skipped, err := importPageContent(ctx, tx, data.PageContent)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import page content: %w", err)
	}
	result.Counts["pageContent"] = updated
	result.PageContentSkipped = skipped

	if err := tx.Commit(ctx); err != nil {
		return ImportResult{}, err
	}
	return result, nil
}

// importPageContent is its own function, not a bulkWrite call: page_content
// rows are seeded by migrations (see page_content_repository.go), not
// created through the API, so there's no INSERT/DELETE path for it in either
// import mode — only UPDATE for keys that already exist on this server. A
// key present in the import but not on this server is reported as skipped
// rather than erroring the whole import, since two environments can
// legitimately be a migration or two apart.
func importPageContent(ctx context.Context, tx pgx.Tx, items []model.PageContent) (updated int, skipped []string, err error) {
	for _, item := range items {
		tag, err := tx.Exec(ctx, `UPDATE page_content SET value = $1, updated_at = now() WHERE key = $2`, item.Value, item.Key)
		if err != nil {
			return 0, nil, err
		}
		if tag.RowsAffected() == 0 {
			skipped = append(skipped, item.Key)
			continue
		}
		updated++
	}
	return updated, skipped, nil
}

// --- generic read/write helpers ------------------------------------------

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func queryAll[T any](ctx context.Context, db querier, query string, scan func(pgx.CollectableRow) (T, error)) ([]T, error) {
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, scan)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []T{}
	}
	return items, nil
}

// bulkWrite inserts every row of `data` into `table` in one multi-row
// INSERT. When conflictCol is empty it's a plain insert (used in replace
// mode, right after the table was cleared); when set, it's
// "ON CONFLICT (conflictCol) DO UPDATE" (merge mode's upsert-by-id). cols[0]
// must be the id/conflict column. Returns the row count written.
func bulkWrite[T any](ctx context.Context, tx pgx.Tx, table string, cols []string, data []T, argsFn func(T) []any, conflictCol string) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var sb strings.Builder
	args := make([]any, 0, len(data)*len(cols))
	fmt.Fprintf(&sb, "INSERT INTO %s (%s) VALUES ", table, strings.Join(cols, ", "))
	for i, row := range data {
		if i > 0 {
			sb.WriteString(", ")
		}
		rowArgs := argsFn(row)
		placeholders := make([]string, len(rowArgs))
		for j := range rowArgs {
			args = append(args, rowArgs[j])
			placeholders[j] = fmt.Sprintf("$%d", len(args))
		}
		fmt.Fprintf(&sb, "(%s)", strings.Join(placeholders, ", "))
	}
	if conflictCol != "" {
		updateSet := make([]string, 0, len(cols)-1)
		for _, c := range cols[1:] {
			updateSet = append(updateSet, fmt.Sprintf("%s = EXCLUDED.%s", c, c))
		}
		fmt.Fprintf(&sb, " ON CONFLICT (%s) DO UPDATE SET %s", conflictCol, strings.Join(updateSet, ", "))
	}

	if _, err := tx.Exec(ctx, sb.String(), args...); err != nil {
		return 0, err
	}
	return len(data), nil
}

// --- per-entity scan/arg functions ----------------------------------------
// Mirrors each repository's own Scan order exactly (see e.g.
// course_section_repository.go) so a column can't silently drift between
// this file and the hand-written repositories.

func scanCourseSection(row pgx.CollectableRow) (model.CourseSection, error) {
	var m model.CourseSection
	err := row.Scan(&m.ID, &m.Heading, &m.Description, &m.Visible, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func courseSectionArgs(m model.CourseSection) []any {
	return []any{m.ID, m.Heading, m.Description, m.Visible, m.SortOrder}
}

func scanCourse(row pgx.CollectableRow) (model.Course, error) {
	var m model.Course
	err := row.Scan(&m.ID, &m.SectionID, &m.Slug, &m.Name, &m.Description, &m.CoverImage, &m.LessonCount, &m.TimeLength, &m.Price, &m.DisplayStyle, &m.Visible, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func courseArgs(m model.Course) []any {
	return []any{m.ID, m.SectionID, m.Slug, m.Name, m.Description, m.CoverImage, m.LessonCount, m.TimeLength, m.Price, m.DisplayStyle, m.Visible, m.SortOrder}
}

func scanCourseBlock(row pgx.CollectableRow) (model.CourseBlock, error) {
	var m model.CourseBlock
	err := row.Scan(&m.ID, &m.CourseID, &m.BlockName, &m.Description, &m.BlockCover, &m.LessonCount, &m.TimeLength, &m.Price, &m.DisplayStyle, &m.Visible, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func courseBlockArgs(m model.CourseBlock) []any {
	return []any{m.ID, m.CourseID, m.BlockName, m.Description, m.BlockCover, m.LessonCount, m.TimeLength, m.Price, m.DisplayStyle, m.Visible, m.SortOrder}
}

func scanLesson(row pgx.CollectableRow) (model.Lesson, error) {
	var m model.Lesson
	err := row.Scan(&m.ID, &m.CourseBlockID, &m.CourseID, &m.Name, &m.Description, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func lessonArgs(m model.Lesson) []any {
	return []any{m.ID, m.CourseBlockID, m.CourseID, m.Name, m.Description, m.SortOrder}
}

// exportMasterclassColumns/scanExportMasterclass duplicate
// repository.masterclassColumns/scanMasterclass (unexported there, so not
// reachable from this package) — keep both in sync if either changes.
const exportMasterclassColumns = "id, slug, title, description, description2, ending_text, duration, price, cover_image, status, created_at, updated_at"

func scanExportMasterclass(row pgx.CollectableRow) (model.Masterclass, error) {
	var m model.Masterclass
	err := row.Scan(&m.ID, &m.Slug, &m.Title, &m.Description, &m.Description2, &m.EndingText, &m.Duration, &m.Price, &m.CoverImage, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func masterclassArgs(m model.Masterclass) []any {
	return []any{m.ID, m.Slug, m.Title, m.Description, m.Description2, m.EndingText, m.Duration, m.Price, m.CoverImage, m.Status}
}

func scanTeacher(row pgx.CollectableRow) (model.Teacher, error) {
	var m model.Teacher
	err := row.Scan(&m.ID, &m.Name, &m.Photo, &m.Description, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func teacherArgs(m model.Teacher) []any {
	return []any{m.ID, m.Name, m.Photo, m.Description, m.SortOrder}
}

func scanGalleryPhoto(row pgx.CollectableRow) (model.GalleryPhoto, error) {
	var m model.GalleryPhoto
	err := row.Scan(&m.ID, &m.Image, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func galleryPhotoArgs(m model.GalleryPhoto) []any {
	return []any{m.ID, m.Image, m.SortOrder}
}

// exportBlogPostColumns/scanExportBlogPost duplicate
// repository.blogPostColumns/scanBlogPost (unexported there) — see comment
// on exportMasterclassColumns above.
const exportBlogPostColumns = "id, slug, title, cover_image, category, tags, author, published_at, content, status, created_at, updated_at"

func scanExportBlogPost(row pgx.CollectableRow) (model.BlogPost, error) {
	var m model.BlogPost
	err := row.Scan(&m.ID, &m.Slug, &m.Title, &m.CoverImage, &m.Category, &m.Tags, &m.Author, &m.PublishedAt, &m.Content, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func blogPostArgs(m model.BlogPost) []any {
	return []any{m.ID, m.Slug, m.Title, m.CoverImage, m.Category, m.Tags, m.Author, m.PublishedAt, m.Content, m.Status}
}

func scanFAQItem(row pgx.CollectableRow) (model.FAQItem, error) {
	var m model.FAQItem
	err := row.Scan(&m.ID, &m.Question, &m.Answer, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func faqItemArgs(m model.FAQItem) []any {
	return []any{m.ID, m.Question, m.Answer, m.SortOrder}
}

func scanFeature(row pgx.CollectableRow) (model.Feature, error) {
	var m model.Feature
	err := row.Scan(&m.ID, &m.Page, &m.Icon, &m.Title, &m.Description, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func featureArgs(m model.Feature) []any {
	return []any{m.ID, m.Page, m.Icon, m.Title, m.Description, m.SortOrder}
}

func scanAboutItem(row pgx.CollectableRow) (model.AboutItem, error) {
	var m model.AboutItem
	err := row.Scan(&m.ID, &m.Badge, &m.Description, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func aboutItemArgs(m model.AboutItem) []any {
	return []any{m.ID, m.Badge, m.Description, m.SortOrder}
}

func scanSocialLink(row pgx.CollectableRow) (model.SocialLink, error) {
	var m model.SocialLink
	err := row.Scan(&m.ID, &m.Label, &m.Href, &m.Disclaimer, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
func socialLinkArgs(m model.SocialLink) []any {
	return []any{m.ID, m.Label, m.Href, m.Disclaimer, m.SortOrder}
}

func scanPageContent(row pgx.CollectableRow) (model.PageContent, error) {
	var m model.PageContent
	err := row.Scan(&m.Key, &m.Label, &m.Value, &m.Type, &m.UpdatedAt)
	return m, err
}
