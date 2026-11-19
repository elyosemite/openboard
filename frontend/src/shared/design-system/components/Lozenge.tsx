import type { ReactNode } from 'react';

export type LozengeTone = 'neutral' | 'success' | 'danger' | 'info';

const toneClasses: Record<LozengeTone, string> = {
  neutral: 'lozenge',
  success: 'lozenge lozenge--success',
  danger: 'lozenge lozenge--removed',
  info: 'lozenge lozenge--inprogress',
};

interface LozengeProps {
  tone?: LozengeTone;
  children: ReactNode;
}

export function Lozenge({ tone = 'neutral', children }: LozengeProps) {
  return <span className={toneClasses[tone]}>{children}</span>;
}
