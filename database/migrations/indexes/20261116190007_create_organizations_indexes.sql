--liquibase formatted sql

--changeset openboard:20261116190007
--comment: Listing an Organization's Workspaces and a User's Organizations.
CREATE INDEX ix_workspaces_organization_id ON workspaces (organization_id);
CREATE INDEX ix_organization_memberships_user_id ON organization_memberships (user_id);

--rollback DROP INDEX ix_organization_memberships_user_id;
--rollback DROP INDEX ix_workspaces_organization_id;
