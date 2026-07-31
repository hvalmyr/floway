-- +goose Up
CREATE TABLE leads (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT NOT NULL,
    phone          TEXT NOT NULL,
    email          TEXT NOT NULL DEFAULT '',
    contact_method TEXT NOT NULL,
    source         TEXT NOT NULL,
    request_type   TEXT NOT NULL,
    related_id     BIGINT,
    status         TEXT NOT NULL DEFAULT 'new',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_leads_status ON leads(status);
CREATE INDEX idx_leads_request_type ON leads(request_type);

-- +goose Down
DROP TABLE leads;
