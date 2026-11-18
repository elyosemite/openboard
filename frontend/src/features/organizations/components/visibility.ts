import type { WorkspaceVisibility } from '../api/organizationsApi';

export const visibilityLabels: Record<WorkspaceVisibility, string> = {
  open: 'Open',
  restricted: 'Restricted',
};
