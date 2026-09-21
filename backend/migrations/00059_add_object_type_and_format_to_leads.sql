-- +goose Up
-- Fields the decor site's lead form needs (object type, format, ad-tracking
-- params) that the school's leads never populate — all nullable/defaulted so
-- existing rows are unaffected. object_type_id is nullable (not every lead
-- names a type — trial-lesson/gift-certificate leads from the school never
-- will), unlike landing_pages.object_type_id which is NOT NULL because every
-- landing page always belongs to exactly one type.
ALTER TABLE leads ADD COLUMN object_type_id BIGINT REFERENCES object_types(id) ON DELETE SET NULL;
ALTER TABLE leads ADD COLUMN format TEXT NOT NULL DEFAULT '';
ALTER TABLE leads ADD CONSTRAINT leads_format_check CHECK (format IN ('', 'season', 'event'));
ALTER TABLE leads ADD COLUMN utm_source TEXT NOT NULL DEFAULT '';
ALTER TABLE leads ADD COLUMN utm_medium TEXT NOT NULL DEFAULT '';
ALTER TABLE leads ADD COLUMN utm_campaign TEXT NOT NULL DEFAULT '';
ALTER TABLE leads ADD COLUMN utm_content TEXT NOT NULL DEFAULT '';
ALTER TABLE leads ADD COLUMN utm_term TEXT NOT NULL DEFAULT '';
ALTER TABLE leads ADD COLUMN yclid TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE leads DROP COLUMN yclid;
ALTER TABLE leads DROP COLUMN utm_term;
ALTER TABLE leads DROP COLUMN utm_content;
ALTER TABLE leads DROP COLUMN utm_campaign;
ALTER TABLE leads DROP COLUMN utm_medium;
ALTER TABLE leads DROP COLUMN utm_source;
ALTER TABLE leads DROP CONSTRAINT leads_format_check;
ALTER TABLE leads DROP COLUMN format;
ALTER TABLE leads DROP COLUMN object_type_id;
