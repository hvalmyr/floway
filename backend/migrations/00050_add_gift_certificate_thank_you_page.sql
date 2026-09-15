-- +goose Up
-- Fourth thank_you_pages variant (see 00047) — the gift-certificate request
-- form on pages/sertifikaty.vue used to borrow the "masterclass" variant's
-- copy/redirect; this gives it its own row so the content and the
-- lead's requestType both accurately reflect a certificate request.
INSERT INTO thank_you_pages (variant, title, subtitle, description) VALUES
    ('gift_certificate', 'Заявка на подарочный сертификат отправлена', 'Скоро свяжемся, чтобы оформить сертификат', 'Наш менеджер свяжется с вами в течение дня — уточнит детали и удобный способ оплаты. Сертификат можно оформить на конкретный курс или мастер-класс либо на произвольную сумму — расскажем обо всех вариантах.');

-- +goose Down
DELETE FROM thank_you_page_faq_items WHERE variant = 'gift_certificate';
DELETE FROM thank_you_page_photos WHERE variant = 'gift_certificate';
DELETE FROM thank_you_pages WHERE variant = 'gift_certificate';
