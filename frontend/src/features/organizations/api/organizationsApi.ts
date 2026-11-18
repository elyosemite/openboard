import { apiRequest } from '../../../shared/api/httpClient';

export const PAGE_SIZE = 20;

export type OrganizationRole = 'owner' | 'contributor';
export type WorkspaceVisibility = 'open' | 'restricted';

export interface Organization {
  id: string;
  name: string;
  currentUserRole: OrganizationRole;
}

export interface Workspace {
  id: string;
  organizationId: string;
  name: string;
  visibility: WorkspaceVisibility;
  createdAt: string;
}

export interface Page<T> {
  items: T[];
  page: number;
  totalPages: number;
}

export interface NewWorkspace {
  organizationId: string;
  name: string;
  visibility: WorkspaceVisibility;
}

interface PageDto<T> {
  items: T[];
  page: number;
  total_pages: number;
}

interface OrganizationDto {
  id: string;
  name: string;
  current_user_role: OrganizationRole;
}

interface WorkspaceDto {
  id: string;
  organization_id: string;
  name: string;
  visibility: WorkspaceVisibility;
  created_at: string;
}

function toOrganization(dto: OrganizationDto): Organization {
  return { id: dto.id, name: dto.name, currentUserRole: dto.current_user_role };
}

function toWorkspace(dto: WorkspaceDto): Workspace {
  return {
    id: dto.id,
    organizationId: dto.organization_id,
    name: dto.name,
    visibility: dto.visibility,
    createdAt: dto.created_at,
  };
}

function toPage<D, T>(dto: PageDto<D>, toItem: (item: D) => T): Page<T> {
  return { items: dto.items.map(toItem), page: dto.page, totalPages: dto.total_pages };
}

export async function listOrganizations(
  page: number,
  signal?: AbortSignal,
): Promise<Page<Organization>> {
  const query = new URLSearchParams({ page: String(page), page_size: String(PAGE_SIZE) });
  const dto = await apiRequest<PageDto<OrganizationDto>>(`/api/v1/organizations?${query}`, {
    signal,
  });
  return toPage(dto, toOrganization);
}

export async function getOrganization(id: string, signal?: AbortSignal): Promise<Organization> {
  const path = `/api/v1/organizations/${encodeURIComponent(id)}`;
  return toOrganization(await apiRequest<OrganizationDto>(path, { signal }));
}

export async function listWorkspaces(
  organizationId: string,
  page: number,
  signal?: AbortSignal,
): Promise<Page<Workspace>> {
  const query = new URLSearchParams({
    organization_id: organizationId,
    page: String(page),
    page_size: String(PAGE_SIZE),
  });
  const dto = await apiRequest<PageDto<WorkspaceDto>>(`/api/v1/workspaces?${query}`, { signal });
  return toPage(dto, toWorkspace);
}

export async function getWorkspace(id: string, signal?: AbortSignal): Promise<Workspace> {
  const path = `/api/v1/workspaces/${encodeURIComponent(id)}`;
  return toWorkspace(await apiRequest<WorkspaceDto>(path, { signal }));
}

export async function createWorkspace(workspace: NewWorkspace): Promise<Workspace> {
  const dto = await apiRequest<WorkspaceDto>('/api/v1/workspaces', {
    method: 'POST',
    body: {
      organization_id: workspace.organizationId,
      name: workspace.name,
      visibility: workspace.visibility,
    },
  });
  return toWorkspace(dto);
}
