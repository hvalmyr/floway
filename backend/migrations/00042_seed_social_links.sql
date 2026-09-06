-- +goose Up
-- social_links (created empty in 00014) never got an admin UI until now —
-- seed it with the values that were hardcoded in frontend's
-- constants/contact-info.ts so the footer/contacts page doesn't go blank
-- once they switch from that constant to this table.
INSERT INTO social_links (label, href, disclaimer, sort_order) VALUES
    ('Telegram', 'https://t.me/floway', '', 0),
    ('VK', 'https://vk.com/floway', '', 1),
    ('Instagram', 'https://instagram.com/floway', 'Принадлежит компании Meta, признанной экстремистской организацией и запрещённой на территории РФ.', 2);

-- +goose Down
DELETE FROM social_links WHERE label IN ('Telegram', 'VK', 'Instagram');
