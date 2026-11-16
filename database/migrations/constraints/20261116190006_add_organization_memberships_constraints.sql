--liquibase formatted sql

--changeset openboard:20261116190006
--comment: Memberships go with their Organization; Users are never deleted, so they block removal. Role is Organization Owner or Organization Contributor.
ALTER TABLE organization_memberships
    ADD CONSTRAINT fk_organization_memberships_organization FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_organization_memberships_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    ADD CONSTRAINT ck_organization_memberships_role CHECK (role IN ('owner', 'contributor'));

--rollback ALTER TABLE organization_memberships DROP CONSTRAINT ck_organization_memberships_role, DROP CONSTRAINT fk_organization_memberships_user, DROP CONSTRAINT fk_organization_memberships_organization;
