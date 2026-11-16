import { ApiError, apiRequest, errorMessage } from './httpClient';

function stubFetch(response: Response) {
  const fetchMock = vi.fn().mockResolvedValue(response);
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

function problemResponse(problem: object, status: number) {
  return new Response(JSON.stringify(problem), {
    status,
    headers: { 'Content-Type': 'application/problem+json' },
  });
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('apiRequest', () => {
  it('GETs and returns the parsed body', async () => {
    const fetchMock = stubFetch(Response.json({ status: 'ok' }));

    await expect(apiRequest<{ status: string }>('/healthz')).resolves.toEqual({ status: 'ok' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/healthz',
      expect.objectContaining({ method: 'GET', credentials: 'same-origin', signal: undefined }),
    );
  });

  it('sends a body as JSON', async () => {
    const fetchMock = stubFetch(Response.json({ id: '1' }, { status: 201 }));

    await apiRequest('/api/v1/sessions', { method: 'POST', body: { username: 'gabriel' } });

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/sessions',
      expect.objectContaining({
        method: 'POST',
        body: '{"username":"gabriel"}',
        headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
      }),
    );
  });

  it('returns undefined for 204', async () => {
    stubFetch(new Response(null, { status: 204 }));
    await expect(apiRequest('/api/v1/sessions/current', { method: 'DELETE' })).resolves.toBe(
      undefined,
    );
  });

  it('turns problem responses into ApiError with the problem attached', async () => {
    const problem = {
      type: 'https://openboard.dev/problems/not-found',
      title: 'Resource not found',
      status: 404,
    };
    stubFetch(problemResponse(problem, 404));

    const error = await apiRequest('/missing').catch((e: unknown) => e);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({ status: 404, message: 'Resource not found', problem });
  });

  it('turns other failures into a generic ApiError', async () => {
    stubFetch(new Response(null, { status: 502 }));

    const error = await apiRequest('/healthz').catch((e: unknown) => e);

    expect(error).toMatchObject({
      name: 'ApiError',
      status: 502,
      message: 'Request failed with status 502',
      problem: undefined,
    });
  });
});

describe('errorMessage', () => {
  it('joins the problem title and detail', () => {
    const error = new ApiError(429, 'Too many failed sign-in attempts', {
      type: 'x',
      title: 'Too many failed sign-in attempts',
      status: 429,
      detail: 'Try again in a few minutes.',
    });
    expect(errorMessage(error)).toBe(
      'Too many failed sign-in attempts. Try again in a few minutes.',
    );
  });

  it('uses the title alone when there is no detail', () => {
    expect(errorMessage(new ApiError(401, 'Invalid username or password'))).toBe(
      'Invalid username or password.',
    );
  });

  it('explains network failures', () => {
    expect(errorMessage(new TypeError('Failed to fetch'))).toMatch(/Could not reach the server/);
  });
});
