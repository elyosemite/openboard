import type { UseQueryResult } from '@tanstack/react-query';
import { useId, type ReactNode } from 'react';
import { errorMessage } from '../../../shared/api/httpClient';
import { Alert, Button, Spinner } from '../../../shared/design-system';
import type { Page } from '../api/organizationsApi';
import './Organizations.css';

interface ListPanelProps<T> {
  title: string;
  query: UseQueryResult<Page<T>>;
  emptyTitle: string;
  emptyDescription: string;
  /** Shown at the end of the panel's header. */
  action?: ReactNode;
  /** Shown between the header and the list. */
  form?: ReactNode;
  renderItem: (item: T) => ReactNode;
  onPageChange: (page: number) => void;
}

/** A titled panel with one page of a collection and its loading, error and empty states. */
export function ListPanel<T>({ title, query, action, form, ...list }: ListPanelProps<T>) {
  const titleId = useId();
  return (
    <section className="panel panel--clip" aria-labelledby={titleId}>
      <header className="panel__head">
        <h2 id={titleId} className="panel__title">
          {title}
        </h2>
        {action && <div className="list-panel__action">{action}</div>}
      </header>
      {form}
      <ListBody query={query} {...list} />
    </section>
  );
}

type ListBodyProps<T> = Pick<
  ListPanelProps<T>,
  'query' | 'emptyTitle' | 'emptyDescription' | 'renderItem' | 'onPageChange'
>;

function ListBody<T>({
  query,
  emptyTitle,
  emptyDescription,
  renderItem,
  onPageChange,
}: ListBodyProps<T>) {
  if (query.isPending) {
    return (
      <div className="panel__body">
        <Spinner />
      </div>
    );
  }
  if (query.isError) {
    return (
      <div className="panel__body">
        <Alert tone="danger">{errorMessage(query.error)}</Alert>
      </div>
    );
  }
  if (query.data.items.length === 0) {
    return (
      <div className="empty">
        <h3 className="empty__title">{emptyTitle}</h3>
        <p className="empty__desc">{emptyDescription}</p>
      </div>
    );
  }
  return (
    <>
      <ul className="entity-list">{query.data.items.map(renderItem)}</ul>
      <ListPager
        page={query.data.page}
        totalPages={query.data.totalPages}
        onPageChange={onPageChange}
      />
    </>
  );
}

interface ListPagerProps {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

function ListPager({ page, totalPages, onPageChange }: ListPagerProps) {
  if (totalPages <= 1) {
    return null;
  }
  return (
    <footer className="panel__foot">
      <span>
        Page {page} of {totalPages}
      </span>
      <div className="list-pager">
        <Button
          variant="secondary"
          disabled={page <= 1}
          onClick={() => {
            onPageChange(page - 1);
          }}
        >
          Previous
        </Button>
        <Button
          variant="secondary"
          disabled={page >= totalPages}
          onClick={() => {
            onPageChange(page + 1);
          }}
        >
          Next
        </Button>
      </div>
    </footer>
  );
}
