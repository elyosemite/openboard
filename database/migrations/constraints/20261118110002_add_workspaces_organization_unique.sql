--liquibase formatted sql

--changeset openboard:20261118110002
--comment: Lets a Workspace Membership reference its Workspace together with the Workspace's Organization.
ALTER TABLE workspaces
    ADD CONSTRAINT uq_workspaces_organization_id UNIQUE (organization_id, id);

--rollback ALTER TABLE workspaces DROP CONSTRAINT uq_workspaces_organization_id;
