-- +goose Up
-- Search-result title/description, distinct from the on-page title/content —
-- both optional, falling back to the post's existing title/blank when empty
-- (see BlogPost.MetaTitle/MetaDescription's doc comment).
ALTER TABLE blog_posts ADD COLUMN meta_title TEXT NOT NULL DEFAULT '';
ALTER TABLE blog_posts ADD COLUMN meta_description TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE blog_posts DROP COLUMN meta_title;
ALTER TABLE blog_posts DROP COLUMN meta_description;
