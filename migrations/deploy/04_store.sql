CREATE TABLE IF NOT EXISTS "store"
(
    store_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    address    TEXT         NOT NULL,
    zip_code   TEXT         NOT NULL,
    city       TEXT         NOT NULL,
    country    TEXT         NOT NULL,
    store_name TEXT         NOT NULL,
    store_type TEXT         NOT NULL,
    url        TEXT         NOT NULL,
    company_id UUID         NOT NULL,
    is_active  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_store_name ON "store" (store_name);
CREATE INDEX IF NOT EXISTS idx_store_company_id ON "store" (company_id);
CREATE INDEX IF NOT EXISTS idx_store_zip_code ON "store" (zip_code);
