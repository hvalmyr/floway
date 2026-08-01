-- +goose Up
-- 'text' — обычная строка/абзац (может содержать markdown, это решает
-- фронтенд конкретного места, не тип в БД); 'image' — значение это
-- относительный URL загруженной через Garage картинки (или '', тогда
-- фронтенд показывает свою заглушку). Ограничение только для админки —
-- какой виджет показать (textarea vs загрузка файла).
ALTER TABLE page_content ADD COLUMN type TEXT NOT NULL DEFAULT 'text';
ALTER TABLE page_content ADD CONSTRAINT page_content_type_check CHECK (type IN ('text', 'image'));

INSERT INTO page_content (key, label, value, type) VALUES
    ('home_hero_image', 'Главная — фото Hero', '', 'image'),
    ('home_trial_image', 'Главная — фото блока пробного занятия', '', 'image'),
    ('masterclasses_hero_image', 'Мастер-классы — фото Hero', '', 'image'),
    ('contacts_office_image', 'Контакты — фото школы', '', 'image');

-- +goose Down
DELETE FROM page_content WHERE type = 'image';
ALTER TABLE page_content DROP CONSTRAINT page_content_type_check;
ALTER TABLE page_content DROP COLUMN type;
