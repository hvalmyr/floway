-- +goose Up
INSERT INTO page_content (key, label, value) VALUES
    ('gift_certificate_hero_title', 'Подарочный мастер-класс — заголовок Hero', 'Подарите мастер-класс по флористике'),
    ('gift_certificate_hero_lead', 'Подарочный мастер-класс — подзаголовок Hero', 'Сертификат «Фловей» — оригинальный подарок для тех, кто любит цветы и творчество. Получатель сам выберет мастер-класс и удобную дату визита в школу.');

INSERT INTO page_content (key, label, value, type) VALUES
    ('gift_certificate_hero_image', 'Подарочный мастер-класс — фото Hero', '', 'image');

-- +goose Down
DELETE FROM page_content WHERE key IN ('gift_certificate_hero_title', 'gift_certificate_hero_lead', 'gift_certificate_hero_image');
