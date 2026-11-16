# OpenBoard — instruções para o Claude

Contexto completo do produto e da arquitetura: `docs/product-vision.md` (inglês) e `docs/product-vision.pt-BR.md` (português). Leia a seção relevante antes de cada tarefa.

- Use os termos do glossário de domínio em `CONTEXT.md` (código, API, testes, conversas).
- Papéis dos documentos: `README.md` é a vitrine pública do projeto; este `CLAUDE.md` guarda regras recorrentes; decisões arquiteturais ficam em ADRs (`docs/adrs/`, no EventCatalog); `CONTEXT.md` é só glossário. Não duplique conteúdo entre eles.
- O `README.md` existe só em inglês. Mantenha as versões em inglês e português da visão do produto sincronizadas.

## Como trabalhamos

- Entregue em **pequenos incrementos**: uma fatia por vez, que compila, tem testes e pode ser revisada sozinha. Não adiante módulos ou funcionalidades que não foram pedidos.
- Antes de codar algo não trivial, proponha um plano curto e espere aprovação.
- Simplicidade antes de abstração: só introduza camada, biblioteca ou serviço com necessidade concreta.
- Decisões arquiteturais relevantes viram ADR em `docs/adrs/NNNN-titulo/index.mdx` (formato de ADR do EventCatalog, ADR 0008), com `summary` e `owners: [maintainers]`; `pnpm lint` em `docs/` valida.

## Qualidade de código (obrigatório em todo código escrito ou alterado)

- Complexidade ciclomática **≤ 10** por função.
- Complexidade cognitiva **≤ 15** por função.
- Aninhamento **≤ 3** níveis.
- Cobertura de código novo **≥ 80%**.
- Testes **sempre passando**. Nunca desabilite, pule ou afrouxe testes ou regras para fazer passar.
- Análise com **SonarQube local** via MCP `sonarqube-local` (definido em `.mcp.json`, servidor `http://localhost:9000`). **Proibido** usar SonarQube Cloud ou qualquer conta/servidor SonarQube hospedado — inclusive o MCP `sonarqube` de outros projetos. Se o MCP local estiver indisponível, avise o usuário em vez de dar por concluído.
- Subir o servidor local: `docker start sonarqube` (criado com `docker run -d --name sonarqube -p 9000:9000 sonarqube:community`).

## Backend (Go)

- Modular Monolith: cada módulo em `internal/<modulo>/` com camadas domínio / aplicação / infraestrutura / http. Um módulo não importa a infraestrutura de outro.
- Regras de negócio não dependem de PostgreSQL nem de HTTP.
- Sempre: `context.Context` propagado, erros tratados de forma centralizada, logs estruturados (`log/slog`), configuração por variáveis de ambiente.
- Toda e qualquer migration é feita com **Liquibase** (nunca ferramenta de migration em Go nem SQL solto aplicado à mão). As migrations são aplicadas pelo pipeline de CD.
- Autorização sempre validada no backend.
- A Organization é o tenant (ADR 0003): todo dado de domínio pertence a uma Organization e nunca cruza para outra; Workspace só restringe visibilidade dentro dela.
- Ferramentas: `gofmt`, `go vet`, `golangci-lint`, `gocyclo -over 10`, `gocognit -over 15`.

## Banco de dados

- Para o banco de dados você deve consultar as instruções em Database.md

## Frontend (React + TypeScript)

- Organizado por funcionalidade (`src/features/<feature>/`), com Design System centralizado.
- Separar apresentação, estado e comunicação com a API.
- Ferramentas: ESLint, TypeScript em modo `strict`, Prettier.

## Antes de dar uma tarefa por concluída

1. Backend: `gofmt`, `go vet`, `golangci-lint`, `gocyclo`, `gocognit` e `go test -cover ./...` sem erros.
2. Frontend: Prettier, ESLint, `tsc --noEmit` e testes com cobertura sem erros.
3. Análise no SonarQube local (MCP) sem issues novas e com os limites acima atendidos.
4. Não faça push sem pedido explícito.

## Commits

- Faça **commits pequenos**: um commit por passo lógico (uma fatia que compila e passa nos testes), nunca um commit gigante no fim da tarefa.
- Toda mensagem de commit é gerada com a skill **`caveman-commit`** (Conventional Commits, em inglês).

## Como os endopints na API devem ser escritos

Você deve olhar api.md
