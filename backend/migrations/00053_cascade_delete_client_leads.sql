-- +goose Up
-- Deleting a client (new admin "Удалить клиента" action) needs its leads to
-- go with it — leads.client_id is NOT NULL, so it can't just be nulled out,
-- and the FK had no ON DELETE action (default NO ACTION), which would have
-- blocked deletion of virtually every client (every client originates from
-- at least one lead — see LeadService.Create). client_comments/reminders/
-- client_*_tags already cascade on client_id (migrations 00032-00034); this
-- makes leads consistent with them. leads itself has no dependents (nothing
-- references leads(id)), so this is the whole cascade chain.
ALTER TABLE leads DROP CONSTRAINT leads_client_id_fkey;
ALTER TABLE leads ADD CONSTRAINT leads_client_id_fkey
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE leads DROP CONSTRAINT leads_client_id_fkey;
ALTER TABLE leads ADD CONSTRAINT leads_client_id_fkey
    FOREIGN KEY (client_id) REFERENCES clients(id);
