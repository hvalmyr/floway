-- +goose Up
INSERT INTO page_content (key, label, value) VALUES
    ('masterclasses_hero_title', 'Мастер-классы — заголовок Hero', 'Мастер-классы по флористике в свободном графике'),
    ('masterclasses_hero_lead', 'Мастер-классы — подзаголовок Hero', 'Мастер классы по флористике посвящены созданию разных видов букетов и флористических композиций. На каждом занятии вы создаете собственную работу и осваиваете новые приемы и навыки флористики.');

-- +goose Down
DELETE FROM page_content WHERE key IN ('masterclasses_hero_title', 'masterclasses_hero_lead');
