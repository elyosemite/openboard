--liquibase formatted sql

--changeset openboard:20261116190005
--comment: A Workspace belongs to exactly one Organization, has a name and is Open or Restricted.
ALTER TABLE workspaces
    ADD CONSTRAINT fk_workspaces_organization FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    ADD CONSTRAINT ck_workspaces_name_not_blank CHECK (btrim(name) <> ''),
    ADD CONSTRAINT ck_workspaces_visibility CHECK (visibility IN ('open', 'restricted'));

--rollback ALTER TABLE workspaces DROP CONSTRAINT ck_workspaces_visibility, DROP CONSTRAINT ck_workspaces_name_not_blank, DROP CONSTRAINT fk_workspaces_organization;
