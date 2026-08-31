package service

import (
	"context"

	"floway-backend/internal/model"
)

// Narrow interfaces scoped to exactly what CourseCatalogService needs from
// each entity's repository — this service doesn't own CRUD for any of them
// (that's CourseSectionService/CourseService/CourseBlockService/LessonService's
// job), it only aggregates reads for the two public "catalog" endpoints.
type CourseSectionListRepository interface {
	List(ctx context.Context) ([]model.CourseSection, error)
}

type CourseLookupRepository interface {
	FindBySlug(ctx context.Context, slug string) (model.Course, error)
}

type CourseBatchBySectionRepository interface {
	ListBySectionIDs(ctx context.Context, sectionIDs []int64) ([]model.Course, error)
}

type CourseBlockListRepository interface {
	ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error)
}

type CourseBlockBatchRepository interface {
	ListByCourseIDs(ctx context.Context, courseIDs []int64) ([]model.CourseBlock, error)
}

type LessonBatchRepository interface {
	ListByCourseBlockIDs(ctx context.Context, courseBlockIDs []int64) ([]model.Lesson, error)
	// ListByCourseID is the course-without-blocks path — see
	// GetFullBySlug's zero-blocks branch and model.Lesson's doc comment.
	ListByCourseID(ctx context.Context, courseID int64) ([]model.Lesson, error)
}

// CourseCatalogService assembles the two public "read the whole tree in one
// shot" responses: the homepage section listing (sections -> courses ->
// blocks, no lesson text) and a single course page (course -> blocks ->
// lessons). Both are fixed-size batches of queries regardless of how many
// children exist, same reasoning as the old CourseDetailService this
// replaces.
type CourseCatalogService struct {
	sections     CourseSectionListRepository
	courseLookup CourseLookupRepository
	coursesBatch CourseBatchBySectionRepository
	blocks       CourseBlockListRepository
	blocksBatch  CourseBlockBatchRepository
	lessons      LessonBatchRepository
}

func NewCourseCatalogService(
	sections CourseSectionListRepository,
	courseLookup CourseLookupRepository,
	coursesBatch CourseBatchBySectionRepository,
	blocks CourseBlockListRepository,
	blocksBatch CourseBlockBatchRepository,
	lessons LessonBatchRepository,
) *CourseCatalogService {
	return &CourseCatalogService{
		sections:     sections,
		courseLookup: courseLookup,
		coursesBatch: coursesBatch,
		blocks:       blocks,
		blocksBatch:  blocksBatch,
		lessons:      lessons,
	}
}

// syntheticBlock builds a CourseBlock-shaped entry from a course's OWN card
// fields — used when the course has zero real CourseBlock rows (splitting
// into blocks is optional, see model.Course's doc comment), so the rest of
// this service (and the frontend, which only knows how to render a list of
// block-shaped cards) doesn't need a separate code path for "simple"
// courses. BlockName stays blank, same convention a real single block used
// under the old mandatory-block model — the public pages only show a
// block's name as a label when a course actually has more than one.
func syntheticBlock(course model.Course) model.CourseBlock {
	return model.CourseBlock{
		CourseID:     course.ID,
		Description:  course.Description,
		BlockCover:   course.CoverImage,
		LessonCount:  course.LessonCount,
		TimeLength:   course.TimeLength,
		Price:        course.Price,
		DisplayStyle: course.DisplayStyle,
		Visible:      true,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}
}

