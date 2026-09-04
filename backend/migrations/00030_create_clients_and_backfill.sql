-- +goose Up
CREATE TABLE clients (
    id               BIGSERIAL PRIMARY KEY,
    name             TEXT NOT NULL,
    phone            TEXT NOT NULL DEFAULT '',
    phone_normalized TEXT NOT NULL DEFAULT '',
    email            TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_clients_phone_normalized ON clients(phone_normalized) WHERE phone_normalized <> '';
CREATE INDEX idx_clients_email_lower ON clients(lower(email)) WHERE email <> '';

ALTER TABLE leads ADD COLUMN client_id BIGINT REFERENCES clients(id);

-- Backfill: one Client per distinct normalized phone number. Every existing
-- lead has a non-blank phone (name+phone are required at creation — see
-- LeadService.Create), so phone is a reliable grouping key for 100% of rows.
-- Email is deliberately NOT used to merge groups here: it's optional/often
-- blank on old leads, and transitively merging groups connected only by
-- email is a graph problem not worth solving in a one-off SQL migration. A
-- historical client who used the same email but two different phone numbers
-- ends up with two Client rows instead of one — safe (no data loss), just
-- an imperfect merge that can be fixed by hand later if needed. Email dedup
-- DOES apply going forward for all new leads (see LeadService.Create).
WITH grouped AS (
    SELECT
        regexp_replace(phone, '[^0-9]', '', 'g')                      AS norm_phone,
        (array_agg(name  ORDER BY created_at DESC))[1]                AS latest_name,
        (array_agg(phone ORDER BY created_at DESC))[1]                AS latest_phone,
        COALESCE((array_agg(email ORDER BY created_at DESC)
                  FILTER (WHERE email <> ''))[1], '')                 AS latest_email
    FROM leads
    GROUP BY regexp_replace(phone, '[^0-9]', '', 'g')
),
inserted AS (
    INSERT INTO clients (name, phone, phone_normalized, email)
    SELECT latest_name, latest_phone, norm_phone, latest_email
    FROM grouped
    RETURNING id, phone_normalized
)
UPDATE leads
SET client_id = inserted.id
FROM inserted
WHERE regexp_replace(leads.phone, '[^0-9]', '', 'g') = inserted.phone_normalized;

ALTER TABLE leads ALTER COLUMN client_id SET NOT NULL;
CREATE INDEX idx_leads_client_id ON leads(client_id);

-- +goose Down
ALTER TABLE leads DROP COLUMN client_id;
DROP TABLE clients;
