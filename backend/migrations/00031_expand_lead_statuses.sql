-- +goose Up
-- Flags leads whose status was auto-migrated from the old 3-value enum so
-- the admin can manually reclassify them (see the UPDATE below) instead of
-- guessing silently or building a whole separate review screen for it.
ALTER TABLE leads ADD COLUMN needs_status_review BOOLEAN NOT NULL DEFAULT false;

-- The old enum's 'closed' had no won/lost distinction. Default to
-- 'closed_won' (product decision) and flag every affected row for review
-- so conversion-rate numbers can be corrected by hand where needed.
UPDATE leads SET status = 'closed_won', needs_status_review = true WHERE status = 'closed';

CREATE INDEX idx_leads_needs_status_review ON leads(needs_status_review) WHERE needs_status_review = true;

-- +goose Down
-- Lossy by nature: any lead moved to a status that didn't exist in the old
-- 3-value enum (waiting_client/booked/postponed/closed_lost) has no
-- faithful old-schema equivalent. This Down is for local dev rollback only
-- — goose only ever runs "up" via the migrate compose service in
-- production (see docker-compose.yml).
UPDATE leads SET status = 'closed' WHERE status IN ('closed_won', 'closed_lost');
ALTER TABLE leads DROP COLUMN needs_status_review;
