import { Link, NavLink, Outlet } from 'react-router';
import { CurrentUserMenu } from '../features/auth';
import { Icon } from '../shared/design-system';
import './App.css';

export function AppLayout() {
  return (
    <div className="workspace">
      <aside className="sidebar workspace__sidebar">
        <Link to="/" className="wordmark wordmark--sm">
          <span className="wordmark__glyph wordmark__glyph--neutral">
            <Icon name="grid" />
          </span>
          OpenBoard
        </Link>
        <nav className="nav-group" aria-label="Main">
          <NavLink to="/" end className="nav-item">
            <Icon name="grid" />
            <span className="nav-item__text">Home</span>
          </NavLink>
        </nav>
      </aside>
      <main className="workspace__main">
        <header className="workspace__topbar">
          <span className="workspace__crumb">OpenBoard</span>
          <span />
          <div className="workspace__topbar-end">
            <CurrentUserMenu />
          </div>
        </header>
        <div className="app-content">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
