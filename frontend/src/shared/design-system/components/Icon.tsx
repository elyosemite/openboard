import spriteUrl from 'ocaml-design-system/sprite.svg?url';

interface IconProps {
  /** Name of the icon in the design system's sprite, without the `i-` prefix. */
  name: string;
}

/** A decorative icon: the text next to it gives the meaning. */
export function Icon({ name }: IconProps) {
  return (
    <svg className="i" aria-hidden="true" focusable="false">
      <use href={`${spriteUrl}#i-${name}`} />
    </svg>
  );
}
