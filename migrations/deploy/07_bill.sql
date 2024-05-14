CREATE TABLE IF NOT EXISTS "bill"
(
    bill_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id       UUID         NOT NULL,
    store_id      UUID         NOT NULL,
    amount        TEXT         NOT NULL,
    bill_state    TEXT         NOT NULL DEFAULT 'create',
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bill_user_id ON "bill" (user_id);
CREATE INDEX IF NOT EXISTS idx_bill_store_id ON "bill" (store_id);
CREATE INDEX IF NOT EXISTS idx_bill_bill_state ON "bill" (bill_state);
