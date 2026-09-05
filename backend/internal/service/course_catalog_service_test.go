package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeCourseLookupRepo struct {
	courses map[string]model.Course
}

func (f *fakeCourseLookupRepo) FindBySlug(ctx context.Context, slug string) (model.Course, error) {
	course, ok := f.courses[slug]
	if !ok {
		return model.Course{}, service.ErrNotFound
	}
	return course, nil
}

type fakeCourseBlockListRepo struct {
	blocksByCourseID map[int64][]model.CourseBlock
}

func (f *fakeCourseBlockListRepo) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error) {
	return f.blocksByCourseID[courseID], nil
}

// ptr64 lets test tables write a literal int64 where model.Lesson now wants
// *int64 (CourseBlockID/CourseID) — a lesson belongs to exactly one parent,
// see model.Lesson's doc comment.
func ptr64(v int64) *int64 { return &v }

type fakeLessonBatchRepo struct {
	calls           int
	lastBlockIDs    []int64
	lessonsByBlock  map[int64][]model.Lesson
	courseCalls     int
	lastCourseID    int64
	lessonsByCourse map[int64][]model.Lesson
}

func (f *fakeLessonBatchRepo) ListByCourseBlockIDs(ctx context.Context, blockIDs []int64) ([]model.Lesson, error) {
	f.calls++
	f.lastBlockIDs = blockIDs
	var out []model.Lesson
	for _, id := range blockIDs {
		out = append(out, f.lessonsByBlock[id]...)
	}
	return out, nil
}

func (f *fakeLessonBatchRepo) ListByCourseID(ctx context.Context, courseID int64) ([]model.Lesson, error) {
	f.courseCalls++
	f.lastCourseID = courseID
	return f.lessonsByCourse[courseID], nil
}

type fakeCourseSectionListRepo struct {
	sections []model.CourseSection
}

func (f *fakeCourseSectionListRepo) List(ctx context.Context) ([]model.CourseSection, error) {
	return f.sections, nil
}

type fakeCourseBatchBySectionRepo struct {
	calls              int
	lastSectionIDs     []int64
	coursesBySectionID map[int64][]model.Course
}

func (f *fakeCourseBatchBySectionRepo) ListBySectionIDs(ctx context.Context, sectionIDs []int64) ([]model.Course, error) {
	f.calls++
	f.lastSectionIDs = sectionIDs
	var out []model.Course
	for _, id := range sectionIDs {
		out = append(out, f.coursesBySectionID[id]...)
	}
	return out, nil
}

type fakeCourseFAQListRepo struct {
	itemsByCourseID map[int64][]model.CourseFAQItem
}

func (f *fakeCourseFAQListRepo) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseFAQItem, error) {
	return f.itemsByCourseID[courseID], nil
}

type fakeCourseBlockBatchRepo struct {
	calls            int
	lastCourseIDs    []int64
	blocksByCourseID map[int64][]model.CourseBlock
}

func (f *fakeCourseBlockBatchRepo) ListByCourseIDs(ctx context.Context, courseIDs []int64) ([]model.CourseBlock, error) {
	f.calls++
	f.lastCourseIDs = courseIDs
	var out []model.CourseBlock
	for _, id := range courseIDs {
		out = append(out, f.blocksByCourseID[id]...)
	}
	return out, nil
}

