import { useId, type ReactNode } from 'react';
import './Card.css';

interface CardProps {
  title: string;
  children: ReactNode;
}

export function Card({ title, children }: CardProps) {
  const titleId = useId();
  return (
    <section className="ob-card" aria-labelledby={titleId}>
      <h2 id={titleId} className="ob-card__title">
        {title}
      </h2>
      <div className="ob-card__body">{children}</div>
    </section>
  );
}
