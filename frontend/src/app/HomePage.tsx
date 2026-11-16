import { useCurrentUser } from '../features/auth';
import { Card } from '../shared/design-system';

export function HomePage() {
  const { data: user } = useCurrentUser();
  return (
    <Card title={`Welcome, ${user?.name ?? ''}`}>
      <p>You are signed in as {user?.username}.</p>
    </Card>
  );
}
