-- +goose Up
-- No author column, by design: single-admin system, nothing to attribute.
CREATE TABLE client_comments (
    id         BIGSERIAL PRIMARY KEY,
    client_id  BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    text       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_client_comments_client_id ON client_comments(client_id, created_at);

-- +goose Down
DROP TABLE client_comments;
