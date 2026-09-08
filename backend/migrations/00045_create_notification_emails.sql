-- +goose Up
-- Admin-editable list of recipients for new-lead email notifications —
-- replaces the single NOTIFY_EMAIL_TO env var (see EmailNotifier), which
-- could only ever hold one address. The env var stays as a one-time seed on
-- startup for existing deployments (see cmd/api/main.go) but this table is
-- now the source of truth.
CREATE TABLE notification_emails (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE notification_emails;
