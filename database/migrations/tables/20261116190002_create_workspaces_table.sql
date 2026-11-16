--liquibase formatted sql

--changeset openboard:20261116190002
--comment: A Workspace divides one Organization and is Open or Restricted (CONTEXT.md: Workspace).
CREATE TABLE workspaces (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID         NOT NULL,
    name            VARCHAR(255) NOT NULL,
    visibility      VARCHAR(16)  NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

--rollback DROP TABLE workspaces;
