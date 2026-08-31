-- +goose Up
-- Full redesign of the courses domain: homepage course listings become
-- admin-editable "course sections" (heading + description + list of
-- courses), each course optionally splits into named blocks, and each
-- block carries its own free-text price/duration/lesson-count plus its own
-- cover image and lesson list directly (no separate "curriculum" level).
-- The old courses/course_blocks/lessons tables are empty in every
-- environment this migration has run against, so this is a clean
-- replacement, not a data-preserving one.
DROP TABLE lessons;
DROP TABLE course_blocks;
DROP TABLE courses;

CREATE TABLE course_sections (
    id          BIGSERIAL PRIMARY KEY,
    heading     TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE courses (
    id          BIGSERIAL PRIMARY KEY,
    section_id  BIGINT NOT NULL REFERENCES course_sections(id) ON DELETE CASCADE,
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_courses_section_id ON courses(section_id);

-- lesson_count/time_length/price are free text ("7 занятий", "30 часов",
-- "38 500 ₽") so the admin can phrase them however the site copy needs —
-- including "was/now" discount prices directly in the price string — no
-- separate old_price column needed.
CREATE TABLE course_blocks (
    id           BIGSERIAL PRIMARY KEY,
    course_id    BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    block_name   TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    block_cover  TEXT NOT NULL DEFAULT '',
    lesson_count TEXT NOT NULL DEFAULT '',
    time_length  TEXT NOT NULL DEFAULT '',
    price        TEXT NOT NULL DEFAULT '',
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_course_blocks_course_id ON course_blocks(course_id);

CREATE TABLE lessons (
    id              BIGSERIAL PRIMARY KEY,
    course_block_id BIGINT NOT NULL REFERENCES course_blocks(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_lessons_course_block_id ON lessons(course_block_id);

-- +goose Down
DROP TABLE lessons;
DROP TABLE course_blocks;
DROP TABLE courses;
DROP TABLE course_sections;

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

CREATE TABLE course_blocks (
    id             BIGSERIAL PRIMARY KEY,
    course_id      BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title          TEXT NOT NULL,
    lessons_count  INTEGER NOT NULL DEFAULT 0,
    hours          INTEGER NOT NULL DEFAULT 0,
    price          INTEGER NOT NULL DEFAULT 0,
    old_price      INTEGER,
    sort_order     INTEGER NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_course_blocks_course_id ON course_blocks(course_id);

CREATE TABLE lessons (
    id              BIGSERIAL PRIMARY KEY,
    course_block_id BIGINT NOT NULL REFERENCES course_blocks(id) ON DELETE CASCADE,
    number          INTEGER NOT NULL,
    title           TEXT NOT NULL,
    topics          TEXT NOT NULL DEFAULT '',
    outcomes        TEXT NOT NULL DEFAULT '',
    duration_hours  INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_lessons_course_block_id ON lessons(course_block_id);
