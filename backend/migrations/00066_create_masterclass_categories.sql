-- +goose Up
-- Lets an admin group masterclasses under tabs on the public page (e.g.
-- "Основные" vs. "Новогодние") instead of always showing one flat list —
-- same flat-dictionary shape as object_types, minus the "visible" flag
-- (every category is shown as its own tab; there's no draft state to hide).
CREATE TABLE masterclass_categories (
    id         BIGSERIAL PRIMARY KEY,
    slug       TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO masterclass_categories (slug, name, sort_order) VALUES
    ('osnovnye', 'Основные', 0),
    ('novogodnie', 'Новогодние', 1);

-- Nullable + ON DELETE SET NULL (not CASCADE, unlike landing_pages/
-- object_type_id): deleting a category should never take masterclasses with
-- it, same reasoning as courses/course_blocks.custom_display_style_id.
ALTER TABLE masterclasses ADD COLUMN category_id BIGINT REFERENCES masterclass_categories(id) ON DELETE SET NULL;

UPDATE masterclasses SET category_id = (SELECT id FROM masterclass_categories WHERE slug = 'osnovnye');

-- +goose Down
ALTER TABLE masterclasses DROP COLUMN category_id;
DROP TABLE masterclass_categories;
