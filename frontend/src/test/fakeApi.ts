type Reply = Response | ((init: RequestInit) => Response);

/**
 * Replaces fetch with a fake API keyed by "METHOD path". Unknown routes fail
 * like a network error so tests notice unexpected calls.
 */
export function stubApi(routes: Record<string, Reply>) {
  const fetchMock = vi.fn((input: string, init: RequestInit = {}) => {
    const key = `${init.method ?? 'GET'} ${input}`;
    const reply = routes[key];
    if (reply === undefined) {
      return Promise.reject(new TypeError(`unexpected request ${key}`));
    }
    return Promise.resolve(typeof reply === 'function' ? reply(init) : reply.clone());
  });
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

export function problem(status: number, slug: string, title: string, extra: object = {}) {
  return new Response(
    JSON.stringify({ type: `https://openboard.dev/problems/${slug}`, title, status, ...extra }),
    { status, headers: { 'Content-Type': 'application/problem+json' } },
  );
}

export function userDto(overrides: object = {}) {
  return {
    id: 'u1',
    username: 'admin',
    name: 'Administrator',
    is_administrator: true,
    must_change_password: false,
    ...overrides,
  };
}

export const unauthenticated = () => problem(401, 'unauthenticated', 'Not signed in');
export const noContent = () => new Response(null, { status: 204 });
