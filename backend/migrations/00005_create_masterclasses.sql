-- +goose Up
CREATE TABLE masterclasses (
    id                BIGSERIAL PRIMARY KEY,
    slug              TEXT NOT NULL UNIQUE,
    title             TEXT NOT NULL,
    short_description TEXT NOT NULL DEFAULT '',
    full_description  TEXT NOT NULL DEFAULT '',
    ending_text       TEXT NOT NULL DEFAULT '',
    duration          TEXT NOT NULL DEFAULT '',
    price_group       INTEGER NOT NULL DEFAULT 0,
    price_individual  INTEGER NOT NULL DEFAULT 0,
    price_description TEXT NOT NULL DEFAULT '',
    cover_image       TEXT NOT NULL DEFAULT '',
    status            TEXT NOT NULL DEFAULT 'active',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE masterclasses;
