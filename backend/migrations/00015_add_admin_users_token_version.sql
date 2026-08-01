-- +goose Up
-- Bumped on logout (and could be bumped manually to kick a compromised
-- session) so an already-issued JWT stops being accepted before its natural
-- expiry — architecture review finding #16: logout only ever cleared the
-- client-side cookie, a copied/leaked token stayed valid for the rest of its
-- 12h TTL regardless.
ALTER TABLE admin_users ADD COLUMN token_version INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE admin_users DROP COLUMN token_version;
