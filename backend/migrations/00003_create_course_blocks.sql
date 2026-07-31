-- +goose Up
CREATE TABLE course_blocks (
    id             BIGSERIAL PRIMARY KEY,
    course_id      BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title          TEXT NOT NULL,
    lessons_count  INTEGER NOT NULL DEFAULT 0,
    hours          INTEGER NOT NULL DEFAULT 0,
    price          INTEGER NOT NULL DEFAULT 0,
    sort_order     INTEGER NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_course_blocks_course_id ON course_blocks(course_id);

-- +goose Down
DROP TABLE course_blocks;
