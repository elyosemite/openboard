import { Navigate, Outlet, useLocation } from 'react-router';
import { errorMessage } from '../../../shared/api/httpClient';
import { Alert } from '../../../shared/design-system';
import { useCurrentUser } from '../hooks/useAuth';
import { CHANGE_PASSWORD_PATH, SIGN_IN_PATH, type SignInState } from '../paths';

/**
 * Renders the child routes only for a signed-in User. Signed-out visitors go
 * to sign-in; a pending Temporary Password sends the User to change it first.
 */
export function RequireAuth() {
  const { data: user, isPending, isError, error } = useCurrentUser();
  const location = useLocation();

  if (isPending) {
    return <p role="status">Loading…</p>;
  }
  if (isError) {
    return <Alert tone="danger">{errorMessage(error)}</Alert>;
  }
  if (!user) {
    const state: SignInState = { from: location.pathname };
    return <Navigate to={SIGN_IN_PATH} replace state={state} />;
  }
  if (user.mustChangePassword && location.pathname !== CHANGE_PASSWORD_PATH) {
    return <Navigate to={CHANGE_PASSWORD_PATH} replace />;
  }
  return <Outlet />;
}
