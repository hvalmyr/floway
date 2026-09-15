-- +goose Up
-- Admin-defined color styles (arbitrary bg/text hex pair) for courses that
-- need a one-off look the fixed 6-combo palette doesn't cover (e.g. a
-- New Year or autumn seasonal course) — alongside, not replacing, the
-- existing display_style enum on courses/course_blocks (migrations
-- 00017/00020). When custom_display_style_id is set it takes precedence
-- over display_style for that row's card color (see CourseCard.vue).
CREATE TABLE custom_display_styles (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    bg_color   TEXT NOT NULL,
    text_color TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE courses ADD COLUMN custom_display_style_id BIGINT REFERENCES custom_display_styles(id) ON DELETE SET NULL;
ALTER TABLE course_blocks ADD COLUMN custom_display_style_id BIGINT REFERENCES custom_display_styles(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE course_blocks DROP COLUMN custom_display_style_id;
ALTER TABLE courses DROP COLUMN custom_display_style_id;
DROP TABLE custom_display_styles;
