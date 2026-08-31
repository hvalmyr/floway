-- +goose Up
-- The masterclasses page's features-block heading/lead ("Почему стоит
-- выбрать мастер-класс «ФлоВей»...") was hardcoded in
-- app/pages/masterclasses/index.vue, unlike the equivalent home-page
-- heading (home_features_heading/lead, added in migration 00017) — moving
-- it into page_content so it's admin-editable too, same as its home-page
-- counterpart.
INSERT INTO page_content (key, label, value) VALUES
    ('masterclasses_features_heading', 'Мастер-классы — заголовок блока преимуществ', 'Почему стоит выбрать мастер-класс “ФлоВей”'),
    ('masterclasses_features_lead', 'Мастер-классы — подзаголовок блока преимуществ', 'Разовое занятие без обязательств: приходите, когда удобно, и уходите с готовой работой в руках.');

-- +goose Down
DELETE FROM page_content WHERE key IN ('masterclasses_features_heading', 'masterclasses_features_lead');
