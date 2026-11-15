<div align="center">

# OpenBoard

**Planeje, acompanhe e entregue seu trabalho, nos seus próprios servidores.**

Uma plataforma open source e self-hosted de gerenciamento de projetos para equipes que querem o poder do Jira ou do Azure DevOps Boards sem abrir mão dos próprios dados.

[![CI](https://github.com/elyosemite/openboard/actions/workflows/ci.yml/badge.svg)](https://github.com/elyosemite/openboard/actions/workflows/ci.yml)

[English](README.md) | **Português (Brasil)**

</div>

---

## Por que o OpenBoard?

- **Seus dados continuam seus.** Rode na sua própria infraestrutura com um comando. Nem quem administra a instalação consegue ler o trabalho da sua organização sem receber acesso.
- **Self-hosted primeiro, pronto para a nuvem.** O mesmo código atende uma equipe num único servidor e milhares de organizações na nuvem.
- **Acesso que dá para explicar.** Dois níveis, organização e workspace, e nada escondido no meio. "Por que ela não vê esse projeto?" tem resposta em uma linha.
- **Feito para times de engenharia.** Work items, ciclos, boards e documentação num só lugar, com o rigor que engenheiros esperam: contrato OpenAPI, migrations testadas e quality gates rígidos.
- **Open source.** Audite, estenda e ajude a decidir para onde ele vai.

## Status

> [!NOTE]
> O OpenBoard está em **desenvolvimento inicial** e ainda não está pronto para uso em produção. As fundações estão prontas; as primeiras funcionalidades de produto são as próximas. Feedback e ideias são bem-vindos.

## O que já existe

- **Stack com um comando.** `docker compose up` sobe o PostgreSQL, aplica as migrations e roda a API e o app web.
- **Base de API pronta para produção.** Serviço em Go com probes de saúde e prontidão, respostas de erro padronizadas no [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457), logs estruturados e graceful shutdown.
- **Base do app web.** App em React e TypeScript organizado por funcionalidade, com design system próprio e temas claro e escuro.
- **Evolução segura do banco.** Migrations versionadas com Liquibase, e o CI comprova que cada uma pode ser desfeita.
- **Qualidade em toda mudança.** O CI roda lint, limites de complexidade, cobertura de testes, varredura de vulnerabilidades, checagem do contrato da API e build das imagens em todo pull request.
- **Um domínio desenhado.** Organizations, workspaces, teams e projects foram modelados de antemão, com glossário compartilhado e decisões registradas.

## Roadmap

As próximas grandes funcionalidades, em ordem:

1. **Organizations, Workspaces e Teams.** Estrutura multi-tenant com workspaces abertos e restritos, para cada área da empresa ter seu espaço sob uma única administração.
2. **Projects e Work Items.** Epics, features, tasks e bugs com descrição rica, comentários, anexos e relações como bloqueia, depende de e subtarefas.
3. **Views.** Lista, board Kanban, calendário e timeline, com filtros salvos e compartilhados.
4. **Eras.** Ciclos com prazo definido (como sprints), com objetivos, acompanhamento de progresso e gráficos de burndown e burnup.
5. **Pages.** Documentação ao lado do trabalho, em que uma ideia numa página vira um work item com um clique.

A visão completa do produto está em [`docs/product-vision.pt-BR.md`](docs/product-vision.pt-BR.md).

## Comece rápido

Requisito: [Docker](https://docs.docker.com/get-docker/).

```bash
git clone https://github.com/elyosemite/openboard.git
cd openboard
docker compose up -d --build
```

- App web: http://localhost:3000
- API: http://localhost:8080 (`/healthz`, `/readyz`)

## Construído com

Go · React · TypeScript · PostgreSQL · Liquibase · Docker · GitHub Actions

## Para quem contribui

O desenvolvimento requer Go 1.26+, Node 24+, pnpm 10+ e [Task](https://taskfile.dev).

```bash
task check   # todos os quality gates: formatação, lint, complexidade, testes com cobertura
task --list  # todas as tarefas disponíveis
```

| Pasta | O que tem lá |
|---|---|
| `backend/` | API em Go (modular monolith) |
| `frontend/` | App web em React + TypeScript |
| `database/` | Migrations Liquibase e a imagem que as aplica |
| `docs/adr/` | Registros de decisões de arquitetura |
| `docs/api/openapi.yaml` | Contrato da API |
| `CONTEXT.md` | Glossário do domínio |

Padrões: [Design de API](api.md) · [Migrations de banco](Database.md)

Uma licença será escolhida antes da primeira release.
