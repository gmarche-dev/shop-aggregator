CREATE TABLE IF NOT EXISTS "company"
(
    company_id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_name  TEXT         NOT NULL UNIQUE,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_company_company_name ON "company" (company_name);
