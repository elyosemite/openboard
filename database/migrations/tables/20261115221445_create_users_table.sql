--liquibase formatted sql

--changeset openboard:20261115221445
--comment: Users are global to the Instance (CONTEXT.md: User, Username, Temporary Password, Administrator).
CREATE TABLE users (
    id                   UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username             VARCHAR(32)  NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    password_hash        TEXT         NOT NULL,
    must_change_password BOOLEAN      NOT NULL DEFAULT FALSE,
    is_administrator     BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

--rollback DROP TABLE users;
