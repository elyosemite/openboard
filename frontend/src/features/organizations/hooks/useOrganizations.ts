import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  createOrganization,
  createWorkspace,
  getOrganization,
  getWorkspace,
  listOrganizations,
  listWorkspaces,
  updateWorkspace,
  type WorkspaceChange,
} from '../api/organizationsApi';

const workspacesKey = (organizationId: string) => ['workspaces', organizationId] as const;

export function useOrganizations(page: number) {
  return useQuery({
    queryKey: ['organizations', page],
    queryFn: ({ signal }) => listOrganizations(page, signal),
  });
}

export function useCreateOrganization() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createOrganization,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['organizations'] }),
  });
}

export function useOrganization(id: string) {
  return useQuery({
    queryKey: ['organization', id],
    queryFn: ({ signal }) => getOrganization(id, signal),
  });
}

export function useWorkspaces(organizationId: string, page: number) {
  return useQuery({
    queryKey: [...workspacesKey(organizationId), page],
    queryFn: ({ signal }) => listWorkspaces(organizationId, page, signal),
  });
}

export function useWorkspace(id: string) {
  return useQuery({
    queryKey: ['workspace', id],
    queryFn: ({ signal }) => getWorkspace(id, signal),
  });
}

export function useUpdateWorkspace(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (change: WorkspaceChange) => updateWorkspace(id, change),
    onSuccess: (workspace) => {
      queryClient.setQueryData(['workspace', id], workspace);
      return queryClient.invalidateQueries({ queryKey: workspacesKey(workspace.organizationId) });
    },
  });
}

export function useCreateWorkspace(organizationId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createWorkspace,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: workspacesKey(organizationId) }),
  });
}
