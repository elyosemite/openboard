--liquibase formatted sql

--changeset openboard:20261115221446
--comment: Server-side sessions (ADR 0006). id is the SHA-256 of the cookie token, never the token itself.
CREATE TABLE sessions (
    id           BYTEA       PRIMARY KEY,
    user_id      UUID        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL
);

--rollback DROP TABLE sessions;
