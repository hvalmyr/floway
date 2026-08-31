-- +goose Up
-- Lets the admin panel tell which specific course/masterclass a lead was
-- about without cross-referencing related_id back to that entity's own
-- table by hand — the frontend now sends the slug straight from whichever
-- page/card the visitor was looking at (see ApplyForm.vue's relatedSlug
-- prop). Empty for trial_lesson leads, which have no specific entity.
ALTER TABLE leads ADD COLUMN related_slug TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE leads DROP COLUMN related_slug;
