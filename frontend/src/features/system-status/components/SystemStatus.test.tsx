import { render, screen } from '@testing-library/react';
import { SystemStatus } from './SystemStatus';

describe('SystemStatus', () => {
  it('shows a neutral badge while loading', () => {
    render(<SystemStatus state={{ kind: 'loading' }} />);
    expect(screen.getByText('Checking API…')).toHaveClass('ob-badge--neutral');
  });

  it('shows a success badge when online', () => {
    render(<SystemStatus state={{ kind: 'online' }} />);
    expect(screen.getByText('API online')).toHaveClass('ob-badge--success');
  });

  it('shows a danger badge and the reason when offline', () => {
    render(<SystemStatus state={{ kind: 'offline', message: 'Service unavailable' }} />);
    expect(screen.getByText('API offline')).toHaveClass('ob-badge--danger');
    expect(screen.getByText('Service unavailable')).toBeInTheDocument();
  });
});
