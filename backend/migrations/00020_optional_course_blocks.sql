-- +goose Up
-- Splitting a course into CourseBlocks is now optional (it used to be
-- mandatory — every course had to have at least one block, even a blank
-- single one, purely so it had a card to render). The course itself now
-- carries its own card fields; CourseCatalogService uses them only when the
-- course has zero CourseBlock rows, and uses the real blocks (one card
-- each) otherwise. See model.Course's doc comment.
ALTER TABLE courses ADD COLUMN cover_image TEXT NOT NULL DEFAULT '';
ALTER TABLE courses ADD COLUMN lesson_count TEXT NOT NULL DEFAULT '';
ALTER TABLE courses ADD COLUMN time_length TEXT NOT NULL DEFAULT '';
ALTER TABLE courses ADD COLUMN price TEXT NOT NULL DEFAULT '';
ALTER TABLE courses ADD COLUMN display_style TEXT NOT NULL DEFAULT 'blue-beige';
ALTER TABLE courses ADD CONSTRAINT courses_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown'));

-- Every course so far has exactly one block (the old mandatory-block
-- world) except courses that were deliberately split into several — copy
-- each single-block course's block fields up onto the course itself, then
-- drop that now-redundant block row so the course becomes a plain
-- zero-block course under the new model. Courses with 2+ blocks are
-- untouched: they keep rendering one card per block, same as before.
UPDATE courses c
SET cover_image = b.block_cover, lesson_count = b.lesson_count, time_length = b.time_length,
    price = b.price, display_style = b.display_style
FROM course_blocks b
WHERE b.course_id = c.id
  AND (SELECT COUNT(*) FROM course_blocks b2 WHERE b2.course_id = c.id) = 1;

DELETE FROM course_blocks
WHERE course_id IN (
    SELECT course_id FROM course_blocks GROUP BY course_id HAVING COUNT(*) = 1
);

-- +goose Down
INSERT INTO course_blocks (course_id, block_name, description, block_cover, lesson_count, time_length, price, display_style, visible, sort_order)
SELECT id, '', '', cover_image, lesson_count, time_length, price, display_style, true, 0
FROM courses
WHERE id NOT IN (SELECT DISTINCT course_id FROM course_blocks);

ALTER TABLE courses DROP CONSTRAINT courses_display_style_check;
ALTER TABLE courses DROP COLUMN display_style;
ALTER TABLE courses DROP COLUMN price;
ALTER TABLE courses DROP COLUMN time_length;
ALTER TABLE courses DROP COLUMN lesson_count;
ALTER TABLE courses DROP COLUMN cover_image;
