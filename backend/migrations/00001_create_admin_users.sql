-- +goose Up
CREATE TABLE admin_users (
    id            BIGSERIAL PRIMARY KEY,
    login         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE admin_users;
