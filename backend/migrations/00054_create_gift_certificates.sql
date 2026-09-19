-- +goose Up
-- Выданные подарочные сертификаты. Номер (#DDMMYYNN) генерируется в
-- GiftCertificateService.Create; уникальность обеспечена индексом, без
-- блокировок — админ один (см. обоснование у ClientComment в model.go).
CREATE TABLE gift_certificates (
    id         BIGSERIAL PRIMARY KEY,
    number     TEXT NOT NULL,
    issued_at  DATE NOT NULL DEFAULT CURRENT_DATE,
    kind       TEXT NOT NULL CHECK (kind IN ('amount', 'course', 'masterclass', 'any_masterclass')),
    value      TEXT NOT NULL DEFAULT '', -- пусто для any_masterclass
    recipient  TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX gift_certificates_number_idx ON gift_certificates (number);

-- +goose Down
DROP TABLE gift_certificates;
