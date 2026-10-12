-- +goose Up
-- Which social network a "found us on social media" lead came from — blank
-- for every lead with a different source. 'instagram' backs the frontend's
-- disclaimer-safe "Запрещённая соц. сеть" label (Meta is recognized as
-- extremist and banned in Russia) rather than naming the platform directly.
ALTER TABLE leads ADD COLUMN social_network TEXT NOT NULL DEFAULT '';
ALTER TABLE leads ADD CONSTRAINT leads_social_network_check CHECK (social_network IN ('', 'telegram', 'vk', 'pinterest', 'instagram', 'youtube'));

-- +goose Down
ALTER TABLE leads DROP CONSTRAINT leads_social_network_check;
ALTER TABLE leads DROP COLUMN social_network;
