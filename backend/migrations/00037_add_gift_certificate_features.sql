-- +goose Up
-- The gift-certificate page's advantages block was hardcoded in
-- app/pages/sertifikaty.vue, unlike its home/masterclasses counterparts —
-- moving it into features/page_content so it's admin-editable too, same as
-- the other two (see migrations 00017 and 00028).
INSERT INTO page_content (key, label, value) VALUES
    ('gift_certificate_features_heading', 'Подарочные сертификаты — заголовок блока преимуществ', 'Почему это отличный подарок'),
    ('gift_certificate_features_lead', 'Подарочные сертификаты — подзаголовок блока преимуществ', 'Дарите не вещь, а впечатление — тёплый опыт создания своими руками.');

INSERT INTO features (page, icon, title, description, sort_order) VALUES
    ('gift_certificate', 'gift', 'Любая сумма или мастер-класс', 'Оформим сертификат на конкретный мастер-класс или на сумму — получатель выберет сам.', 0),
    ('gift_certificate', 'flex-start', 'Свободная дата записи', 'Получатель сам выберет удобный день и время — без привязки к конкретной дате.', 1),
    ('gift_certificate', 'checklist', 'Все материалы включены', 'Всё необходимое для мастер-класса уже включено в стоимость сертификата.', 2),
    ('gift_certificate', 'calendar-check', 'Долгий срок действия', 'Достаточно времени, чтобы выбрать удобный момент и записаться.', 3);

-- +goose Down
DELETE FROM page_content WHERE key IN ('gift_certificate_features_heading', 'gift_certificate_features_lead');
DELETE FROM features WHERE page = 'gift_certificate';
