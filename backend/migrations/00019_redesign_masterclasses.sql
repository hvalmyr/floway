-- +goose Up
-- Not every masterclass needs two prices (group/individual) — that was the
-- common case, not a rule. Collapsing to one free-text price column (same
-- pattern as course_blocks.price) lets an admin phrase it however the
-- specific masterclass needs ("3000₽", "3000₽ или 4500₽", "от 2500₽"),
-- instead of the form always asking for a second price nobody fills in.
ALTER TABLE masterclasses RENAME COLUMN short_description TO description;
ALTER TABLE masterclasses RENAME COLUMN full_description TO description2;
ALTER TABLE masterclasses ADD COLUMN price TEXT NOT NULL DEFAULT '';
UPDATE masterclasses SET price = CASE
    WHEN price_individual > 0 THEN price_group || '₽ или ' || price_individual || '₽'
    WHEN price_group > 0 THEN price_group || '₽'
    ELSE ''
END;
ALTER TABLE masterclasses DROP COLUMN price_group;
ALTER TABLE masterclasses DROP COLUMN price_individual;
ALTER TABLE masterclasses DROP COLUMN price_description;

INSERT INTO masterclasses (slug, title, description, description2, ending_text, duration, price, cover_image, status)
VALUES
    (
        'krugliy-buket',
        'Круглый букет',
        'Круглый букет — это один из самых важных навыков в современной флористике и отличная отправная точка для знакомства с профессией. На этом мастер-классе вы научитесь создавать гармоничные круглые букеты, освоите основы спиральной техники и поймете, как правильно располагать цветы, чтобы композиция получалась объемной, аккуратной и устойчивой.',
        'Запишитесь на мастер-класс по созданию круглого букета и сделайте первый шаг в мир профессиональной флористики под руководством опытного преподавателя.',
        '3000 ₽ — практическое занятие, после которого собранный букет остается в школе. 4500 ₽ — с возможностью забрать готовый букет с собой.',
        '2-3 часа',
        '3000₽ или 4500₽',
        '',
        'active'
    ),
    (
        'raskidistiy-buket',
        'Раскидистый букет',
        'Раскидистый букет — один из самых эффектных и востребованных современных форматов флористики. На этом мастер-классе вы научитесь создавать легкие, воздушные композиции с естественной природной формой, которые выглядят объемно, стильно и профессионально.',
        'Запишитесь на мастер-класс по созданию раскидистого букета и научитесь собирать современные цветочные композиции, которые станут настоящим украшением любого интерьера или прекрасным подарком.',
        'После окончания занятия вы заберете собранный букет с собой.',
        '4 часа',
        '5500₽',
        '',
        'active'
    );

-- +goose Down
DELETE FROM masterclasses WHERE slug IN ('krugliy-buket', 'raskidistiy-buket');
ALTER TABLE masterclasses ADD COLUMN price_group INTEGER NOT NULL DEFAULT 0;
ALTER TABLE masterclasses ADD COLUMN price_individual INTEGER NOT NULL DEFAULT 0;
ALTER TABLE masterclasses ADD COLUMN price_description TEXT NOT NULL DEFAULT '';
ALTER TABLE masterclasses DROP COLUMN price;
ALTER TABLE masterclasses RENAME COLUMN description TO short_description;
ALTER TABLE masterclasses RENAME COLUMN description2 TO full_description;
