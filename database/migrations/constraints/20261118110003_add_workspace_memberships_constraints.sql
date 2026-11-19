--liquibase formatted sql

--changeset openboard:20261118110003
--comment: A Workspace Membership lasts only while its Workspace and the User's Organization Membership in the same Organization exist. Role is Workspace Owner, Contributor or Viewer.
ALTER TABLE workspace_memberships
    ADD CONSTRAINT fk_workspace_memberships_workspace FOREIGN KEY (organization_id, workspace_id) REFERENCES workspaces (organization_id, id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_workspace_memberships_organization_membership FOREIGN KEY (organization_id, user_id) REFERENCES organization_memberships (organization_id, user_id) ON DELETE CASCADE,
    ADD CONSTRAINT ck_workspace_memberships_role CHECK (role IN ('owner', 'contributor', 'viewer'));

--rollback ALTER TABLE workspace_memberships DROP CONSTRAINT ck_workspace_memberships_role, DROP CONSTRAINT fk_workspace_memberships_organization_membership, DROP CONSTRAINT fk_workspace_memberships_workspace;
