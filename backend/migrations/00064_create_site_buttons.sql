-- +goose Up
-- Admin-editable text/style/link for the site's static standalone CTA
-- buttons (fixed destination, not a per-item link like a course card's
-- "go to course" or a form submit button) — see UiButton.vue for the two
-- variants. One row per button, keyed by a fixed slug the frontend looks
-- up via useSiteButtons(); rows are never created/deleted through the API,
-- only updated via PUT /api/v1/site-buttons/{key}.
CREATE TABLE site_buttons (
    key        TEXT PRIMARY KEY,
    label      TEXT NOT NULL,
    text       TEXT NOT NULL,
    variant    TEXT NOT NULL DEFAULT 'primary',
    url        TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO site_buttons (key, label, text, variant, url) VALUES
    ('home_hero_courses', 'Главная — Hero — «Курсы»', 'Курсы', 'primary', '/#courses'),
    ('home_hero_trial', 'Главная — Hero — «Пробное занятие»', 'Пробное занятие', 'outline', '/#trial'),
    ('sertifikaty_apply', 'Сертификаты — Hero — «Оставить заявку»', 'Оставить заявку', 'primary', '#apply'),
    ('sertifikaty_masterclasses', 'Сертификаты — Hero — «Мастер-классы»', 'Мастер-классы', 'outline', '/masterclasses'),
    ('course_detail_apply', 'Страница курса — «Оставить заявку»', 'Оставить заявку', 'primary', '#apply'),
    ('thank_you_home', 'Страница благодарности — «Вернуться на главную»', 'Вернуться на главную', 'primary', '/'),
    ('masterclasses_hero_list', 'Мастер-классы — Hero — «Мастер-классы»', 'Мастер-классы', 'primary', '#masterclasses-list'),
    ('masterclasses_hero_apply', 'Мастер-классы — Hero — «Оставить заявку»', 'Оставить заявку', 'outline', '#apply'),
    ('blog_back', 'Страница статьи блога — «Назад к блогу»', 'Назад к блогу', 'outline', '/blog');

-- +goose Down
DROP TABLE site_buttons;
