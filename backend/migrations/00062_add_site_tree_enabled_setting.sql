-- +goose Up
-- Admin toggle for whether AmbientTreeBackground.vue's animated 3D branch
-- renders on the school site — see layouts/default.vue and bare.vue's
-- `treeEnabled` gate. Off by default, matching the current live site
-- (migration 00061's predecessor, "perf(frontend): disable
-- AmbientTreeBackground WebGL scene", turned it off entirely for mobile
-- TBT/SEO) — admins can flip it back on from "Тексты сайта" once desired.
ALTER TABLE page_content DROP CONSTRAINT page_content_type_check;
ALTER TABLE page_content ADD CONSTRAINT page_content_type_check
    CHECK (type IN ('text', 'image', 'icon', 'number', 'boolean'));

INSERT INTO page_content (key, label, value, type) VALUES
    ('site_tree_enabled', 'Показывать анимированное дерево на сайте', 'false', 'boolean');

-- +goose Down
DELETE FROM page_content WHERE key = 'site_tree_enabled';
ALTER TABLE page_content DROP CONSTRAINT page_content_type_check;
ALTER TABLE page_content ADD CONSTRAINT page_content_type_check
    CHECK (type IN ('text', 'image', 'icon', 'number'));
