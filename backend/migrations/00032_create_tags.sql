-- +goose Up
-- Two independent tag types get two independent tables (not one shared
-- `tags` table with a `type` discriminator) so a query can never mix a
-- product tag into a client-type slot by mistake, and each type gets its
-- own case-insensitive uniqueness constraint for free.
CREATE TABLE product_tags (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_product_tags_name_lower ON product_tags (lower(name));

CREATE TABLE client_type_tags (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_client_type_tags_name_lower ON client_type_tags (lower(name));

CREATE TABLE client_product_tags (
    client_id      BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    product_tag_id BIGINT NOT NULL REFERENCES product_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (client_id, product_tag_id)
);

CREATE TABLE client_client_type_tags (
    client_id          BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    client_type_tag_id BIGINT NOT NULL REFERENCES client_type_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (client_id, client_type_tag_id)
);

-- +goose Down
DROP TABLE client_client_type_tags;
DROP TABLE client_product_tags;
DROP TABLE client_type_tags;
DROP TABLE product_tags;
