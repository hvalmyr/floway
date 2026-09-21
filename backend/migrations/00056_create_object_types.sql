-- +goose Up
-- Shared dictionary of property types the decor site's landing pages, lead
-- form, and portfolio filter all draw from (see landing_pages.object_type_id,
-- leads.object_type_id, gallery_photos.object_type_id below) — one place to
-- rename/add/hide a type so it updates everywhere at once.
CREATE TABLE object_types (
    id         BIGSERIAL PRIMARY KEY,
    slug       TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    visible    BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE object_types;
