import { render, screen } from '@testing-library/react';
import { Alert, Button, Card, Icon, Lozenge, TextField } from '.';

describe('Lozenge', () => {
  it('defaults to the neutral tone', () => {
    render(<Lozenge>Idle</Lozenge>);
    expect(screen.getByText('Idle')).toHaveAttribute('class', 'lozenge');
  });

  it('applies the requested tone', () => {
    render(<Lozenge tone="danger">Down</Lozenge>);
    expect(screen.getByText('Down')).toHaveClass('lozenge', 'lozenge--removed');
  });
});

describe('Icon', () => {
  it('points at the named symbol and stays hidden from assistive technology', () => {
    const { container } = render(<Icon name="info" />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('aria-hidden', 'true');
    expect(svg?.querySelector('use')?.getAttribute('href')).toMatch(/#i-info$/);
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

describe('Button', () => {
  it('is a primary, non-submitting button by default', () => {
    render(<Button>Save</Button>);
    const button = screen.getByRole('button', { name: 'Save' });
    expect(button).toHaveAttribute('type', 'button');
    expect(button).toHaveClass('btn', 'btn--primary');
    expect(button).toBeEnabled();
    expect(button).not.toHaveAttribute('aria-busy');
  });

  it('blocks clicks while busy', () => {
    render(
      <Button type="submit" variant="secondary" busy>
        Saving
      </Button>,
    );
    const button = screen.getByRole('button', { name: 'Saving' });
    expect(button).toBeDisabled();
    expect(button).toHaveAttribute('aria-busy', 'true');
    expect(button).toHaveClass('btn', 'btn--outline');
  });

  it('honours disabled', () => {
    render(<Button disabled>Off</Button>);
    expect(screen.getByRole('button', { name: 'Off' })).toBeDisabled();
  });
});

describe('TextField', () => {
  it('labels the input', () => {
    render(<TextField label="Username" />);
    const input = screen.getByLabelText('Username');
    expect(input).not.toHaveAttribute('aria-invalid');
    expect(input).not.toHaveAttribute('aria-describedby');
    expect(input.parentElement).not.toHaveClass('is-invalid');
  });

  it('links the error message and marks the input invalid', () => {
    render(<TextField label="Password" error="Too short" />);
    const input = screen.getByLabelText('Password');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAccessibleDescription('Too short');
    expect(input.parentElement).toHaveClass('input', 'is-invalid');
  });
});

describe('Alert', () => {
  it('announces info politely by default', () => {
    render(<Alert>Heads up</Alert>);
    expect(screen.getByRole('status')).toHaveTextContent('Heads up');
  });

  it('interrupts for danger', () => {
    render(<Alert tone="danger">Failed</Alert>);
    const alert = screen.getByRole('alert');
    expect(alert).toHaveTextContent('Failed');
    expect(alert.parentElement).toHaveClass('section-message', 'section-message--error');
  });
});
