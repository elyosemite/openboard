--liquibase formatted sql

--changeset openboard:20261116190001
--comment: The Organization is the tenant (ADR 0003); every domain record belongs to one.
CREATE TABLE organizations (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

--rollback DROP TABLE organizations;