func TestCourseCatalogService_GetFullBySlug_AssemblesBlocksAndLessons(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"osnovy": {ID: 1, Slug: "osnovy", Name: "Основы", Visible: true},
	}}
	blocks := &fakeCourseBlockListRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		1: {
			{ID: 10, CourseID: 1, BlockName: "Букеты", Visible: true},
			{ID: 11, CourseID: 1, BlockName: "Композиции", Visible: true},
		},
	}}
	lessons := &fakeLessonBatchRepo{lessonsByBlock: map[int64][]model.Lesson{
		10: {{ID: 100, CourseBlockID: ptr64(10), Name: "Занятие 1"}},
		11: {{ID: 110, CourseBlockID: ptr64(11), Name: "Занятие 1"}, {ID: 111, CourseBlockID: ptr64(11), Name: "Занятие 2"}},
	}}
	faqItems := &fakeCourseFAQListRepo{itemsByCourseID: map[int64][]model.CourseFAQItem{
		1: {{ID: 1000, CourseID: 1, Question: "Сколько длится курс?", Answer: "3 месяца"}},
	}}
	svc := service.NewCourseCatalogService(nil, courses, nil, blocks, nil, lessons, faqItems)

	detail, err := svc.GetFullBySlug(context.Background(), "osnovy")

	require.NoError(t, err)
	assert.Equal(t, "Основы", detail.Name)
	assert.Equal(t, 2, detail.BlockCount)
	require.Len(t, detail.Blocks, 2)
	assert.Equal(t, "Букеты", detail.Blocks[0].BlockName)
	require.Len(t, detail.Blocks[0].Lessons, 1)
	assert.Equal(t, "Композиции", detail.Blocks[1].BlockName)
	require.Len(t, detail.Blocks[1].Lessons, 2)
	require.Len(t, detail.FAQItems, 1)
	assert.Equal(t, "Сколько длится курс?", detail.FAQItems[0].Question)

	// The whole point: one batched lesson lookup, not one per block.
	assert.Equal(t, 1, lessons.calls)
	assert.ElementsMatch(t, []int64{10, 11}, lessons.lastBlockIDs)
}

// Regression test: a block with zero lessons must come back as an empty
// slice, not nil — a nil Lessons field serializes to JSON `null`, which
// crashed the frontend's `block.lessons.length` (see [slug].vue).
func TestCourseCatalogService_GetFullBySlug_BlockWithNoLessonsGetsEmptySlice(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"svadebnaya": {ID: 2, Slug: "svadebnaya", Name: "Свадебная флористика", Visible: true},
	}}
	blocks := &fakeCourseBlockListRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		2: {{ID: 21, CourseID: 2, BlockName: "Для флористов", Visible: true}},
	}}
	lessons := &fakeLessonBatchRepo{lessonsByBlock: map[int64][]model.Lesson{}}
	svc := service.NewCourseCatalogService(nil, courses, nil, blocks, nil, lessons, &fakeCourseFAQListRepo{})

	detail, err := svc.GetFullBySlug(context.Background(), "svadebnaya")

	require.NoError(t, err)
	require.Len(t, detail.Blocks, 1)
	assert.NotNil(t, detail.Blocks[0].Lessons)
	assert.Empty(t, detail.Blocks[0].Lessons)
}

func TestCourseCatalogService_GetFullBySlug_PropagatesNotFound(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{}}
	svc := service.NewCourseCatalogService(nil, courses, nil, &fakeCourseBlockListRepo{}, nil, &fakeLessonBatchRepo{}, nil)

	_, err := svc.GetFullBySlug(context.Background(), "does-not-exist")

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestCourseCatalogService_GetFullBySlug_CourseWithNoBlocks(t *testing.T) {
	// Splitting into blocks is optional (model.Course's doc comment) — a
	// course with zero CourseBlock rows still renders as exactly one card,
	// built from the course's own fields, and its lessons are its own
	// (course_id, not course_block_id — see model.Lesson's doc comment).
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"empty": {ID: 2, Slug: "empty", Name: "Empty", Visible: true, Price: "2000₽"},
	}}
	lessons := &fakeLessonBatchRepo{lessonsByCourse: map[int64][]model.Lesson{
		2: {{ID: 200, CourseID: ptr64(2), Name: "Занятие 1"}},
	}}
	svc := service.NewCourseCatalogService(nil, courses, nil, &fakeCourseBlockListRepo{}, nil, lessons, &fakeCourseFAQListRepo{})

	detail, err := svc.GetFullBySlug(context.Background(), "empty")

	require.NoError(t, err)
	assert.Equal(t, 1, detail.BlockCount)
	require.Len(t, detail.Blocks, 1)
	assert.Equal(t, "2000₽", detail.Blocks[0].Price)
	require.Len(t, detail.Blocks[0].Lessons, 1)
	assert.Equal(t, "Занятие 1", detail.Blocks[0].Lessons[0].Name)
	assert.Equal(t, 0, lessons.calls, "should never query lessons by block when there are no real blocks")
	assert.Equal(t, 1, lessons.courseCalls)
	assert.Equal(t, int64(2), lessons.lastCourseID)
}

func TestCourseCatalogService_GetFullBySlug_PropagatesLessonRepoError(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"osnovy": {ID: 1, Slug: "osnovy"},
	}}
	blocks := &fakeCourseBlockListRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		1: {{ID: 10, CourseID: 1}},
	}}
	svc := service.NewCourseCatalogService(nil, courses, nil, blocks, nil, failingLessonRepo{}, nil)

	_, err := svc.GetFullBySlug(context.Background(), "osnovy")

	require.Error(t, err)
}

