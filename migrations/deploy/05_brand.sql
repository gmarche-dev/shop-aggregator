CREATE TABLE IF NOT EXISTS "brand"
(
    brand_id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand_name  TEXT         NOT NULL UNIQUE,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_brand_brand_name ON "brand" (brand_name);

INSERT INTO brand (brand_id,brand_name) VALUES ('c2a2dea3-4fb0-4411-b395-bb1d14c92c0b','bulk');