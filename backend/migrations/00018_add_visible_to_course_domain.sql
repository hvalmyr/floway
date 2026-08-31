-- +goose Up
-- Lets an admin hide a section/course/block from the public site without
-- deleting it (and losing its blocks/lessons in the process). Existing rows
-- backfill to true (nothing was hidden before this existed).
ALTER TABLE course_sections ADD COLUMN visible BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE courses ADD COLUMN visible BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE course_blocks ADD COLUMN visible BOOLEAN NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE course_blocks DROP COLUMN visible;
ALTER TABLE courses DROP COLUMN visible;
ALTER TABLE course_sections DROP COLUMN visible;
