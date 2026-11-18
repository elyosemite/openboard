import { useId, type InputHTMLAttributes } from 'react';

interface TextFieldProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'id'> {
  label: string;
  /** Message shown under the field; also marks the input as invalid. */
  error?: string | undefined;
}

export function TextField({ label, error, ...rest }: TextFieldProps) {
  const id = useId();
  const errorId = `${id}-error`;
  const invalid = error !== undefined && error !== '';
  return (
    <div className="field">
      <label htmlFor={id} className="field__label">
        {label}
      </label>
      <div className={invalid ? 'input is-invalid' : 'input'}>
        <input
          {...rest}
          id={id}
          aria-invalid={invalid || undefined}
          aria-describedby={invalid ? errorId : undefined}
        />
      </div>
      {invalid && (
        <span id={errorId} className="field__help field__help--error">
          {error}
        </span>
      )}
    </div>
  );
}
