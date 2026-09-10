-- +goose Up
-- Same 4-color card style already used for course blocks/courses
-- (migrations 00017/00020) — reused for blog cards on /blog instead of a
-- separate palette.
ALTER TABLE blog_posts ADD COLUMN display_style TEXT NOT NULL DEFAULT 'blue-beige';
ALTER TABLE blog_posts ADD CONSTRAINT blog_posts_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown'));

-- +goose Down
ALTER TABLE blog_posts DROP CONSTRAINT blog_posts_display_style_check;
ALTER TABLE blog_posts DROP COLUMN display_style;
