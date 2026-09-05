-- +goose Up
-- Per-course FAQ block, shown after the apply form on the course page
-- (frontend/app/pages/courses/[slug].vue). Exactly one block per course —
-- its title/intro text live directly on courses, its Q&A items in their own
-- table (mirrors course_blocks: one course_id-scoped list, same reasoning
-- as every other "list of sub-items" entity in this schema). Deliberately
-- NOT the global faq_items table (migration 00008) — that one has no course
-- concept and is rendered once, on the homepage.
ALTER TABLE courses ADD COLUMN faq_title TEXT NOT NULL DEFAULT '';
ALTER TABLE courses ADD COLUMN faq_description TEXT NOT NULL DEFAULT '';
-- Defaults to hidden — an admin should opt in once the title/items are
-- actually filled in, rather than a brand-new course showing an empty FAQ
-- section.
ALTER TABLE courses ADD COLUMN faq_visible BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE course_faq_items (
    id         BIGSERIAL PRIMARY KEY,
    course_id  BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    question   TEXT NOT NULL,
    answer     TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX course_faq_items_course_id_idx ON course_faq_items (course_id);

-- +goose Down
DROP TABLE course_faq_items;
ALTER TABLE courses DROP COLUMN faq_visible;
ALTER TABLE courses DROP COLUMN faq_description;
ALTER TABLE courses DROP COLUMN faq_title;
