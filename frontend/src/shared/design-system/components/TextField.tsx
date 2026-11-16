import { useId, type InputHTMLAttributes } from 'react';
import './TextField.css';

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
    <div className="ob-field">
      <label htmlFor={id} className="ob-field__label">
        {label}
      </label>
      <input
        {...rest}
        id={id}
        className="ob-field__input"
        aria-invalid={invalid || undefined}
        aria-describedby={invalid ? errorId : undefined}
      />
      {invalid && (
        <p id={errorId} className="ob-field__error">
          {error}
        </p>
      )}
    </div>
  );
}
