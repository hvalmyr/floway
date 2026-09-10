-- +goose Up
-- Square hero photo for the thank-you page's hero block (title/description/
-- CTAs left, 1:1 photo right — same layout as Hero.vue), one per variant.
ALTER TABLE thank_you_pages ADD COLUMN hero_image TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE thank_you_pages DROP COLUMN hero_image;
