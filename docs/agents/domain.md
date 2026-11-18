# Documentação de domínio

Como as skills de engenharia devem consumir a documentação de domínio deste repositório. O OpenBoard usa **contexto único**.

## Antes de explorar, leia

- **`CONTEXT.md`** na raiz: o glossário do domínio.
- **`docs/adrs/NNNN-titulo/index.mdx`**: os ADRs que tocam a área em que você vai trabalhar.

Os ADRs seguem o formato do EventCatalog (ADR 0008). Não use `docs/adr/`.

## Estrutura de arquivos

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

## Use o vocabulário do glossário

Quando sua saída nomear um conceito do domínio (título de issue, proposta de refatoração, hipótese, nome de teste), use o termo definido em `CONTEXT.md`. Não use os sinônimos que o glossário manda evitar.

Se o conceito ainda não estiver no glossário, há duas possibilidades: você está inventando uma linguagem que o projeto não usa, e deve reconsiderar, ou existe uma lacuna real, que deve ser anotada para `/grill-with-docs`.

## Aponte conflitos com ADRs

Se sua saída contradisser um ADR existente, diga isso explicitamente em vez de ignorá-lo:

> _Contradiz o ADR 0007 (Project é uma fronteira estrutural), mas vale reabrir porque…_
