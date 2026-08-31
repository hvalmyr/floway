-- +goose Up
-- Per-block homepage card color style, picked by the admin per course block
-- (see CourseCard.vue): background/text color pair drawn from the design
-- system's 4 standard colors (blue/beige/white/brown) — only the 3 that
-- combine into readable pairs are used here (white isn't part of any of the
-- 4 combos the brief specifies).
ALTER TABLE course_blocks ADD COLUMN display_style TEXT NOT NULL DEFAULT 'blue-beige';
ALTER TABLE course_blocks ADD CONSTRAINT course_blocks_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown'));

-- Vertical-photo carousel between "Преимущества" and "О школе" on the
-- homepage — same flat id/image/sort_order shape as teachers.
CREATE TABLE gallery_photos (
    id         BIGSERIAL PRIMARY KEY,
    image      TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO page_content (key, label, value) VALUES
    ('home_features_heading', 'Главная — заголовок блока преимуществ', 'Почему стоит учиться в школе «ФлоВей»?'),
    ('home_features_lead', 'Главная — подзаголовок блока преимуществ', 'Мы создали школу, в которой удобно учиться, легко развиваться и получать реальные практические навыки флористики.');

-- +goose Down
DELETE FROM page_content WHERE key IN ('home_features_heading', 'home_features_lead');
DROP TABLE gallery_photos;
ALTER TABLE course_blocks DROP CONSTRAINT course_blocks_display_style_check;
ALTER TABLE course_blocks DROP COLUMN display_style;
