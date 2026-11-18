interface SpinnerProps {
  label?: string;
}

/** Announces that content is loading. */
export function Spinner({ label = 'Loading' }: SpinnerProps) {
  return (
    <span role="status" aria-label={label}>
      <span className="spinner spinner--sm" aria-hidden="true" />
    </span>
  );
}
