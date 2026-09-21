-- +goose Up
-- Decor site's own thank-you-page variant (п. 8.2 ТЗ decor-site) — same
-- table/convention as the school's course/masterclass/trial_lesson/
-- gift_certificate rows (migration 00047), just a new row key. Placeholder
-- copy per the TZ's structure (менеджер свяжется в течение 30 минут →
-- бесплатный выезд и замеры → визуал и КП) — editable via
-- PUT /api/v1/thank-you-pages/decor without a deploy.
INSERT INTO thank_you_pages (variant, title, subtitle, description) VALUES
    ('decor', 'Заявка принята', 'Мы свяжемся с вами в течение 30 минут', 'Дальше — бесплатный выезд на объект и замеры, затем визуал оформления и коммерческое предложение в виде презентации. Если вопросы появятся раньше — пишите в любой удобный мессенджер.');

-- +goose Down
DELETE FROM thank_you_pages WHERE variant = 'decor';
