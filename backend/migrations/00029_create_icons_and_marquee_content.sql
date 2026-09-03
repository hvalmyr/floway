-- +goose Up
-- Uploaded SVG icon library — lets the admin attach a custom icon anywhere
-- a "built-in" icon key (FEATURE_ICONS on the frontend) was previously the
-- only option, without a deploy. `svg` holds the sanitized markup itself
-- (not a Garage/S3 URL): icons are small text, and storing the markup
-- directly means the public GET returns everything a page needs to render
-- one inline (recolorable via currentColor, unlike an <img src>), no extra
-- per-icon file fetch. Sanitization happens once, at upload time, in
-- IconService — every reader downstream (this table's content) is already
-- safe to inject with v-html.
CREATE TABLE icons (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    svg        TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- A feature/page_content "icon" value is either a bare FEATURE_ICONS key
-- (e.g. "gift") or "icon:<id>" referencing a row in the new table above —
-- one convention, parsed the same way everywhere (see AppIcon.vue) — so no
-- column/FK changes are needed on `features` itself (`icon TEXT` already
-- fits both forms).
ALTER TABLE page_content DROP CONSTRAINT page_content_type_check;
ALTER TABLE page_content ADD CONSTRAINT page_content_type_check CHECK (type IN ('text', 'image', 'icon'));

-- The gift-certificate marquee's text/icon were hardcoded in
-- GiftCertificateMarquee.vue (a fallback default masked the missing key —
-- see usePageContent's `text(key, fallback)`) — moving them here so both
-- are admin-editable, same as every other page_content-backed block.
INSERT INTO page_content (key, label, value, type) VALUES
    ('gift_certificate_marquee_text', 'Бегущая строка сертификатов — текст', 'Подарочные сертификаты — порадуйте близких цветами', 'text'),
    ('gift_certificate_marquee_icon', 'Бегущая строка сертификатов — иконка', 'gift', 'icon');

-- +goose Down
DELETE FROM page_content WHERE key IN ('gift_certificate_marquee_text', 'gift_certificate_marquee_icon');
ALTER TABLE page_content DROP CONSTRAINT page_content_type_check;
ALTER TABLE page_content ADD CONSTRAINT page_content_type_check CHECK (type IN ('text', 'image'));
DROP TABLE icons;
