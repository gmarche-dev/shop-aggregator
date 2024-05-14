CREATE TABLE IF NOT EXISTS "auth"
(
    user_id    UUID PRIMARY KEY,
    token      TEXT         NOT NULL,
    is_active  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_auth_token ON "auth" (token);
