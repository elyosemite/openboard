import { useEffect, useState } from 'react';
import { getHealth, type HealthStatus } from '../api/getHealth';

export type HealthState =
  { kind: 'loading' } | { kind: 'online' } | { kind: 'offline'; message: string };

type FetchHealth = (signal: AbortSignal) => Promise<HealthStatus>;

function messageOf(error: unknown): string {
  return error instanceof Error ? error.message : 'Unknown error';
}

export function useHealth(fetchHealth: FetchHealth = getHealth): HealthState {
  const [state, setState] = useState<HealthState>({ kind: 'loading' });

  useEffect(() => {
    const controller = new AbortController();
    fetchHealth(controller.signal).then(
      () => {
        setState({ kind: 'online' });
      },
      (error: unknown) => {
        if (!controller.signal.aborted) {
          setState({ kind: 'offline', message: messageOf(error) });
        }
      },
    );
    return () => {
      controller.abort();
    };
  }, [fetchHealth]);

  return state;
}
