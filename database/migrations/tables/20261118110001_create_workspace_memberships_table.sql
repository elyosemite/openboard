--liquibase formatted sql

--changeset openboard:20261118110001
--comment: A Workspace Membership links one User to one Workspace with one role (CONTEXT.md: Workspace Membership).
CREATE TABLE workspace_memberships (
    workspace_id    UUID        NOT NULL,
    organization_id UUID        NOT NULL,
    user_id         UUID        NOT NULL,
    role            VARCHAR(16) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (workspace_id, user_id)
);

--rollback DROP TABLE workspace_memberships;
