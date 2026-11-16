export interface FieldError {
  field: string;
  code: string;
  message: string;
}

export interface ProblemDetails {
  type: string;
  title: string;
  status: number;
  detail?: string;
  instance?: string;
  errors?: FieldError[];
}

export class ApiError extends Error {
  readonly status: number;
  readonly problem: ProblemDetails | undefined;

  constructor(status: number, message: string, problem?: ProblemDetails) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.problem = problem;
  }
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  /** Sent as JSON. */
  body?: unknown;
  signal?: AbortSignal;
}

const PROBLEM_CONTENT_TYPE = 'application/problem+json';

/**
 * Calls the API and returns the JSON body, or undefined for 204 responses.
 * Error responses become an ApiError.
 */
export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, signal } = options;
  const headers: Record<string, string> = { Accept: 'application/json' };
  const init: RequestInit = { method, headers, signal, credentials: 'same-origin' };
  if (body !== undefined) {
    headers['Content-Type'] = 'application/json';
    init.body = JSON.stringify(body);
  }

  const response = await fetch(path, init);
  if (!response.ok) {
    throw await toApiError(response);
  }
  if (response.status === 204) {
    return undefined as T;
  }
  return (await response.json()) as T;
}

async function toApiError(response: Response): Promise<ApiError> {
  const contentType = response.headers.get('Content-Type') ?? '';
  if (contentType.includes(PROBLEM_CONTENT_TYPE)) {
    const problem = (await response.json()) as ProblemDetails;
    return new ApiError(response.status, problem.title, problem);
  }
  return new ApiError(response.status, `Request failed with status ${String(response.status)}`);
}

/** A sentence to show the User for any failure of apiRequest. */
export function errorMessage(error: unknown): string {
  if (!(error instanceof ApiError)) {
    return 'Could not reach the server. Check your connection and try again.';
  }
  const detail = error.problem?.detail;
  return detail ? `${error.message}. ${detail}` : `${error.message}.`;
}
