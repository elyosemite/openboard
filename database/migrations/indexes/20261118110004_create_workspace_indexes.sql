--liquibase formatted sql

--changeset openboard:20261118110004
--comment: Workspace names are unique within an Organization ignoring case; the second index serves removing an Organization Membership.
CREATE UNIQUE INDEX ux_workspaces_organization_name ON workspaces (organization_id, lower(name));
CREATE INDEX ix_workspace_memberships_organization_user ON workspace_memberships (organization_id, user_id);

--rollback DROP INDEX ix_workspace_memberships_organization_user;
--rollback DROP INDEX ux_workspaces_organization_name;
