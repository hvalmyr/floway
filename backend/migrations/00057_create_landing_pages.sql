-- +goose Up
-- One admin-creatable/copyable/renameable/hideable/deletable page per
-- object type (see model.LandingPage's doc comment for why this isn't
-- built on courses/course_sections). faq_title/faq_description/faq_visible
-- follow Course's inline-FAQ-fields-on-the-parent-row convention rather than
-- a separate settings table, since page_faq_settings' fixed page-name CHECK
-- doesn't fit a dynamically admin-created set of pages.
CREATE TABLE landing_pages (
    id                BIGSERIAL PRIMARY KEY,
    object_type_id    BIGINT NOT NULL REFERENCES object_types(id) ON DELETE CASCADE,
    slug              TEXT NOT NULL UNIQUE,
    h1                TEXT NOT NULL,
    meta_title        TEXT NOT NULL DEFAULT '',
    meta_description  TEXT NOT NULL DEFAULT '',
    faq_title         TEXT NOT NULL DEFAULT '',
    faq_description   TEXT NOT NULL DEFAULT '',
    faq_visible       BOOLEAN NOT NULL DEFAULT false,
    visible           BOOLEAN NOT NULL DEFAULT true,
    sort_order        INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_landing_pages_object_type_id ON landing_pages(object_type_id);

-- +goose Down
DROP TABLE landing_pages;
