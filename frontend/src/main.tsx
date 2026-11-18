import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import 'ocaml-design-system/fonts.css';
import 'ocaml-design-system/tokens.css';
import 'ocaml-design-system/components.css';
import { App } from './app/App';
import './shared/design-system/styles/global.css';

const container = document.getElementById('root');
if (!container) {
  throw new Error('Root element #root not found');
}

createRoot(container).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
