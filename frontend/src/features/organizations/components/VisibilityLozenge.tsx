import { Lozenge } from '../../../shared/design-system';
import type { WorkspaceVisibility } from '../api/organizationsApi';
import { visibilityLabels } from './visibility';

export function VisibilityLozenge({ visibility }: { visibility: WorkspaceVisibility }) {
  return (
    <Lozenge tone={visibility === 'open' ? 'success' : 'neutral'}>
      {visibilityLabels[visibility]}
    </Lozenge>
  );
}
