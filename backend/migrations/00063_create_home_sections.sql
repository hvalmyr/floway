-- +goose Up
-- Admin-controlled order/visibility for the homepage's blocks below Hero
-- (Hero itself always renders first and isn't reorderable/hideable — see
-- pages/index.vue). One row per block, keyed by a fixed slug the frontend
-- maps to a component; rows are never created/deleted through the API,
-- only reordered and shown/hidden via PUT /api/v1/home-sections/{id}.
CREATE TABLE home_sections (
    id         BIGSERIAL PRIMARY KEY,
    key        TEXT NOT NULL UNIQUE,
    visible    BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO home_sections (key, sort_order) VALUES
    ('features', 0),
    ('gallery', 1),
    ('courses', 2),
    ('gift_certificates', 3),
    ('trial', 4),
    ('about', 5),
    ('teachers', 6),
    ('reviews', 7),
    ('faq', 8);

-- +goose Down
DROP TABLE home_sections;
