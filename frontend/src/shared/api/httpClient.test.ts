import { ApiError, getJson } from './httpClient';

function stubFetch(response: Response) {
  const fetchMock = vi.fn().mockResolvedValue(response);
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('getJson', () => {
  it('returns the parsed body on success', async () => {
    const fetchMock = stubFetch(Response.json({ status: 'ok' }));

    await expect(getJson<{ status: string }>('/healthz')).resolves.toEqual({ status: 'ok' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/healthz',
      expect.objectContaining({ signal: undefined }),
    );
  });

  it('turns problem responses into ApiError with the problem attached', async () => {
    const problem = {
      type: 'https://openboard.dev/problems/not-found',
      title: 'Resource not found',
      status: 404,
    };
    stubFetch(
      new Response(JSON.stringify(problem), {
        status: 404,
        headers: { 'Content-Type': 'application/problem+json' },
      }),
    );

    const error = await getJson('/missing').catch((e: unknown) => e);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({ status: 404, message: 'Resource not found', problem });
  });

  it('turns other failures into a generic ApiError', async () => {
    stubFetch(new Response(null, { status: 502 }));

    const error = await getJson('/healthz').catch((e: unknown) => e);

    expect(error).toMatchObject({
      name: 'ApiError',
      status: 502,
      message: 'Request failed with status 502',
      problem: undefined,
    });
  });
});
