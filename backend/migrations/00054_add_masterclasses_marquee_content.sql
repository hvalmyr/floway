-- +goose Up
-- Same pattern as the gift-certificate marquee (see 00029): the
-- masterclasses-page marquee's text/icon are admin-editable page_content
-- rows rather than hardcoded, with fallback defaults in
-- MasterclassesMarquee.vue covering the case where this migration hasn't
-- run yet.
INSERT INTO page_content (key, label, value, type) VALUES
    ('masterclasses_marquee_text', 'Бегущая строка мастер-классов — текст', 'Мастер-классы по флористике — запишитесь прямо сейчас', 'text'),
    ('masterclasses_marquee_icon', 'Бегущая строка мастер-классов — иконка', 'tulips', 'icon');

-- +goose Down
DELETE FROM page_content WHERE key IN ('masterclasses_marquee_text', 'masterclasses_marquee_icon');
