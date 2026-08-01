-- +goose Up
-- Контакты/реквизиты были хардкодом во фронтенде (constants/contact-info.ts)
-- — переносим в page_content (одиночные поля) и новую таблицу social_links
-- (повторяющийся список, тот же паттерн, что features/about_items).
-- Реальные значения — из того же constants/contact-info.ts, 1-в-1, не
-- придумано заново.
INSERT INTO page_content (key, label, value, type) VALUES
    ('contact_phone', 'Контакты — телефон', '+7 985 226 19 48', 'text'),
    ('contact_email', 'Контакты — почта', 'floway-mos@mail.ru', 'text'),
    ('contact_telegram_url', 'Контакты — ссылка на Telegram', 'https://t.me/floway', 'text'),
    ('contact_whatsapp_url', 'Контакты — ссылка на Whatsapp', 'https://wa.me/79852261948', 'text'),
    ('contact_max_url', 'Контакты — ссылка на Max', '', 'text'),
    ('contact_address', 'Контакты — адрес', 'г. Москва, Новинский бульвар, 18Б', 'text'),
    ('contact_metro_stations', 'Контакты — станции метро (через запятую)', 'Смоленская, Баррикадная, Краснопресненская', 'text'),
    ('contact_directions', 'Контакты — как пройти (поддерживает markdown, абзацы через пустую строку)', 'Чтобы попасть к нам, зайдите в арку рядом с магазином «ВкусВилл». Во дворе поверните налево и идите вдоль сквера с детской площадкой. В конце сквера поверните направо — перед вами будет узкая дорожка, ведущая ко второму подъезду.

Подойдите ко второму подъезду, затем поверните налево. Наша дверь находится в самом углу жёлтого углового здания. На двери вы увидите вывеску «ФлоВей».

Если по пути возникнут вопросы или не получится найти вход, просто напишите или позвоните нам — мы с удовольствием подскажем и поможем найти дорогу.', 'text'),
    ('legal_entity_name', 'Реквизиты — юрлицо', 'ООО "МЭДЖИК ГАРДЕН"', 'text'),
    ('legal_inn', 'Реквизиты — ИНН', '9704145697', 'text'),
    ('legal_ogrn', 'Реквизиты — ОГРН', '1227700362779', 'text');

CREATE TABLE social_links (
    id         BIGSERIAL PRIMARY KEY,
    label      TEXT NOT NULL,
    href       TEXT NOT NULL,
    disclaimer TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO social_links (label, href, disclaimer, sort_order) VALUES
    ('Telegram', 'https://t.me/floway', '', 1),
    ('VK', 'https://vk.com/floway', '', 2),
    ('Instagram', 'https://instagram.com/floway', 'Instagram принадлежит компании Meta, признанной экстремистской организацией и запрещённой на территории РФ.', 3);

-- +goose Down
DROP TABLE social_links;
DELETE FROM page_content WHERE key IN (
    'contact_phone', 'contact_email', 'contact_telegram_url', 'contact_whatsapp_url',
    'contact_max_url', 'contact_address', 'contact_metro_stations', 'contact_directions',
    'legal_entity_name', 'legal_inn', 'legal_ogrn'
);
