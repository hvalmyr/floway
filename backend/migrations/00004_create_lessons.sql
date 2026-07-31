-- +goose Up
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

-- +goose Down
DROP TABLE lessons;
