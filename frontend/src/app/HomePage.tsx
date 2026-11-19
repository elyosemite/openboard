import { useCurrentUser } from '../features/auth';
import { OrganizationList } from '../features/organizations';
import { Card } from '../shared/design-system';
import './App.css';

export function HomePage() {
  const { data: user } = useCurrentUser();
  return (
    <div className="home-page">
      <Card title={`Welcome, ${user?.name ?? ''}`}>
        <p>You are signed in as {user?.username}.</p>
      </Card>
      <OrganizationList canCreate={user?.isAdministrator ?? false} />
    </div>
  );
}
