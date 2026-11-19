import type { ButtonHTMLAttributes } from 'react';

export type ButtonVariant = 'primary' | 'secondary';

const variantClasses: Record<ButtonVariant, string> = {
  primary: 'btn btn--primary',
  secondary: 'btn btn--outline',
};

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  /** Shows the action is running and blocks further clicks. */
  busy?: boolean;
}

export function Button({
  variant = 'primary',
  busy = false,
  type = 'button',
  disabled,
  children,
  ...rest
}: ButtonProps) {
  return (
    <button
      {...rest}
      type={type}
      className={variantClasses[variant]}
      disabled={disabled === true || busy}
      aria-busy={busy || undefined}
    >
      {children}
    </button>
  );
}
