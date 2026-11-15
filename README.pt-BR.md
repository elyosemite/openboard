# OpenBoard

[English](README.md) | **Português (Brasil)**

## Como rodar localmente

Pré-requisitos: Docker. Para desenvolver: Go 1.26+, Node 24+ e pnpm 10+.

```bash
docker compose up -d --build   # postgres + liquibase + backend + frontend
```

- Frontend: http://localhost:3000
- API: http://localhost:8080 (`/healthz`, `/readyz`)

| Pasta | Conteúdo |
|---|---|
| `backend/` | API em Go (modular monolith) |
| `frontend/` | App React + TypeScript |
| `database/` | Migrations Liquibase e a imagem que as aplica |
| `docs/` | ADRs (`docs/adr`) e contrato OpenAPI (`docs/api/openapi.yaml`) |

Verificações de qualidade: `task check` (veja o `Taskfile.yml`; instale o [Task](https://taskfile.dev) com `go install github.com/go-task/task/v3/cmd/task@latest`). Padrões do projeto: [`api.md`](api.md) e [`Database.md`](Database.md).

## 1. Visão geral

O **OpenBoard** será uma plataforma open source de gerenciamento de projetos e acompanhamento de trabalho (*Project Management & Work Tracking*), inspirada em soluções consolidadas como Azure DevOps Boards, Jira e Plane.

Seu propósito é oferecer uma alternativa moderna, flexível e extensível para que empresas, equipes de engenharia e organizações possam planejar, organizar, acompanhar e analisar seu trabalho em um único ambiente.

A plataforma será desenvolvida com foco em duas modalidades de implantação: **Self-hosted e Cloud**. O modelo Self-hosted permitirá que organizações instalem e executem sua própria instância do OpenBoard em servidores internos ou infraestrutura cloud de sua preferência, mantendo controle sobre seus dados, configurações e políticas de acesso.

O OpenBoard deverá atender desde pequenas equipes até organizações com estruturas complexas, múltiplos departamentos, equipes e projetos, oferecendo recursos para colaboração, planejamento, execução, documentação e análise de atividades.

Sua arquitetura será organizada em subdomínios funcionais, com responsabilidades bem definidas, permitindo que a plataforma evolua de maneira modular sem comprometer a consistência do domínio central.

## 2. Estrutura organizacional

A estrutura organizacional será o principal fundamento do OpenBoard. Ela definirá como usuários, equipes e projetos se relacionam dentro da plataforma.

Os quatro conceitos centrais serão Organizations, Workspaces, Teams e Projects.

### 2.1. Organizations

Organizations representarão as organizações ou empresas que utilizam o OpenBoard.

Cada organização terá seu próprio contexto administrativo, podendo gerenciar usuários, membros, configurações e diferentes espaços de trabalho.

Uma Organization será responsável por centralizar a administração e estabelecer as fronteiras organizacionais de acesso e colaboração.

### 2.2. Workspaces

Workspaces serão espaços de trabalho independentes dentro de uma Organization, permitindo separar diferentes áreas, departamentos, unidades de negócio ou contextos de trabalho.

Cada Workspace poderá possuir suas próprias equipes, projetos, configurações e estruturas de organização.

O objetivo é permitir que uma mesma organização mantenha diferentes ambientes de trabalho sem misturar atividades, informações e responsabilidades.

### 2.3. Teams

Teams representarão grupos de pessoas que trabalham em conjunto dentro de um Workspace.

Cada equipe poderá possuir membros, responsáveis, configurações e projetos associados.

As equipes serão fundamentais para organizar responsabilidades, distribuir atividades e acompanhar a execução do trabalho.

### 2.4. Projects

Projects serão as unidades centrais de gerenciamento e execução do trabalho.

Cada projeto poderá possuir seus próprios Work Items, Eras, Modules, Views, Pages e indicadores de Analytics.

Os projetos poderão estar associados a uma ou mais equipes, permitindo diferentes modelos de colaboração conforme as necessidades da organização.

## 3. Subdomínios funcionais

A plataforma será organizada em subdomínios que representam diferentes capacidades do sistema. Alguns deles constituirão o núcleo inicial do produto, enquanto outros serão desenvolvidos progressivamente.

### 3.1. Work Items

Work Items serão a unidade fundamental de trabalho do OpenBoard.

Representarão atividades que precisam ser planejadas, executadas e acompanhadas ao longo do ciclo de vida de um projeto.

Cada item poderá conter título, descrição, responsável, prioridade, status, prazo, estimativa de esforço, labels, anexos e comentários.

Um editor de texto rico permitirá documentar requisitos, registrar decisões técnicas, compartilhar imagens e anexar arquivos.

Também será possível estabelecer relações entre diferentes Work Items, como dependências, bloqueios, itens relacionados, subtarefas e hierarquias de trabalho.

Tipos de Work Items poderão incluir Epics, Features, Tasks, Bugs e Subtasks, com possibilidade de personalização conforme as necessidades de cada organização.

### 3.2. Eras

Eras serão períodos ou ciclos de execução utilizados para organizar o trabalho de uma equipe.

Conceitualmente, serão semelhantes aos Sprints do Scrum e aos Cycles de ferramentas como Plane.

Cada Era terá data de início, data de encerramento, objetivos e um conjunto de Work Items associados.

As equipes poderão planejar suas atividades, acompanhar o progresso e avaliar o que foi concluído durante determinado período.

O acompanhamento incluirá indicadores como progresso geral, itens concluídos e pendentes, estimativas de esforço, Burndown Charts e Burnup Charts.

### 3.3. Modules

Modules permitirão dividir projetos maiores em unidades menores e mais gerenciáveis.

Cada módulo poderá representar uma funcionalidade, subsistema, componente de software, departamento ou qualquer outra divisão lógica definida pela organização.

Os módulos agruparão Work Items relacionados, permitindo acompanhar o progresso de iniciativas específicas sem perder a visão geral do projeto.

### 3.4. Views

Views oferecerão diferentes formas de visualizar e organizar os Work Items.

As visualizações poderão incluir:

- List View: organização tabular dos itens.
- Board View: visualização Kanban organizada por status.
- Calendar View: acompanhamento de prazos e datas.
- Timeline View: visualização temporal das atividades e dependências.

Os usuários poderão criar filtros personalizados utilizando propriedades como responsável, prioridade, status, módulo, Era, labels e datas.

As Views poderão ser salvas, compartilhadas entre membros da equipe e reutilizadas.

### 3.5. Pages

Pages serão espaços de documentação integrados ao gerenciamento de projetos.

Permitirão documentar ideias, requisitos, decisões arquiteturais, atas de reuniões, especificações técnicas e outros conhecimentos relacionados ao projeto.

Cada página contará com um editor de texto rico, suporte a imagens, hyperlinks, anexos, formatação avançada e referências a Work Items existentes.

Uma ideia registrada em uma página poderá ser transformada em um Work Item, aproximando documentação e execução.

A arquitetura também deverá permitir futuras integrações com recursos de inteligência artificial para auxiliar na criação, organização e sumarização de documentos.

### 3.6. Analytics

Analytics fornecerá uma visão consolidada dos dados dos projetos, permitindo que equipes e gestores acompanhem a evolução do trabalho por meio de indicadores e visualizações.

Entre os indicadores previstos estão:

- Distribuição de Work Items por status, prioridade e responsável.
- Velocidade de execução das equipes.
- Progresso das Eras.
- Burndown e Burnup.
- Lead Time e Cycle Time.
- Identificação de atividades bloqueadas.
- Evolução das entregas ao longo do tempo.

O objetivo será facilitar a identificação de gargalos, compreender a capacidade operacional das equipes e apoiar decisões baseadas em dados.

## 4. Subdomínios de suporte e infraestrutura funcional

Além dos subdomínios diretamente relacionados à organização e execução do trabalho, o OpenBoard contará com capacidades transversais responsáveis por sustentar a operação da plataforma.

Esses subdomínios serão desenvolvidos progressivamente.

### 4.1. Workflow Engine

Mecanismo responsável por definir e executar os fluxos de trabalho dos Work Items.

Permitirá configurar estados, transições, regras e restrições, possibilitando que cada organização adapte seus processos operacionais.

### 4.2. Automations

Sistema de automações baseado em eventos e condições.

Permitirá executar ações automaticamente a partir de determinados acontecimentos, como alterações de status, atribuições, vencimentos e criação de Work Items.

### 4.3. Integrations

Camada responsável pela integração com ferramentas e serviços externos.

Poderá oferecer integrações com GitHub, GitLab, Bitbucket, Slack, Microsoft Teams e outras plataformas utilizadas pelas equipes.

### 4.4. RBAC

Sistema de controle de acesso baseado em papéis (*Role-Based Access Control*).

Será responsável por gerenciar permissões e controlar quais recursos e operações estarão disponíveis para cada usuário dentro de uma Organization, Workspace, Team ou Project.

### 4.5. Audit Logs

Sistema de auditoria responsável por registrar eventos e alterações relevantes dentro da plataforma.

Permitirá rastrear operações, identificar responsáveis por alterações e consultar o histórico de atividades.

### 4.6. Notifications

Sistema de notificações responsável por informar os usuários sobre acontecimentos relevantes.

Incluirá notificações sobre atribuições, comentários, menções, alterações de status, prazos e outras atividades relacionadas aos projetos.

## 5. Princípios do projeto

O desenvolvimento do OpenBoard será orientado pelos seguintes princípios:

- **Open Source:** código-fonte aberto, permitindo auditoria, colaboração e evolução pela comunidade.
- **Self-hosted First:** projetado para funcionar em infraestrutura própria, sem dependência obrigatória de serviços externos.
- **Cloud Ready:** arquitetura preparada para oferecer uma modalidade Cloud.
- **Data Ownership:** respeito à propriedade e ao controle dos dados pelas organizações.
- **Modular Architecture:** subdomínios com responsabilidades e limites bem definidos.
- **Extensibility:** arquitetura preparada para integrações, plugins, automações e funcionalidades desenvolvidas pela comunidade.
- **Multi-tenant:** suporte a múltiplas organizações e espaços de trabalho, com isolamento adequado de dados e permissões.
- **Customizable Workflows:** possibilidade de adaptar fluxos de trabalho, propriedades e processos às necessidades de cada organização.
- **Developer Experience:** experiência especialmente cuidadosa para equipes de engenharia de software, sem limitar o uso a esse público.

## 6. Objetivo de longo prazo

O OpenBoard pretende se tornar uma plataforma completa de gerenciamento de trabalho, capaz de atender diferentes modelos de organização e metodologias de execução.

Mais do que reproduzir funcionalidades existentes em ferramentas comerciais, o projeto busca oferecer uma alternativa aberta, independente e extensível, na qual organizações possam construir seus próprios processos de trabalho sem abrir mão da propriedade dos dados, da flexibilidade operacional e da liberdade tecnológica.

## 7. Especificações técnicas e arquitetura

O OpenBoard deverá ser desenvolvido seguindo boas práticas de engenharia de software, priorizando manutenibilidade, segurança, escalabilidade, testabilidade e simplicidade arquitetural.

As decisões técnicas deverão considerar tanto a modalidade Self-hosted quanto a modalidade Cloud, evitando dependências desnecessárias de provedores específicos e garantindo que a aplicação possa ser executada em diferentes ambientes de infraestrutura.

### 7.1. Stack tecnológica

As tecnologias principais do projeto serão:

- **Backend:** Go (Golang).
- **Frontend:** React.
- **Linguagem do frontend:** TypeScript.
- **Build e dependências do frontend:** Node.js, utilizando um gerenciador de pacotes consistente.
- **Containerização:** Docker e Docker Compose.
- **CI/CD:** GitHub Actions.
- **Controle de versão:** Git e GitHub.
- **Banco de dados principal:** PostgreSQL.

A escolha das tecnologias deverá priorizar estabilidade, maturidade do ecossistema, facilidade de manutenção e disponibilidade de ferramentas para desenvolvimento e observabilidade.

### 7.2. Arquitetura do backend

O backend será desenvolvido em Go, utilizando uma arquitetura organizada por domínios e responsabilidades.

O projeto deverá adotar princípios de Domain-Driven Design (DDD) quando agregarem valor à organização do domínio, mantendo separação clara entre regras de negócio, aplicação, infraestrutura e interfaces de comunicação.

A organização do código deverá favorecer modularidade, baixo acoplamento e alta coesão, evitando tanto a criação de um monólito desorganizado quanto a fragmentação prematura em microserviços.

A arquitetura inicial deverá seguir uma abordagem de **Modular Monolith**, na qual os subdomínios possuem limites bem definidos e podem evoluir independentemente dentro de uma única aplicação.

Os principais módulos deverão incluir:

- Organizations
- Workspaces
- Teams
- Projects
- Work Items
- Eras
- Modules
- Views
- Pages
- Analytics
- Workflow Engine
- Automations
- Integrations
- RBAC
- Audit Logs
- Notifications

Cada módulo deverá manter suas próprias responsabilidades, regras de negócio e interfaces de comunicação.

O backend deverá disponibilizar uma API HTTP bem definida, preferencialmente RESTful, com contratos consistentes, versionamento e documentação OpenAPI.

Requisitos adicionais:

- Tratamento centralizado de erros.
- Validação de entradas e regras de negócio.
- Logs estruturados.
- Context propagation e suporte a cancelamento de operações.
- Gerenciamento adequado de concorrência.
- Configuração por variáveis de ambiente.
- Migrations versionadas para o banco de dados.
- Testes unitários, de integração e de arquitetura.
- Graceful shutdown.
- Health checks e readiness checks.

### 7.3. Arquitetura do frontend

O frontend será desenvolvido em React com TypeScript.

A aplicação deverá adotar uma arquitetura modular, organizada por funcionalidades e domínios, evitando uma estrutura excessivamente centralizada ou baseada apenas em componentes genéricos.

A interface deverá ser responsiva, acessível e consistente, priorizando uma experiência de uso intuitiva para equipes que trabalham diariamente com gerenciamento de projetos.

A estrutura deverá contemplar:

- Componentes reutilizáveis.
- Design System centralizado.
- Gerenciamento consistente de estado.
- Camada de comunicação com a API.
- Tratamento padronizado de erros.
- Controle de autenticação e autorização na interface.
- Rotas organizadas por contexto funcional.
- Formulários com validação.
- Estratégia de cache e sincronização de dados.
- Testes unitários e de componentes.
- Testes de integração e end-to-end.

O frontend deverá manter uma separação clara entre apresentação, estado da aplicação, regras de interação e comunicação com o backend.

### 7.4. Banco de dados e persistência

O PostgreSQL será o banco de dados relacional principal da plataforma.

A modelagem deverá refletir os limites do domínio e preservar a integridade dos relacionamentos entre Organizations, Workspaces, Teams, Projects e demais entidades.

Requisitos:

- Migrations versionadas e reproduzíveis.
- Integridade referencial.
- Índices planejados de acordo com os padrões de consulta.
- Transações para operações que exijam consistência.
- Estratégias de paginação para coleções extensas.
- Auditoria de alterações relevantes.
- Isolamento adequado dos dados entre organizações e workspaces.
- Estratégias de backup e restauração.

O acesso ao banco de dados deverá ocorrer por meio de uma camada de persistência bem definida, evitando que regras de negócio dependam diretamente de detalhes específicos do PostgreSQL.

### 7.5. Containerização e infraestrutura

O OpenBoard deverá ser distribuído utilizando containers Docker.

O repositório deverá possuir Dockerfiles específicos para backend e frontend, priorizando imagens pequenas, builds reproduzíveis e segurança durante a execução.

Os Dockerfiles deverão seguir boas práticas como:

- Multi-stage builds.
- Imagens base oficiais e confiáveis.
- Execução com usuário não privilegiado.
- Separação entre dependências de build e runtime.
- Redução da superfície de ataque.
- Configuração por variáveis de ambiente.
- Health checks quando aplicáveis.

O projeto deverá disponibilizar um Docker Compose para facilitar a execução local e a instalação self-hosted.

A infraestrutura inicial deverá permitir executar os componentes essenciais da plataforma sem exigir serviços externos obrigatórios.

### 7.6. CI/CD e automação de desenvolvimento

O processo de integração e entrega contínuas será implementado utilizando GitHub Actions.

O pipeline deverá automatizar verificações de qualidade, testes e geração de artefatos, garantindo que alterações no código não comprometam a estabilidade do projeto.

O pipeline deverá contemplar:

- Verificação de formatação e lint.
- Análise estática de código.
- Compilação do backend.
- Build do frontend.
- Testes unitários.
- Testes de integração.
- Verificação de dependências e vulnerabilidades.
- Build das imagens Docker.
- Publicação de imagens no GitHub Container Registry (GHCR).
- Geração de artefatos versionados.

O fluxo deverá considerar diferentes etapas, como Pull Requests, integração na branch principal, criação de releases e publicação de versões.

As releases deverão possuir versionamento consistente, preferencialmente seguindo Semantic Versioning.

### 7.7. Qualidade e testes

A qualidade deverá ser tratada como parte integrante do processo de desenvolvimento.

O projeto deverá possuir uma estratégia de testes em diferentes níveis:

- **Unit Tests:** validação isolada de regras de negócio e componentes.
- **Integration Tests:** validação da comunicação entre módulos e infraestrutura.
- **End-to-End Tests:** validação dos principais fluxos de utilização da plataforma.
- **Architecture Tests:** verificação de dependências e limites arquiteturais.

O código deverá passar por verificações automatizadas antes de ser integrado à branch principal.

Também deverão existir padrões de contribuição, revisão de código e documentação técnica para orientar colaboradores.

### 7.8. Segurança

Segurança deverá ser considerada desde a concepção da arquitetura.

A plataforma deverá implementar mecanismos adequados de autenticação, autorização e isolamento de dados.

Entre os requisitos iniciais:

- Armazenamento seguro de credenciais.
- Proteção contra acesso não autorizado.
- Validação e sanitização de entradas.
- Proteção contra vulnerabilidades comuns, incluindo as categorias relevantes do OWASP Top 10.
- Gerenciamento seguro de secrets.
- Princípio do menor privilégio.
- Controle de acesso entre Organizations, Workspaces, Teams e Projects.
- Registro de operações sensíveis.
- Atualização e monitoramento de dependências vulneráveis.

A autorização deverá ser validada no backend, não podendo depender exclusivamente dos controles implementados na interface.

### 7.9. Observabilidade e operação

O OpenBoard deverá possuir mecanismos de observabilidade que permitam compreender o comportamento da aplicação em ambientes de desenvolvimento e produção.

A implementação deverá contemplar:

- Logs estruturados.
- Métricas de aplicação e infraestrutura.
- Rastreamento distribuído quando necessário.
- Health checks.
- Monitoramento de erros e falhas.
- Identificação de operações lentas.
- Correlação de requisições.

A arquitetura deverá permitir futuras integrações com OpenTelemetry e ferramentas como Prometheus e Grafana.

### 7.10. Documentação e manutenção

O repositório deverá manter documentação técnica atualizada, incluindo:

- Instruções de instalação e execução local.
- Guia de configuração self-hosted.
- Arquitetura geral da aplicação.
- Documentação dos módulos e seus limites.
- Contratos e documentação da API.
- Estratégia de testes.
- Processo de contribuição.
- Procedimentos de build, release e atualização.

As decisões arquiteturais relevantes deverão ser documentadas por meio de Architecture Decision Records (ADRs), permitindo compreender os motivos das escolhas técnicas e sua evolução ao longo do projeto.

### 7.11. Princípios arquiteturais

O desenvolvimento técnico do OpenBoard deverá seguir estes princípios:

1. **Simplicidade antes da complexidade:** não introduzir abstrações, serviços ou tecnologias sem uma necessidade concreta.
2. **Modularidade:** manter os domínios isolados e com responsabilidades bem definidas.
3. **Baixo acoplamento:** evitar dependências desnecessárias entre módulos.
4. **Alta coesão:** manter regras e funcionalidades relacionadas dentro de seus respectivos domínios.
5. **Testabilidade:** projetar componentes que possam ser testados de maneira isolada.
6. **Segurança por padrão:** tratar segurança como requisito arquitetural, não como funcionalidade posterior.
7. **Portabilidade:** evitar dependências obrigatórias de provedores específicos de infraestrutura.
8. **Evolução incremental:** permitir que a plataforma cresça sem exigir reescritas arquiteturais constantes.
