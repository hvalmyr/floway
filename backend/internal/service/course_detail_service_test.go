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

type fakeLessonBatchRepo struct {
	calls          int
	lastBlockIDs   []int64
	lessonsByBlock map[int64][]model.Lesson
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

func TestCourseDetailService_GetFullBySlug_AssemblesModulesAndLessons(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"osnovy": {ID: 1, Slug: "osnovy", Title: "Основы"},
	}}
	blocks := &fakeCourseBlockListRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		1: {
			{ID: 10, CourseID: 1, Title: "Букеты"},
			{ID: 11, CourseID: 1, Title: "Композиции"},
		},
	}}
	lessons := &fakeLessonBatchRepo{lessonsByBlock: map[int64][]model.Lesson{
		10: {{ID: 100, CourseBlockID: 10, Title: "Занятие 1"}},
		11: {{ID: 110, CourseBlockID: 11, Title: "Занятие 1"}, {ID: 111, CourseBlockID: 11, Title: "Занятие 2"}},
	}}
	svc := service.NewCourseDetailService(courses, blocks, lessons)

	detail, err := svc.GetFullBySlug(context.Background(), "osnovy")

	require.NoError(t, err)
	assert.Equal(t, "Основы", detail.Title)
	require.Len(t, detail.Modules, 2)
	assert.Equal(t, "Букеты", detail.Modules[0].Title)
	require.Len(t, detail.Modules[0].Lessons, 1)
	assert.Equal(t, "Композиции", detail.Modules[1].Title)
	require.Len(t, detail.Modules[1].Lessons, 2)

	// The whole point: one batched lesson lookup, not one per block.
	assert.Equal(t, 1, lessons.calls)
	assert.ElementsMatch(t, []int64{10, 11}, lessons.lastBlockIDs)
}

func TestCourseDetailService_GetFullBySlug_PropagatesNotFound(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{}}
	svc := service.NewCourseDetailService(courses, &fakeCourseBlockListRepo{}, &fakeLessonBatchRepo{})

	_, err := svc.GetFullBySlug(context.Background(), "does-not-exist")

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestCourseDetailService_GetFullBySlug_CourseWithNoBlocks(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"empty": {ID: 2, Slug: "empty", Title: "Empty"},
	}}
	lessons := &fakeLessonBatchRepo{}
	svc := service.NewCourseDetailService(courses, &fakeCourseBlockListRepo{}, lessons)

	detail, err := svc.GetFullBySlug(context.Background(), "empty")

	require.NoError(t, err)
	assert.Empty(t, detail.Modules)
	assert.Equal(t, 0, lessons.calls, "should not query lessons at all when there are no blocks")
}

func TestCourseDetailService_GetFullBySlug_PropagatesLessonRepoError(t *testing.T) {
	courses := &fakeCourseLookupRepo{courses: map[string]model.Course{
		"osnovy": {ID: 1, Slug: "osnovy"},
	}}
	blocks := &fakeCourseBlockListRepo{blocksByCourseID: map[int64][]model.CourseBlock{
		1: {{ID: 10, CourseID: 1}},
	}}
	svc := service.NewCourseDetailService(courses, blocks, failingLessonRepo{})

	_, err := svc.GetFullBySlug(context.Background(), "osnovy")

	require.Error(t, err)
}

type failingLessonRepo struct{}

func (failingLessonRepo) ListByCourseBlockIDs(ctx context.Context, blockIDs []int64) ([]model.Lesson, error) {
	return nil, errors.New("boom")
}
