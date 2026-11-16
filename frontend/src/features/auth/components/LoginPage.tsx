import { useState, type SubmitEvent } from 'react';
import { Navigate, useLocation } from 'react-router';
import { errorMessage } from '../../../shared/api/httpClient';
import { Alert, Button, Card, TextField } from '../../../shared/design-system';
import { useCurrentUser, useSignIn } from '../hooks/useAuth';
import { CHANGE_PASSWORD_PATH, type SignInState } from '../paths';
import './AuthPage.css';

export function LoginPage() {
  const { data: user } = useCurrentUser();
  const signIn = useSignIn();
  const location = useLocation();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  if (user) {
    const from = (location.state as SignInState | null)?.from ?? '/';
    return <Navigate to={user.mustChangePassword ? CHANGE_PASSWORD_PATH : from} replace />;
  }

  function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    signIn.mutate({ username, password });
  }

  return (
    <main className="auth-page">
      <h1 className="auth-page__brand">OpenBoard</h1>
      <Card title="Sign in">
        <form className="auth-form" onSubmit={handleSubmit}>
          {signIn.isError && <Alert tone="danger">{errorMessage(signIn.error)}</Alert>}
          <TextField
            label="Username"
            name="username"
            autoComplete="username"
            autoCapitalize="none"
            spellCheck={false}
            required
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
            }}
          />
          <TextField
            label="Password"
            name="password"
            type="password"
            autoComplete="current-password"
            required
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
            }}
          />
          <Button type="submit" busy={signIn.isPending}>
            Sign in
          </Button>
        </form>
      </Card>
    </main>
  );
}
