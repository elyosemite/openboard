import { useHealth } from '../hooks/useHealth';
import { SystemStatus } from './SystemStatus';

export function SystemStatusPanel() {
  const state = useHealth();
  return <SystemStatus state={state} />;
}
