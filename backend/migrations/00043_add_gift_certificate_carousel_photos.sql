-- +goose Up
-- Editable photo carousel on the gift-certificates page, right after the
-- hero block — same flat id/image/sort_order shape as gallery_photos (the
-- homepage carousel), just a separate table since it's a different slot on
-- a different page.
CREATE TABLE gift_certificate_carousel_photos (
    id         BIGSERIAL PRIMARY KEY,
    image      TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE gift_certificate_carousel_photos;
