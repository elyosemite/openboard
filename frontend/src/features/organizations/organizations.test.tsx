import { screen, waitFor, within } from '@testing-library/react';
import { problem, stubApi, userDto } from '../../test/fakeApi';
import { renderRoute } from '../../test/renderRoute';

const ACME = 'o1';
const organizationsUrl = (page = '1') => `GET /api/v1/organizations?page=${page}&page_size=20`;
const workspacesUrl = () => `GET /api/v1/workspaces?organization_id=${ACME}&page=1&page_size=20`;

const signedIn = { 'GET /api/v1/users/me': Response.json(userDto()) };

function pageOf(items: object[], page = 1, totalPages = 1) {
  return Response.json({
    items,
    page,
    page_size: 20,
    total_items: items.length,
    total_pages: totalPages,
  });
}

const acme = (role = 'owner') => ({
  id: ACME,
  name: 'Acme',
  created_at: '2026-11-17T12:00:00Z',
  current_user_role: role,
});

const workspaceDto = (overrides: object = {}) => ({
  id: 'w1',
  organization_id: ACME,
  name: 'Default',
  visibility: 'open',
  created_at: '2026-11-17T12:00:00Z',
  updated_at: '2026-11-17T12:00:00Z',
  ...overrides,
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('Organizations on the home page', () => {
  it('lists the Organizations with the role of the signed-in User', async () => {
    stubApi({ ...signedIn, [organizationsUrl()]: pageOf([acme('contributor')]) });
    renderRoute('/');

    const link = await screen.findByRole('link', { name: 'Acme' });
    expect(link).toHaveAttribute('href', `/organizations/${ACME}`);
    expect(screen.getByText('Organization Contributor')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Next' })).not.toBeInTheDocument();
  });

  it('names the Organization Owner role', async () => {
    stubApi({ ...signedIn, [organizationsUrl()]: pageOf([acme()]) });
    renderRoute('/');

    expect(await screen.findByText('Organization Owner')).toBeInTheDocument();
  });

  it('explains when there are no Organizations', async () => {
    stubApi({ ...signedIn, [organizationsUrl()]: pageOf([]) });
    renderRoute('/');

    expect(
      await screen.findByRole('heading', { name: 'No Organizations yet' }),
    ).toBeInTheDocument();
  });

  it('reports a failure to load', async () => {
    stubApi({ ...signedIn, [organizationsUrl()]: problem(500, 'internal-error', 'Unexpected') });
    renderRoute('/');

    expect(await screen.findByRole('alert')).toHaveTextContent('Unexpected');
  });

  it('moves between pages', async () => {
    stubApi({
      ...signedIn,
      [organizationsUrl()]: pageOf([acme()], 1, 2),
      [organizationsUrl('2')]: pageOf([{ ...acme(), id: 'o2', name: 'Globex' }], 2, 2),
    });
    const { user } = renderRoute('/');

    await screen.findByRole('link', { name: 'Acme' });
    expect(screen.getByText('Page 1 of 2')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Previous' })).toBeDisabled();

    await user.click(screen.getByRole('button', { name: 'Next' }));

    expect(await screen.findByRole('link', { name: 'Globex' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Next' })).toBeDisabled();

    await user.click(screen.getByRole('button', { name: 'Previous' }));

    expect(await screen.findByRole('link', { name: 'Acme' })).toBeInTheDocument();
  });
});

describe('Organization page', () => {
  const organizationRoutes = (role = 'owner', workspaces = [workspaceDto()]) => ({
    ...signedIn,
    [`GET /api/v1/organizations/${ACME}`]: Response.json(acme(role)),
    [workspacesUrl()]: pageOf(workspaces),
  });

  it('shows the Workspaces and their visibility', async () => {
    stubApi(
      organizationRoutes('contributor', [
        workspaceDto(),
        workspaceDto({ id: 'w2', name: 'Legal', visibility: 'restricted' }),
      ]),
    );
    renderRoute(`/organizations/${ACME}`);

    expect(await screen.findByRole('heading', { name: 'Acme', level: 1 })).toBeInTheDocument();
    const legal = await screen.findByRole('link', { name: 'Legal' });
    expect(legal).toHaveAttribute('href', '/workspaces/w2');
    expect(within(legal.closest('li') as HTMLElement).getByText('Restricted')).toBeInTheDocument();
    expect(screen.getByText('Open')).toBeInTheDocument();
  });

  it('offers creation only to those who manage the Organization', async () => {
    stubApi(organizationRoutes('contributor'));
    renderRoute(`/organizations/${ACME}`);

    await screen.findByRole('link', { name: 'Default' });
    expect(screen.queryByRole('button', { name: 'New Workspace' })).not.toBeInTheDocument();
  });

  it('explains when no Workspace can be entered', async () => {
    stubApi(organizationRoutes('contributor', []));
    renderRoute(`/organizations/${ACME}`);

    expect(
      await screen.findByRole('heading', { name: 'No Workspaces to show' }),
    ).toBeInTheDocument();
  });

  it('reports an Organization that cannot be seen', async () => {
    stubApi({
      ...signedIn,
      [`GET /api/v1/organizations/${ACME}`]: problem(404, 'not-found', 'Resource not found'),
    });
    renderRoute(`/organizations/${ACME}`);

    expect(await screen.findByRole('alert')).toHaveTextContent('Resource not found');
  });

  async function openForm(routes: Parameters<typeof stubApi>[0]) {
    const fetchMock = stubApi(routes);
    const rendered = renderRoute(`/organizations/${ACME}`);
    await rendered.user.click(await screen.findByRole('button', { name: 'New Workspace' }));
    return { fetchMock, ...rendered };
  }

  it('creates a Restricted Workspace by default and refreshes the list', async () => {
    let created = false;
    const { fetchMock, user } = await openForm({
      ...organizationRoutes(),
      [workspacesUrl()]: () =>
        pageOf(
          created ? [workspaceDto(), workspaceDto({ id: 'w2', name: 'Legal' })] : [workspaceDto()],
        ),
      'POST /api/v1/workspaces': () => {
        created = true;
        return Response.json(workspaceDto({ id: 'w2', name: 'Legal' }), { status: 201 });
      },
    });

    expect(screen.getByRole('button', { name: 'Restricted' })).toHaveAttribute(
      'aria-pressed',
      'true',
    );
    await user.type(screen.getByLabelText('Name'), 'Legal');
    await user.click(screen.getByRole('button', { name: 'Create Workspace' }));

    expect(await screen.findByRole('link', { name: 'Legal' })).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/workspaces',
      expect.objectContaining({
        body: JSON.stringify({ organization_id: ACME, name: 'Legal', visibility: 'restricted' }),
      }),
    );
    await waitFor(() => {
      expect(screen.queryByRole('form', { name: 'New Workspace' })).not.toBeInTheDocument();
    });
  });

  it('creates an Open Workspace when asked', async () => {
    const { fetchMock, user } = await openForm({
      ...organizationRoutes(),
      'POST /api/v1/workspaces': Response.json(workspaceDto(), { status: 201 }),
    });

    await user.click(screen.getByRole('button', { name: 'Open' }));
    expect(screen.getByText('Everyone in the Organization can enter.')).toBeInTheDocument();
    await user.type(screen.getByLabelText('Name'), 'Sales');
    await user.click(screen.getByRole('button', { name: 'Create Workspace' }));

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledWith(
        '/api/v1/workspaces',
        expect.objectContaining({
          body: JSON.stringify({ organization_id: ACME, name: 'Sales', visibility: 'open' }),
        }),
      );
    });
  });

  it.each([
    [
      'a name already in use',
      problem(409, 'workspace-name-taken', 'Workspace name already in use'),
      'Another Workspace in this Organization already uses this name.',
    ],
    [
      'an invalid name',
      problem(422, 'validation-error', 'Validation failed', {
        errors: [{ field: 'name', code: 'invalid', message: 'name must not be blank' }],
      }),
      'name must not be blank',
    ],
  ])('shows %s on the field', async (_case, reply, message) => {
    const { user } = await openForm({
      ...organizationRoutes(),
      'POST /api/v1/workspaces': reply,
    });

    await user.type(screen.getByLabelText('Name'), 'Legal');
    await user.click(screen.getByRole('button', { name: 'Create Workspace' }));

    await waitFor(() => {
      expect(screen.getByLabelText('Name')).toHaveAccessibleDescription(message);
    });
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it.each([
    ['a refused request', problem(403, 'forbidden', 'Not allowed')],
    ['an unreachable server', undefined],
  ])('reports %s above the form', async (_case, reply) => {
    const { user } = await openForm({
      ...organizationRoutes(),
      ...(reply && { 'POST /api/v1/workspaces': reply }),
    });

    await user.type(screen.getByLabelText('Name'), 'Legal');
    await user.click(screen.getByRole('button', { name: 'Create Workspace' }));

    expect(await screen.findByRole('alert')).toBeInTheDocument();
  });

  it('closes the form on cancel', async () => {
    const { user } = await openForm(organizationRoutes());

    await user.click(screen.getByRole('button', { name: 'Cancel' }));

    expect(screen.queryByRole('form', { name: 'New Workspace' })).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'New Workspace' })).toBeInTheDocument();
  });
});

describe('Workspace page', () => {
  it('shows the details of the Workspace', async () => {
    stubApi({
      ...signedIn,
      'GET /api/v1/workspaces/w2': Response.json(
        workspaceDto({ id: 'w2', name: 'Legal', visibility: 'restricted' }),
      ),
    });
    renderRoute('/workspaces/w2');

    expect(await screen.findByRole('heading', { name: 'Legal', level: 1 })).toBeInTheDocument();
    expect(screen.getByRole('region', { name: 'Details' })).toHaveTextContent('Restricted');
    expect(screen.getByRole('link', { name: 'Back to Organization' })).toHaveAttribute(
      'href',
      `/organizations/${ACME}`,
    );
  });

  it('reports a Workspace that cannot be entered', async () => {
    stubApi({
      ...signedIn,
      'GET /api/v1/workspaces/w2': problem(404, 'not-found', 'Resource not found'),
    });
    renderRoute('/workspaces/w2');

    expect(await screen.findByRole('alert')).toHaveTextContent('Resource not found');
  });
});
