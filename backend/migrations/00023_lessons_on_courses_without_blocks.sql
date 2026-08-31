-- +goose Up
-- Lets a course with no blocks have its own lesson list directly, instead of
-- requiring a real (if pointless) CourseBlock row just to unlock lesson
-- editing. A lesson now belongs to exactly one parent — a course_block (the
-- existing, unchanged path for courses split into blocks) OR a course
-- directly (new) — never both, never neither.
ALTER TABLE lessons ALTER COLUMN course_block_id DROP NOT NULL;
ALTER TABLE lessons ADD COLUMN course_id BIGINT REFERENCES courses(id) ON DELETE CASCADE;
ALTER TABLE lessons ADD CONSTRAINT lessons_exactly_one_parent CHECK (
    (course_id IS NOT NULL AND course_block_id IS NULL) OR
    (course_id IS NULL AND course_block_id IS NOT NULL)
);
CREATE INDEX idx_lessons_course_id ON lessons(course_id);

-- +goose Down
DROP INDEX IF EXISTS idx_lessons_course_id;
ALTER TABLE lessons DROP CONSTRAINT lessons_exactly_one_parent;
ALTER TABLE lessons DROP COLUMN course_id;
ALTER TABLE lessons ALTER COLUMN course_block_id SET NOT NULL;
