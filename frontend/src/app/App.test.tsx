import { render, screen } from '@testing-library/react';
import { problem, stubApi } from '../test/fakeApi';
import { App } from './App';

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('App', () => {
  it('opens on sign-in when there is no session', async () => {
    stubApi({ 'GET /api/v1/users/me': problem(401, 'unauthenticated', 'Not signed in') });

    render(<App />);

    expect(await screen.findByRole('heading', { name: 'Sign in' })).toBeInTheDocument();
  });
});
