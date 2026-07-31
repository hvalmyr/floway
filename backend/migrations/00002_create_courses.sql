-- +goose Up
CREATE TABLE courses (
    id               BIGSERIAL PRIMARY KEY,
    slug             TEXT NOT NULL UNIQUE,
    title            TEXT NOT NULL,
    short_description TEXT NOT NULL DEFAULT '',
    full_description TEXT NOT NULL DEFAULT '',
    status           TEXT NOT NULL DEFAULT 'active',
    cover_image      TEXT NOT NULL DEFAULT '',
    gallery          TEXT[] NOT NULL DEFAULT '{}',
    sort_order       INTEGER NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE courses;
