-- +goose Up
-- Content shown after a lead form submission (ApplyForm.vue) once it
-- navigates away to its own page instead of swapping its own success text
-- inline — one row per LeadRequestType variant. Never created/deleted
-- through the API (same convention as page_faq_settings/migration 00040),
-- only updated. Messenger icons and social-network links are deliberately
-- NOT columns here — they reuse the site's existing global contact channels
-- (page_content contact_*_url keys) and social_links table; the show_*
-- columns below only toggle whether that already-editable shared content
-- appears on a given variant.
CREATE TABLE thank_you_pages (
    variant            TEXT PRIMARY KEY,
    title              TEXT NOT NULL DEFAULT '',
    subtitle           TEXT NOT NULL DEFAULT '',
    description        TEXT NOT NULL DEFAULT '',
    show_messengers    BOOLEAN NOT NULL DEFAULT true,
    show_social_links  BOOLEAN NOT NULL DEFAULT true,
    show_blog_link     BOOLEAN NOT NULL DEFAULT false,
    blog_link_text     TEXT NOT NULL DEFAULT '',
    blog_link_url      TEXT NOT NULL DEFAULT '',
    show_carousel      BOOLEAN NOT NULL DEFAULT false,
    show_faq           BOOLEAN NOT NULL DEFAULT false,
    show_community     BOOLEAN NOT NULL DEFAULT false,
    community_text     TEXT NOT NULL DEFAULT '',
    community_url      TEXT NOT NULL DEFAULT '',
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO thank_you_pages (variant, title, subtitle, description) VALUES
    ('course', 'Заявка на курс принята', 'Совсем скоро вы сделаете первый шаг в мир флористики', 'Наш менеджер свяжется с вами в течение дня — расскажет о программе, расписании и оплате, ответит на все вопросы. Если хочется начать готовиться уже сейчас — посмотрите программу курса или почитайте истории наших выпускников.'),
    ('masterclass', 'Заявка на мастер-класс отправлена', 'Ждём вас в студии!', 'Мы свяжемся с вами, чтобы подтвердить дату и время, и пришлём всё нужное — адрес студии, что взять с собой и как удобнее добраться. Если вопросы появятся раньше — пишите в любой удобный для вас мессенджер, ответим быстро.'),
    ('trial_lesson', 'Вы записаны на пробное занятие', 'Совсем скоро увидимся в студии', 'Пробное занятие — это возможность познакомиться со школой без каких-либо обязательств. Менеджер свяжется с вами, чтобы согласовать удобное время и рассказать, как подготовиться. Понравится — расскажем, как перейти на полный курс на выгодных условиях.');

-- Optional photo carousel per variant — same flat id/image/sort_order shape
-- as gallery_photos/gift_certificate_carousel_photos, scoped by variant
-- instead of being its own table per page (mirrors features/page_faq_items).
CREATE TABLE thank_you_page_photos (
    id         BIGSERIAL PRIMARY KEY,
    variant    TEXT NOT NULL,
    image      TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Optional mini-FAQ per variant — same shape as page_faq_items, scoped by
-- variant instead of a static page identifier.
CREATE TABLE thank_you_page_faq_items (
    id         BIGSERIAL PRIMARY KEY,
    variant    TEXT NOT NULL,
    question   TEXT NOT NULL DEFAULT '',
    answer     TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE thank_you_page_faq_items;
DROP TABLE thank_you_page_photos;
DROP TABLE thank_you_pages;
