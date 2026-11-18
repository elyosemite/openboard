import { useId, type ReactNode } from 'react';

interface CardProps {
  title: string;
  children: ReactNode;
}

export function Card({ title, children }: CardProps) {
  const titleId = useId();
  return (
    <section className="panel" aria-labelledby={titleId}>
      <header className="panel__head">
        <h2 id={titleId} className="panel__title">
          {title}
        </h2>
      </header>
      <div className="panel__body">{children}</div>
    </section>
  );
}
