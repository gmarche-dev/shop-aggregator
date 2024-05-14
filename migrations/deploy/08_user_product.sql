CREATE TABLE IF NOT EXISTS "user_product"
(
    user_product_id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id       UUID         NOT NULL,
    user_id          UUID         NOT NULL,
    bill_id          UUID         NOT NULL,
    price            TEXT         NOT NULL,
    product_type     TEXT         NOT NULL DEFAULT 'barcoded_product',
    product_size     TEXT         NOT NULL DEFAULT '',
    size_format      TEXT         NOT NULL DEFAULT '',
    quantity         INTEGER      NOT NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_product_product_id ON "user_product" (product_id);
CREATE INDEX IF NOT EXISTS idx_user_product_bill_id ON "user_product" (bill_id);
CREATE INDEX IF NOT EXISTS idx_user_product_user_id ON "user_product" (user_id);
