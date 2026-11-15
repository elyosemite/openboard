import { SystemStatusPanel } from '../features/system-status';
import './App.css';

export function App() {
  return (
    <div className="app-shell">
      <header className="app-header">
        <h1>OpenBoard</h1>
      </header>
      <main className="app-main">
        <SystemStatusPanel />
      </main>
    </div>
  );
}
