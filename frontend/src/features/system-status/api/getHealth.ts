import { getJson } from '../../../shared/api/httpClient';

export interface HealthStatus {
  status: string;
}

export function getHealth(signal?: AbortSignal): Promise<HealthStatus> {
  return getJson<HealthStatus>('/healthz', signal);
}
