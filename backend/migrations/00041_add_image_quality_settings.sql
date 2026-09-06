-- +goose Up
-- Photo compression settings, editable from the same "Тексты сайта" admin
-- page as everything else in page_content. 'number' is a new type — the
-- admin UI renders it as a bounded numeric input instead of a text/markdown
-- field; the service layer (PageContentService.Update) enforces the 1-100
-- range for these three specific keys.
--
-- Values picked to match what was already hardcoded for the hero image
-- (see UiHeroPicture.vue before this migration): avif needs a much lower
-- number than webp/jpeg for an equivalent look, and mobile can go lower
-- than desktop since screens are smaller and viewed further from the eye.
ALTER TABLE page_content DROP CONSTRAINT page_content_type_check;
ALTER TABLE page_content ADD CONSTRAINT page_content_type_check
    CHECK (type IN ('text', 'image', 'icon', 'number'));

INSERT INTO page_content (key, label, value, type) VALUES
    ('image_quality_desktop', 'Качество фото на десктопе, webp/jpeg (1–100)', '80', 'number'),
    ('image_quality_mobile', 'Качество фото на мобильных, webp/jpeg (1–100)', '55', 'number'),
    ('image_quality_avif', 'Качество фото в формате AVIF (1–100)', '50', 'number');

-- +goose Down
DELETE FROM page_content WHERE key IN ('image_quality_desktop', 'image_quality_mobile', 'image_quality_avif');
ALTER TABLE page_content DROP CONSTRAINT page_content_type_check;
ALTER TABLE page_content ADD CONSTRAINT page_content_type_check CHECK (type IN ('text', 'image', 'icon'));
