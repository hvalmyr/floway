-- +goose Up
-- Lets the admin edit the lead-form title/subhead shown above the shared
-- ApplyForm.vue on each of the three pages that embed it with its own
-- wording (course page, masterclasses page, gift-certificate page) —
-- previously hardcoded per-page literals, see ApplyForm.vue's `title`/`lead`
-- props. The trial-lesson embed on the home page isn't included here: it
-- renders with an empty title on purpose (the section above it already has
-- its own heading).
INSERT INTO page_content (key, label, value) VALUES
    ('course_apply_form_title', 'Форма заявки на курсе — заголовок', 'Записаться на курс'),
    ('course_apply_form_lead', 'Форма заявки на курсе — описание', ''),
    ('masterclasses_apply_form_title', 'Форма заявки на мастер-классах — заголовок', 'Оставить заявку на мастер-класс'),
    ('masterclasses_apply_form_lead', 'Форма заявки на мастер-классах — описание', ''),
    ('gift_certificate_apply_form_title', 'Форма заявки на подарочных сертификатах — заголовок', 'Оставить заявку на подарочный сертификат'),
    ('gift_certificate_apply_form_lead', 'Форма заявки на подарочных сертификатах — описание', 'Свяжемся с вами и оформим сертификат.');

-- +goose Down
DELETE FROM page_content WHERE key IN (
    'course_apply_form_title',
    'course_apply_form_lead',
    'masterclasses_apply_form_title',
    'masterclasses_apply_form_lead',
    'gift_certificate_apply_form_title',
    'gift_certificate_apply_form_lead'
);
