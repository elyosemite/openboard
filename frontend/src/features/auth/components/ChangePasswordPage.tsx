import { useState, type SubmitEvent } from 'react';
import { useNavigate } from 'react-router';
import { ApiError, errorMessage } from '../../../shared/api/httpClient';
import { Alert, Button, Card, TextField } from '../../../shared/design-system';
import { useChangePassword, useCurrentUser } from '../hooks/useAuth';
import './AuthPage.css';

type FieldMessages = Partial<Record<'current_password' | 'new_password', string>>;

/** Field errors from a validation-error problem, or null for any other failure. */
function fieldMessagesOf(error: unknown): FieldMessages | null {
  if (!(error instanceof ApiError) || error.status !== 422) {
    return null;
  }
  const messages: FieldMessages = {};
  for (const { field, message } of error.problem?.errors ?? []) {
    if (field === 'current_password' || field === 'new_password') {
      messages[field] = message;
    }
  }
  return messages;
}

export function ChangePasswordPage() {
  const { data: user } = useCurrentUser();
  const change = useChangePassword();
  const navigate = useNavigate();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [mismatch, setMismatch] = useState(false);

  const fieldMessages = change.isError ? fieldMessagesOf(change.error) : null;

  function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    const confirmed = newPassword === confirmation;
    setMismatch(!confirmed);
    if (!confirmed) {
      return;
    }
    change.mutate(
      { currentPassword, newPassword },
      {
        onSuccess: () => {
          void navigate('/', { replace: true });
        },
      },
    );
  }

  return (
    <div className="auth-page">
      <Card title="Change password">
        <form className="auth-form" onSubmit={handleSubmit}>
          {user?.mustChangePassword && (
            <Alert>
              You signed in with a Temporary Password. Choose a new password to continue.
            </Alert>
          )}
          {change.isError && !fieldMessages && (
            <Alert tone="danger">{errorMessage(change.error)}</Alert>
          )}
          <TextField
            label="Current password"
            name="current-password"
            type="password"
            autoComplete="current-password"
            required
            value={currentPassword}
            error={fieldMessages?.current_password}
            onChange={(e) => {
              setCurrentPassword(e.target.value);
            }}
          />
          <TextField
            label="New password"
            name="new-password"
            type="password"
            autoComplete="new-password"
            required
            value={newPassword}
            error={fieldMessages?.new_password}
            onChange={(e) => {
              setNewPassword(e.target.value);
            }}
          />
          <p className="auth-form__hint">Use 12 to 128 characters. A long passphrase works well.</p>
          <TextField
            label="Confirm new password"
            name="confirm-password"
            type="password"
            autoComplete="new-password"
            required
            value={confirmation}
            error={mismatch ? 'The passwords do not match.' : undefined}
            onChange={(e) => {
              setConfirmation(e.target.value);
            }}
          />
          <Button type="submit" busy={change.isPending}>
            Change password
          </Button>
        </form>
      </Card>
    </div>
  );
}
