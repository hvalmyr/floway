-- +goose Up
-- Lets the portfolio page filter by object type/format (п. 6 ТЗ decor-site) —
-- both nullable/defaulted so existing photos (school's homepage carousel,
-- untagged) keep working unfiltered.
ALTER TABLE gallery_photos ADD COLUMN object_type_id BIGINT REFERENCES object_types(id) ON DELETE SET NULL;
ALTER TABLE gallery_photos ADD COLUMN format TEXT NOT NULL DEFAULT '';
ALTER TABLE gallery_photos ADD CONSTRAINT gallery_photos_format_check CHECK (format IN ('', 'season', 'event'));

-- +goose Down
ALTER TABLE gallery_photos DROP CONSTRAINT gallery_photos_format_check;
ALTER TABLE gallery_photos DROP COLUMN format;
ALTER TABLE gallery_photos DROP COLUMN object_type_id;
