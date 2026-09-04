-- +goose Up
-- remind_at is a DATE, not a TIMESTAMPTZ: "remind me in N days" and the
-- "due today" badge are both date-granularity concepts, which sidesteps
-- time-of-day/timezone complexity entirely.
CREATE TABLE reminders (
    id           BIGSERIAL PRIMARY KEY,
    client_id    BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    remind_at    DATE NOT NULL,
    note         TEXT NOT NULL DEFAULT '',
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_reminders_due ON reminders(remind_at) WHERE completed_at IS NULL;
CREATE INDEX idx_reminders_client_id ON reminders(client_id);

-- +goose Down
DROP TABLE reminders;
