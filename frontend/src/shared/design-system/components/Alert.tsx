import type { ReactNode } from 'react';
import { Icon } from './Icon';

export type AlertTone = 'info' | 'danger';

interface AlertProps {
  tone?: AlertTone;
  children: ReactNode;
}

/** Danger alerts interrupt screen readers; info alerts are announced politely. */
export function Alert({ tone = 'info', children }: AlertProps) {
  const danger = tone === 'danger';
  return (
    <div className={danger ? 'section-message section-message--error' : 'section-message'}>
      <Icon name={danger ? 'error' : 'info'} />
      <div role={danger ? 'alert' : 'status'}>{children}</div>
    </div>
  );
}
