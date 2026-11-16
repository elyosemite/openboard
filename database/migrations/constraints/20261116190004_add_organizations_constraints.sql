--liquibase formatted sql

--changeset openboard:20261116190004
--comment: An Organization name is never blank.
ALTER TABLE organizations
    ADD CONSTRAINT ck_organizations_name_not_blank CHECK (btrim(name) <> '');

--rollback ALTER TABLE organizations DROP CONSTRAINT ck_organizations_name_not_blank;
