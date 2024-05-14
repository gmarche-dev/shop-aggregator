CREATE TABLE IF NOT EXISTS "product"
(
    product_id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ean           TEXT         NOT NULL UNIQUE,
    product_name  TEXT         NOT NULL,
    brand_id      UUID         NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_ean ON "product" (ean);
CREATE INDEX IF NOT EXISTS idx_product_product_name ON "product" (product_name);
CREATE INDEX IF NOT EXISTS idx_product_brand_id ON "product" (brand_id);

INSERT into product (product_id,ean,product_name,brand_id) VALUES
('2e30955b-0f88-43df-8924-1ec21afed0aa','meat','meat', 'c2a2dea3-4fb0-4411-b395-bb1d14c92c0b'),
('eb5be0d0-b3f6-4f2f-b582-9a7dd566b549','vegetable','vegetable', 'c2a2dea3-4fb0-4411-b395-bb1d14c92c0b'),
('3c94d3d7-7bce-40d6-8f11-c8fdeb328d41','fruit','fruit', 'c2a2dea3-4fb0-4411-b395-bb1d14c92c0b');