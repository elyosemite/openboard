import { Badge, Card, type BadgeTone } from '../../../shared/design-system';
import type { HealthState } from '../hooks/useHealth';

const presentation: Record<HealthState['kind'], { tone: BadgeTone; label: string }> = {
  loading: { tone: 'neutral', label: 'Checking API…' },
  online: { tone: 'success', label: 'API online' },
  offline: { tone: 'danger', label: 'API offline' },
};

interface SystemStatusProps {
  state: HealthState;
}

export function SystemStatus({ state }: SystemStatusProps) {
  const { tone, label } = presentation[state.kind];
  return (
    <Card title="System status">
      <Badge tone={tone}>{label}</Badge>
      {state.kind === 'offline' && <p>{state.message}</p>}
    </Card>
  );
}
