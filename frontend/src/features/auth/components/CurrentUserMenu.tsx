import { Button } from '../../../shared/design-system';
import { useCurrentUser, useSignOut } from '../hooks/useAuth';
import './AuthPage.css';

/** The signed-in User's name and the sign-out action, for the app header. */
export function CurrentUserMenu() {
  const { data: user } = useCurrentUser();
  const signOut = useSignOut();
  if (!user) {
    return null;
  }
  return (
    <div className="current-user">
      <span className="current-user__name">{user.name}</span>
      <Button
        variant="secondary"
        busy={signOut.isPending}
        onClick={() => {
          signOut.mutate();
        }}
      >
        Sign out
      </Button>
    </div>
  );
}
