-- +goose Up
-- Lets the admin edit the lead/application form's field labels and select
-- option text (ApplyForm.vue) without a deploy — deliberately NOT a field
-- constructor: which inputs exist and which options a select offers stay
-- fixed in code, only their displayed text moves into page_content, same
-- as every other freeform copy on the site.
INSERT INTO page_content (key, label, value) VALUES
    ('apply_form_name_label', 'Форма заявки — подпись поля «Имя»', 'Имя'),
    ('apply_form_name_placeholder', 'Форма заявки — плейсхолдер поля «Имя»', 'Как вас зовут'),
    ('apply_form_phone_label', 'Форма заявки — подпись поля «Телефон»', 'Номер телефона'),
    ('apply_form_email_label', 'Форма заявки — подпись поля «Почта»', 'Почта'),
    ('apply_form_email_placeholder', 'Форма заявки — плейсхолдер поля «Почта»', 'you@example.com'),
    ('apply_form_contact_method_label', 'Форма заявки — подпись селекта «Способ связи»', 'Как с вами связаться?'),
    ('apply_form_contact_method_call', 'Форма заявки — вариант связи «Звонок»', 'Позвоните мне'),
    ('apply_form_contact_method_telegram', 'Форма заявки — вариант связи «Telegram»', 'Напишите мне в Telegram'),
    ('apply_form_contact_method_whatsapp', 'Форма заявки — вариант связи «Whatsapp»', 'Напишите мне в Whatsapp'),
    ('apply_form_contact_method_max', 'Форма заявки — вариант связи «Max»', 'Напишите мне в Max'),
    ('apply_form_source_label', 'Форма заявки — подпись селекта «Откуда узнали»', 'Как вы о нас узнали?'),
    ('apply_form_source_referral', 'Форма заявки — источник «Рекомендация»', 'По рекомендации'),
    ('apply_form_source_ads', 'Форма заявки — источник «Реклама»', 'Реклама'),
    ('apply_form_source_internet', 'Форма заявки — источник «Интернет»', 'В интернете'),
    ('apply_form_source_social', 'Форма заявки — источник «Соцсети»', 'В социальных сетях'),
    ('apply_form_source_maps', 'Форма заявки — источник «Карты»', 'В картах'),
    ('apply_form_consent_prefix', 'Форма заявки — согласие, текст до ссылки', 'Отправляя форму, вы соглашаетесь с'),
    ('apply_form_consent_link_text', 'Форма заявки — согласие, текст ссылки на политику', 'политикой обработки персональных данных'),
    ('apply_form_consent_suffix', 'Форма заявки — согласие, текст после ссылки', 'и даёте согласие на обработку указанных персональных данных.'),
    ('apply_form_submit_default', 'Форма заявки — текст кнопки отправки (курс/мастер-класс)', 'Отправить заявку'),
    ('apply_form_submit_trial', 'Форма заявки — текст кнопки отправки (пробное занятие)', 'Записаться на занятие'),
    ('apply_form_success_title', 'Форма заявки — заголовок после отправки', 'заявка отправлена'),
    ('apply_form_success_message', 'Форма заявки — текст после отправки', 'Мы свяжемся с вами в ближайшее время удобным для вас способом.');

-- +goose Down
DELETE FROM page_content WHERE key LIKE 'apply_form_%';
