export const ORGANIZATION_ROUTE = '/organizations/:organizationId';
export const WORKSPACE_ROUTE = '/workspaces/:workspaceId';

export const organizationPath = (id: string) => `/organizations/${encodeURIComponent(id)}`;
export const workspacePath = (id: string) => `/workspaces/${encodeURIComponent(id)}`;
