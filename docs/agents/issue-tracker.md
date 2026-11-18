# Issue tracker: markdown local

Issues e PRDs deste repositório ficam em arquivos markdown dentro de `.scratch/`. A pasta é ignorada pelo Git, então esse conteúdo existe só na máquina local.

## Convenções

- Uma funcionalidade por diretório: `.scratch/<feature-slug>/`
- O PRD fica em `.scratch/<feature-slug>/PRD.md`
- As issues de implementação ficam em `.scratch/<feature-slug>/issues/<NN>-<slug>.md`, numeradas a partir de `01`
- O estado de triagem é uma linha `Status:` perto do topo de cada issue (os valores estão em `triage-labels.md`)
- Comentários e histórico da conversa são acrescentados ao fim do arquivo, sob o título `## Comments`

## Quando uma skill disser "publique no issue tracker"

Crie um arquivo em `.scratch/<feature-slug>/`, criando o diretório se necessário.

## Quando uma skill disser "busque o ticket relevante"

Leia o arquivo no caminho indicado. Normalmente o usuário informa o caminho ou o número da issue.