type failingLessonRepo struct{}

func (failingLessonRepo) ListByCourseBlockIDs(ctx context.Context, blockIDs []int64) ([]model.Lesson, error) {
	return nil, errors.New("boom")
}

func (failingLessonRepo) ListByCourseID(ctx context.Context, courseID int64) ([]model.Lesson, error) {
	return nil, errors.New("boom")
}

func TestCourseCatalogService_ListSections_AssemblesCoursesAndBlocks(t *testing.T) {
	sections := &fakeCourseSectionListRepo{sections: []model.CourseSection{
		{ID: 1, Heading: "Курс «Основы флористики»", Visible: true},
		{ID: 2, Heading: "Профильные курсы", Visible: true},
	}}
	courses := &fakeCourseBatchBySectionRepo{coursesBySectionID: map[int64][]model.Course{
		1: {{ID: 10, SectionID: 1, Slug: "osnovy", Name: "Основы флористики", Visible: true}},
		2: {
			{ID: 20, SectionID: 2, Slug: "svadebnaya", Name: "Свадебная флористика", Visible: true},
			// course 21 intentionally has no blocks yet — a simple course
			// with its own card fields, not split into tracks.
			{ID: 21, SectionID: 2, Slug: "kommercheskaya", Name: "Коммерческая флористика", Visible: true, Price: "15000₽"},
		},
	}}
	blocks := &fakeCourseBlockBatchRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		10: {
			{ID: 100, CourseID: 10, BlockName: "Букеты", Visible: true},
			{ID: 101, CourseID: 10, BlockName: "Композиции", Visible: true},
		},
		20: {{ID: 200, CourseID: 20, Visible: true}},
		// course 21 intentionally has no blocks yet.
	}}
	svc := service.NewCourseCatalogService(sections, nil, courses, nil, blocks, nil, nil)

	result, err := svc.ListSections(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, "Курс «Основы флористики»", result[0].Heading)
	assert.Equal(t, 1, result[0].Length)
	require.Len(t, result[0].Courses, 1)
	assert.Equal(t, 2, result[0].Courses[0].BlockCount)
	assert.Len(t, result[0].Courses[0].Blocks, 2)

	assert.Equal(t, "Профильные курсы", result[1].Heading)
	assert.Equal(t, 2, result[1].Length)
	require.Len(t, result[1].Courses, 2)
	assert.Equal(t, 1, result[1].Courses[0].BlockCount)
	// course 21 has no real blocks, but still gets exactly one synthesized
	// card built from its own fields — not zero cards.
	require.Equal(t, 1, result[1].Courses[1].BlockCount)
	require.Len(t, result[1].Courses[1].Blocks, 1)
	assert.Equal(t, "15000₽", result[1].Courses[1].Blocks[0].Price)

	// One batched query for courses (all sections), one batched query for
	// blocks (all courses) — not one per section/course.
	assert.Equal(t, 1, courses.calls)
	assert.Equal(t, 1, blocks.calls)
	assert.ElementsMatch(t, []int64{1, 2}, courses.lastSectionIDs)
	assert.ElementsMatch(t, []int64{10, 20, 21}, blocks.lastCourseIDs)
}

func TestCourseCatalogService_ListSections_SingleCardCollapsesMultiBlockCourse(t *testing.T) {
	sections := &fakeCourseSectionListRepo{sections: []model.CourseSection{
		{ID: 1, Heading: "Секция", Visible: true},
	}}
	courses := &fakeCourseBatchBySectionRepo{coursesBySectionID: map[int64][]model.Course{
		1: {{ID: 10, SectionID: 1, Slug: "osnovy", Name: "Основы флористики", Visible: true, Price: "10000₽", SingleCard: true}},
	}}
	blocks := &fakeCourseBlockBatchRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		10: {
			{ID: 100, CourseID: 10, BlockName: "Букеты", Visible: true},
			{ID: 101, CourseID: 10, BlockName: "Композиции", Visible: true},
		},
	}}
	svc := service.NewCourseCatalogService(sections, nil, courses, nil, blocks, nil, nil)

	result, err := svc.ListSections(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Len(t, result[0].Courses, 1)
	// SingleCard is set, so the two real blocks are ignored in favor of one
	// card built from the course's own fields, same as a zero-block course.
	assert.Equal(t, 1, result[0].Courses[0].BlockCount)
	require.Len(t, result[0].Courses[0].Blocks, 1)
	assert.Equal(t, "10000₽", result[0].Courses[0].Blocks[0].Price)
	assert.Empty(t, result[0].Courses[0].Blocks[0].BlockName)
}

