-- +goose Up
-- Повторяющиеся структурированные списки (иконка+заголовок+описание,
-- бейдж+описание), в отличие от page_content — есть собственная форма,
-- поэтому типизированные сущности, а не generic-блоки.
CREATE TABLE features (
    id          BIGSERIAL PRIMARY KEY,
    page        TEXT NOT NULL,
    icon        TEXT NOT NULL,
    title       TEXT NOT NULL,
    description TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_features_page ON features(page);

CREATE TABLE about_items (
    id          BIGSERIAL PRIMARY KEY,
    badge       TEXT NOT NULL,
    description TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE features;
DROP TABLE about_items;
