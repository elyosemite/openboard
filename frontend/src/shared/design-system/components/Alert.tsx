import type { ReactNode } from 'react';
import './Alert.css';

export type AlertTone = 'info' | 'danger';

interface AlertProps {
  tone?: AlertTone;
  children: ReactNode;
}

/** Danger alerts interrupt screen readers; info alerts are announced politely. */
export function Alert({ tone = 'info', children }: AlertProps) {
  return (
    <div role={tone === 'danger' ? 'alert' : 'status'} className={`ob-alert ob-alert--${tone}`}>
      {children}
    </div>
  );
}
