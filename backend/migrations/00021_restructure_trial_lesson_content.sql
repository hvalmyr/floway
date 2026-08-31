-- +goose Up
-- Splits the home page's trial-lesson block into 4 editable text keys
-- instead of the previous 5 (heading/lead for the section, a separate
-- fixed-label "Пробное занятие" heading, and one free-flowing description
-- replacing the old duration/price fields, which forced two separate
-- single-line values instead of one admin-editable paragraph).
UPDATE page_content SET key = 'trial_section_heading', value = 'Попробуйте флористику на практике'
    WHERE key = 'home_trial_heading';
UPDATE page_content SET key = 'trial_section_description'
    WHERE key = 'home_trial_description';
UPDATE page_content SET key = 'trial_heading'
    WHERE key = 'home_trial_lesson_title';

INSERT INTO page_content (key, label, value, type) VALUES
    ('trial_description', 'Главная — описание пробного занятия (абзацы)',
     E'Продолжительность: 2,5 часа.\n\nСтоимость: 3 000 ₽.', 'text');

DELETE FROM page_content WHERE key IN ('home_trial_duration', 'home_trial_price');

-- +goose Down
DELETE FROM page_content WHERE key = 'trial_description';

INSERT INTO page_content (key, label, value) VALUES
    ('home_trial_duration', 'Главная — длительность пробного занятия', 'Продолжительность: 2,5 часа'),
    ('home_trial_price', 'Главная — стоимость пробного занятия', 'Стоимость: 3 000 ₽');

UPDATE page_content SET key = 'home_trial_lesson_title' WHERE key = 'trial_heading';
UPDATE page_content SET key = 'home_trial_description' WHERE key = 'trial_section_description';
UPDATE page_content SET key = 'home_trial_heading', value = 'Попробуйте флористику на практике'
    WHERE key = 'trial_section_heading';