func TestCourseCatalogService_ListSections_NoSections(t *testing.T) {
	sections := &fakeCourseSectionListRepo{}
	courses := &fakeCourseBatchBySectionRepo{}
	svc := service.NewCourseCatalogService(sections, nil, courses, nil, &fakeCourseBlockBatchRepo{}, nil, nil)

	result, err := svc.ListSections(context.Background())

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, courses.calls, "should not query courses at all when there are no sections")
}

func TestCourseCatalogService_ListSections_DropsHiddenSectionCourseAndBlock(t *testing.T) {
	sections := &fakeCourseSectionListRepo{sections: []model.CourseSection{
		{ID: 1, Heading: "Видимая секция", Visible: true},
		{ID: 2, Heading: "Скрытая секция", Visible: false},
	}}
	courses := &fakeCourseBatchBySectionRepo{coursesBySectionID: map[int64][]model.Course{
		1: {
			{ID: 10, SectionID: 1, Slug: "vidimyi", Name: "Видимый курс", Visible: true},
			{ID: 11, SectionID: 1, Slug: "skrytyi", Name: "Скрытый курс", Visible: false},
		},
	}}
	blocks := &fakeCourseBlockBatchRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		10: {
			{ID: 100, CourseID: 10, BlockName: "Видимый блок", Visible: true},
			{ID: 101, CourseID: 10, BlockName: "Скрытый блок", Visible: false},
		},
	}}
	svc := service.NewCourseCatalogService(sections, nil, courses, nil, blocks, nil, nil)

	result, err := svc.ListSections(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1, "the hidden section should be dropped entirely")
	assert.Equal(t, "Видимая секция", result[0].Heading)
	require.Len(t, result[0].Courses, 1, "the hidden course should be dropped entirely")
	assert.Equal(t, "Видимый курс", result[0].Courses[0].Name)
	require.Len(t, result[0].Courses[0].Blocks, 1, "the hidden block should be dropped")
	assert.Equal(t, "Видимый блок", result[0].Courses[0].Blocks[0].BlockName)
	assert.Equal(t, 1, result[0].Courses[0].BlockCount)
}

func TestCourseCatalogService_GetFullBySlug_HiddenCourseIsNotFound(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"hidden": {ID: 3, Slug: "hidden", Name: "Скрытый курс", Visible: false},
	}}
	svc := service.NewCourseCatalogService(nil, courses, nil, &fakeCourseBlockListRepo{}, nil, &fakeLessonBatchRepo{}, nil)

	_, err := svc.GetFullBySlug(context.Background(), "hidden")

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestCourseCatalogService_GetFullBySlug_DropsHiddenBlocks(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"osnovy": {ID: 1, Slug: "osnovy", Name: "Основы", Visible: true},
	}}
	blocks := &fakeCourseBlockListRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		1: {
			{ID: 10, CourseID: 1, BlockName: "Видимый", Visible: true},
			{ID: 11, CourseID: 1, BlockName: "Скрытый", Visible: false},
		},
	}}
	lessons := &fakeLessonBatchRepo{lessonsByBlock: map[int64][]model.Lesson{
		10: {{ID: 100, CourseBlockID: ptr64(10), Name: "Занятие 1"}},
		11: {{ID: 110, CourseBlockID: ptr64(11), Name: "Занятие 1"}},
	}}
	svc := service.NewCourseCatalogService(nil, courses, nil, blocks, nil, lessons, &fakeCourseFAQListRepo{})

	detail, err := svc.GetFullBySlug(context.Background(), "osnovy")

	require.NoError(t, err)
	assert.Equal(t, 1, detail.BlockCount)
	require.Len(t, detail.Blocks, 1)
	assert.Equal(t, "Видимый", detail.Blocks[0].BlockName)
	assert.ElementsMatch(t, []int64{10}, lessons.lastBlockIDs, "hidden block's lessons should never be queried")
}
