-- +goose Up
-- Two more color pairs (blue-on-brown, brown-on-blue) alongside the
-- existing 4 combos for course blocks/courses/blog posts — see
-- migrations 00017/00020/00048 for the original set.
ALTER TABLE course_blocks DROP CONSTRAINT course_blocks_display_style_check;
ALTER TABLE course_blocks ADD CONSTRAINT course_blocks_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown', 'blue-brown', 'brown-blue'));

ALTER TABLE courses DROP CONSTRAINT courses_display_style_check;
ALTER TABLE courses ADD CONSTRAINT courses_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown', 'blue-brown', 'brown-blue'));

ALTER TABLE blog_posts DROP CONSTRAINT blog_posts_display_style_check;
ALTER TABLE blog_posts ADD CONSTRAINT blog_posts_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown', 'blue-brown', 'brown-blue'));

-- +goose Down
ALTER TABLE course_blocks DROP CONSTRAINT course_blocks_display_style_check;
ALTER TABLE course_blocks ADD CONSTRAINT course_blocks_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown'));

ALTER TABLE courses DROP CONSTRAINT courses_display_style_check;
ALTER TABLE courses ADD CONSTRAINT courses_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown'));

ALTER TABLE blog_posts DROP CONSTRAINT blog_posts_display_style_check;
ALTER TABLE blog_posts ADD CONSTRAINT blog_posts_display_style_check
    CHECK (display_style IN ('blue-beige', 'brown-beige', 'beige-blue', 'beige-brown'));
