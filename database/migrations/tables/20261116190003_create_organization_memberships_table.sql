--liquibase formatted sql

--changeset openboard:20261116190003
--comment: An Organization Membership links one User to one Organization with one role.
CREATE TABLE organization_memberships (
    organization_id UUID        NOT NULL,
    user_id         UUID        NOT NULL,
    role            VARCHAR(16) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (organization_id, user_id)
);

--rollback DROP TABLE organization_memberships;
