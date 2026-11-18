# Domain docs

How the engineering skills should consume this repository's domain documentation. OpenBoard uses a **single context**.

## Before exploring, read

- **`CONTEXT.md`** at the root: the domain glossary.
- **`docs/adrs/NNNN-title/index.mdx`**: the ADRs that touch the area you are about to work in.

ADRs follow the EventCatalog format (ADR 0008). Do not use `docs/adr/`.

## File structure

```
/
├── CONTEXT.md
├── docs/adrs/
│   ├── 0003-organization-is-the-tenant/index.mdx
│   └── 0007-project-is-a-structural-boundary/index.mdx
└── backend/internal/
    ├── identity/
    └── organizations/
```

## Use the glossary's vocabulary

When your output names a domain concept (issue title, refactoring proposal, hypothesis, test name), use the term defined in `CONTEXT.md`. Do not use the synonyms the glossary says to avoid.

If the concept is not in the glossary yet, either you are inventing language the project does not use, and should reconsider, or there is a real gap, which should be noted for `/grill-with-docs`.

## Flag ADR conflicts

If your output contradicts an existing ADR, say so explicitly instead of ignoring it:

> _Contradicts ADR 0007 (Project is a structural boundary), but worth reopening because…_
