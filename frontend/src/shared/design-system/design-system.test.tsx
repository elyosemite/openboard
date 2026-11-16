import { render, screen } from '@testing-library/react';
import { Alert, Badge, Button, Card, TextField } from '.';

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

describe('Button', () => {
  it('is a primary, non-submitting button by default', () => {
    render(<Button>Save</Button>);
    const button = screen.getByRole('button', { name: 'Save' });
    expect(button).toHaveAttribute('type', 'button');
    expect(button).toHaveClass('ob-button--primary');
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
    expect(button).toHaveClass('ob-button--secondary');
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
  });

  it('links the error message and marks the input invalid', () => {
    render(<TextField label="Password" error="Too short" />);
    const input = screen.getByLabelText('Password');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAccessibleDescription('Too short');
  });
});

describe('Alert', () => {
  it('announces info politely by default', () => {
    render(<Alert>Heads up</Alert>);
    expect(screen.getByRole('status')).toHaveTextContent('Heads up');
  });

  it('interrupts for danger', () => {
    render(<Alert tone="danger">Failed</Alert>);
    expect(screen.getByRole('alert')).toHaveClass('ob-alert--danger');
  });
});
