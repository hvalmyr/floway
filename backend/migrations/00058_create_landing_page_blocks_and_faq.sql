-- +goose Up
-- Repeatable image+text sections within one landing page (photo examples,
-- format descriptions) — same shape as course_blocks minus the
-- lesson_count/time_length/price fields, which don't apply here.
CREATE TABLE landing_page_blocks (
    id              BIGSERIAL PRIMARY KEY,
    landing_page_id BIGINT NOT NULL REFERENCES landing_pages(id) ON DELETE CASCADE,
    image           TEXT NOT NULL DEFAULT '',
    text            TEXT NOT NULL DEFAULT '',
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_landing_page_blocks_landing_page_id ON landing_page_blocks(landing_page_id);

-- A landing page's own Q&A items — the block's title/description/visible
-- flag are plain LandingPage fields (see migration 00057), same convention
-- as course_faq_items/Course.
CREATE TABLE landing_page_faq_items (
    id              BIGSERIAL PRIMARY KEY,
    landing_page_id BIGINT NOT NULL REFERENCES landing_pages(id) ON DELETE CASCADE,
    question        TEXT NOT NULL,
    answer          TEXT NOT NULL,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_landing_page_faq_items_landing_page_id ON landing_page_faq_items(landing_page_id);

-- +goose Down
DROP TABLE landing_page_faq_items;
DROP TABLE landing_page_blocks;
