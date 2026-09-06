-- +goose Up
-- Page-scoped FAQ block for the masterclasses and gift-certificate pages —
-- same title+description+visible+items shape as a course's own FAQ block
-- (migration 00039), but keyed by a fixed page identifier instead of a
-- course id, since these are static pages, not rows in a table with their
-- own admin CRUD. Deliberately separate from the global, unscoped faq_items
-- table (migration 00008, rendered once on the homepage) — that one predates
-- any page concept and stays untouched; "home" is intentionally not a valid
-- page here.
CREATE TABLE page_faq_settings (
    page        TEXT PRIMARY KEY CHECK (page IN ('masterclasses', 'gift_certificate')),
    title       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    visible     BOOLEAN NOT NULL DEFAULT false,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO page_faq_settings (page) VALUES ('masterclasses'), ('gift_certificate');

CREATE TABLE page_faq_items (
    id         BIGSERIAL PRIMARY KEY,
    page       TEXT NOT NULL REFERENCES page_faq_settings(page) ON DELETE CASCADE,
    question   TEXT NOT NULL,
    answer     TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX page_faq_items_page_idx ON page_faq_items (page);

-- +goose Down
DROP TABLE page_faq_items;
DROP TABLE page_faq_settings;
