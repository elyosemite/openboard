import { render, screen } from '@testing-library/react';
import { Badge, Card } from '.';

describe('Badge', () => {
  it('defaults to the neutral tone', () => {
    render(<Badge>Idle</Badge>);
    expect(screen.getByText('Idle')).toHaveClass('ob-badge', 'ob-badge--neutral');
  });

  it('applies the requested tone', () => {
    render(<Badge tone="danger">Down</Badge>);
    expect(screen.getByText('Down')).toHaveClass('ob-badge--danger');
  });
});

describe('Card', () => {
  it('renders an accessible region named by its title', () => {
    render(
      <Card title="System status">
        <p>content</p>
      </Card>,
    );
    expect(screen.getByRole('region', { name: 'System status' })).toHaveTextContent('content');
  });
});
