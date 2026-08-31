-- +goose Up
-- Lets the admin collapse a multi-block course into ONE homepage card
-- (using the course's own cover/lessonCount/timeLength/price/displayStyle,
-- same fields a zero-block course already renders from — see
-- syntheticBlock() in course_catalog_service.go) instead of one card per
-- block. Site-only: the admin's own course/block list pages are unaffected,
-- they keep listing every block as its own row regardless of this flag.
ALTER TABLE courses ADD COLUMN single_card BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE courses DROP COLUMN single_card;
