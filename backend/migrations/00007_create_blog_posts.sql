-- +goose Up
CREATE TABLE blog_posts (
    id           BIGSERIAL PRIMARY KEY,
    slug         TEXT NOT NULL UNIQUE,
    title        TEXT NOT NULL,
    cover_image  TEXT NOT NULL DEFAULT '',
    category     TEXT NOT NULL DEFAULT '',
    tags         TEXT[] NOT NULL DEFAULT '{}',
    author       TEXT NOT NULL DEFAULT '',
    published_at TIMESTAMPTZ,
    content      TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'draft',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE blog_posts;
