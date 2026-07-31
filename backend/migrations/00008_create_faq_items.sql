-- +goose Up
CREATE TABLE faq_items (
    id         BIGSERIAL PRIMARY KEY,
    question   TEXT NOT NULL,
    answer     TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE faq_items;
