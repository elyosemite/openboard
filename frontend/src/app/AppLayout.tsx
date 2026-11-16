import { Outlet } from 'react-router';
import { CurrentUserMenu } from '../features/auth';
import './App.css';

export function AppLayout() {
  return (
    <div className="app-shell">
      <header className="app-header">
        <h1>OpenBoard</h1>
        <CurrentUserMenu />
      </header>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  );
}
