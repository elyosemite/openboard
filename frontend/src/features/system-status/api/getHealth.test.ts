import { getHealth } from './getHealth';

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('getHealth', () => {
  it('queries the liveness probe', async () => {
    const fetchMock = vi.fn().mockResolvedValue(Response.json({ status: 'ok' }));
    vi.stubGlobal('fetch', fetchMock);
    const controller = new AbortController();

    await expect(getHealth(controller.signal)).resolves.toEqual({ status: 'ok' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/healthz',
      expect.objectContaining({ signal: controller.signal }),
    );
  });
});
