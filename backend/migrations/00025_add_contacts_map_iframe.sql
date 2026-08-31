-- +goose Up
-- Yandex map widget src on /contacts was hardcoded in the frontend
-- (unlike the rest of contact_* content, which already lived in
-- page_content but the frontend templates never actually read it —
-- fixed alongside this migration) — moving it here makes the map
-- editable from the admin panel like every other contact detail.
INSERT INTO page_content (key, label, value) VALUES
    ('contact_map_iframe_url', 'Контакты — ссылка на карту (iframe src)', 'https://yandex.ru/map-widget/v1/?um=constructor%3A059e5fd8df226009106acdf8b132a08a57922b5a9e616451b6a1d8e0e07eeba2&source=constructor');

-- +goose Down
DELETE FROM page_content WHERE key = 'contact_map_iframe_url';