// GetFullBySlug is the public course-page aggregation: one course, its
// blocks, and each block's lessons — three queries total regardless of
// block count. A hidden course is reported as not found (same as if the
// slug didn't exist at all — "hide" means unpublish, not "leave off the
// homepage but still directly reachable"); hidden blocks are dropped from
// the result, same as ListSections does for the homepage cards.
func (s *CourseCatalogService) GetFullBySlug(ctx context.Context, slug string) (model.CourseWithBlocks, error) {
	course, err := s.courseLookup.FindBySlug(ctx, slug)
	if err != nil {
		return model.CourseWithBlocks{}, err
	}
	if !course.Visible {
		return model.CourseWithBlocks{}, ErrNotFound
	}

	allBlocks, err := s.blocks.ListByCourseID(ctx, course.ID)
	if err != nil {
		return model.CourseWithBlocks{}, err
	}
	blocks := make([]model.CourseBlock, 0, len(allBlocks))
	for _, block := range allBlocks {
		if block.Visible {
			blocks = append(blocks, block)
		}
	}
	if len(blocks) == 0 {
		block := syntheticBlock(course)
		// The course's own lessons (edited directly, no CourseBlock
		// involved — see model.Lesson's doc comment), attached to the
		// synthetic block so the rest of this response — and the frontend,
		// which only knows how to render a block's .lessons — needs no
		// separate code path for a blockless course.
		lessons, err := s.lessons.ListByCourseID(ctx, course.ID)
		if err != nil {
			return model.CourseWithBlocks{}, err
		}
		return model.CourseWithBlocks{
			Course:     course,
			BlockCount: 1,
			Blocks:     []model.CourseBlockWithLessons{{CourseBlock: block, Lessons: lessons}},
		}, nil
	}

	blockIDs := make([]int64, len(blocks))
	for i, block := range blocks {
		blockIDs[i] = block.ID
	}

	lessons, err := s.lessons.ListByCourseBlockIDs(ctx, blockIDs)
	if err != nil {
		return model.CourseWithBlocks{}, err
	}

	lessonsByBlock := make(map[int64][]model.Lesson, len(blocks))
	for _, lesson := range lessons {
		// Every lesson returned by ListByCourseBlockIDs has a non-nil
		// CourseBlockID — it's the WHERE clause that selected them.
		lessonsByBlock[*lesson.CourseBlockID] = append(lessonsByBlock[*lesson.CourseBlockID], lesson)
	}

	blocksWithLessons := make([]model.CourseBlockWithLessons, len(blocks))
	for i, block := range blocks {
		blocksWithLessons[i] = model.CourseBlockWithLessons{CourseBlock: block, Lessons: lessonsByBlock[block.ID]}
	}

	return model.CourseWithBlocks{Course: course, BlockCount: len(blocks), Blocks: blocksWithLessons}, nil
}

// ListSections is the public homepage aggregation: every course section,
// its courses, and each course's blocks (no lesson text — the homepage
// only needs a block's cover/lessonCount/timeLength/blockName) — three
// queries total regardless of section/course count. Hidden sections,
// courses, and blocks are all dropped at their respective level (a hidden
// section never even gets its courses fetched; a hidden course's blocks
// still count toward nothing since the course itself is gone).
func (s *CourseCatalogService) ListSections(ctx context.Context) ([]model.CourseSectionWithCourses, error) {
	allSections, err := s.sections.List(ctx)
	if err != nil {
		return nil, err
	}
	sections := make([]model.CourseSection, 0, len(allSections))
	for _, section := range allSections {
		if section.Visible {
			sections = append(sections, section)
		}
	}
	if len(sections) == 0 {
		return []model.CourseSectionWithCourses{}, nil
	}

	sectionIDs := make([]int64, len(sections))
	for i, section := range sections {
		sectionIDs[i] = section.ID
	}

	allCourses, err := s.coursesBatch.ListBySectionIDs(ctx, sectionIDs)
	if err != nil {
		return nil, err
	}
	courses := make([]model.Course, 0, len(allCourses))
	for _, course := range allCourses {
		if course.Visible {
			courses = append(courses, course)
		}
	}

	var allBlocks []model.CourseBlock
	if len(courses) > 0 {
		courseIDs := make([]int64, len(courses))
		for i, course := range courses {
			courseIDs[i] = course.ID
		}
		allBlocks, err = s.blocksBatch.ListByCourseIDs(ctx, courseIDs)
		if err != nil {
			return nil, err
		}
	}

	blocksByCourse := make(map[int64][]model.CourseBlock, len(courses))
	for _, block := range allBlocks {
		if block.Visible {
			blocksByCourse[block.CourseID] = append(blocksByCourse[block.CourseID], block)
		}
	}

	coursesBySection := make(map[int64][]model.CourseSummary, len(sections))
	for _, course := range courses {
		courseBlocks := blocksByCourse[course.ID]
		blockCount := len(courseBlocks)
		if blockCount == 0 || course.SingleCard {
			// No real blocks, or the admin chose to collapse a multi-block
			// course into one card — either way, one card from the course's
			// own fields instead (see syntheticBlock's doc comment).
			courseBlocks = []model.CourseBlock{syntheticBlock(course)}
			blockCount = 1
		}
		coursesBySection[course.SectionID] = append(coursesBySection[course.SectionID], model.CourseSummary{
			Course:     course,
			BlockCount: blockCount,
			Blocks:     courseBlocks,
		})
	}

	result := make([]model.CourseSectionWithCourses, len(sections))
	for i, section := range sections {
		sectionCourses := coursesBySection[section.ID]
		result[i] = model.CourseSectionWithCourses{
			CourseSection: section,
			Length:        len(sectionCourses),
			Courses:       sectionCourses,
		}
	}

	return result, nil
}
