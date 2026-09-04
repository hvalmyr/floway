-- +goose Up
-- Default matches the neutral gray chip both tag types already rendered
-- as (bg-gray-100) before colors existed, so existing tags don't jump to
-- an arbitrary color the first time this ships.
ALTER TABLE product_tags ADD COLUMN color TEXT NOT NULL DEFAULT '#f3f4f6';
ALTER TABLE client_type_tags ADD COLUMN color TEXT NOT NULL DEFAULT '#f3f4f6';

-- +goose Down
ALTER TABLE product_tags DROP COLUMN color;
ALTER TABLE client_type_tags DROP COLUMN color;
