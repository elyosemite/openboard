import { render, screen } from '@testing-library/react';
import { App } from './App';

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('App', () => {
  it('renders the shell and reports the API as online', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(Response.json({ status: 'ok' })));

    render(<App />);

    expect(screen.getByRole('heading', { level: 1, name: 'OpenBoard' })).toBeInTheDocument();
    expect(await screen.findByText('API online')).toBeInTheDocument();
  });
});
