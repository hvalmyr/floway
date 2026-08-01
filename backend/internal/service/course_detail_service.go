package service

import (
	"context"

	"floway-backend/internal/model"
)

// Narrow interfaces scoped to exactly what GetFullBySlug needs from each
// entity's repository — CourseDetailService doesn't own CRUD for any of
// them (that's CourseService/CourseBlockService/LessonService's job), it
// only aggregates reads.
type CourseLookupRepository interface {
	FindBySlug(ctx context.Context, slug string) (model.Course, error)
}

type CourseBlockListRepository interface {
	ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error)
}

type LessonBatchRepository interface {
	ListByCourseBlockIDs(ctx context.Context, courseBlockIDs []int64) ([]model.Lesson, error)
}

// CourseDetailService assembles the public "course page" response — one
// course, its blocks, and each block's lessons — in three queries total
// regardless of block count. Previously this composition lived in
// courseHandler.getFullBySlug directly (one query per block, untestable
// under the project's service-layer-only testing convention — see
// architecture review finding #2).
type CourseDetailService struct {
	courses CourseLookupRepository
	blocks  CourseBlockListRepository
	lessons LessonBatchRepository
}

func NewCourseDetailService(courses CourseLookupRepository, blocks CourseBlockListRepository, lessons LessonBatchRepository) *CourseDetailService {
	return &CourseDetailService{courses: courses, blocks: blocks, lessons: lessons}
}

func (s *CourseDetailService) GetFullBySlug(ctx context.Context, slug string) (model.CourseDetail, error) {
	course, err := s.courses.FindBySlug(ctx, slug)
	if err != nil {
		return model.CourseDetail{}, err
	}

	blocks, err := s.blocks.ListByCourseID(ctx, course.ID)
	if err != nil {
		return model.CourseDetail{}, err
	}
	if len(blocks) == 0 {
		return model.CourseDetail{Course: course, Modules: []model.CourseModule{}}, nil
	}

	blockIDs := make([]int64, len(blocks))
	for i, block := range blocks {
		blockIDs[i] = block.ID
	}

	lessons, err := s.lessons.ListByCourseBlockIDs(ctx, blockIDs)
	if err != nil {
		return model.CourseDetail{}, err
	}

	lessonsByBlock := make(map[int64][]model.Lesson, len(blocks))
	for _, lesson := range lessons {
		lessonsByBlock[lesson.CourseBlockID] = append(lessonsByBlock[lesson.CourseBlockID], lesson)
	}

	modules := make([]model.CourseModule, len(blocks))
	for i, block := range blocks {
		modules[i] = model.CourseModule{CourseBlock: block, Lessons: lessonsByBlock[block.ID]}
	}

	return model.CourseDetail{Course: course, Modules: modules}, nil
}
