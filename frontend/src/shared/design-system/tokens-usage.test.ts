const sources = import.meta.glob<string>(
  ['/src/**/*.css', '/src/**/*.tsx', '!/src/**/*.test.tsx'],
  {
    query: '?raw',
    import: 'default',
    eager: true,
  },
);

const forbidden: { rule: string; pattern: RegExp }[] = [
  { rule: 'hexadecimal color', pattern: /#[0-9a-f]{3,8}\b/i },
  { rule: 'functional color', pattern: /\b(?:rgb|hsl)a?\(/i },
  { rule: 'palette primitive', pattern: /--palette-/ },
  { rule: 'token defined by OpenBoard', pattern: /^\s*--[\w-]+\s*:/m },
];

// ADR 0011: styles come only from the OCaml Design System's semantic tokens.
describe('design system usage', () => {
  it('finds the source files', () => {
    expect(Object.keys(sources)).toContain('/src/shared/design-system/styles/global.css');
  });

  describe.each(Object.entries(sources))('%s', (_path, source) => {
    it.each(forbidden)('has no $rule', ({ pattern }) => {
      expect(source).not.toMatch(pattern);
    });
  });
});
