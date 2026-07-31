-- +goose Up
CREATE TABLE teachers (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    photo       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE teachers;
