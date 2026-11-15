<div align="center">

# OpenBoard

**Plan, track and ship your work, on your own servers.**

An open source, self-hosted project management platform for teams that want the power of Jira or Azure DevOps Boards without giving up their data.

[![CI](https://github.com/elyosemite/openboard/actions/workflows/ci.yml/badge.svg)](https://github.com/elyosemite/openboard/actions/workflows/ci.yml)

**English** | [Português (Brasil)](README.pt-BR.md)

</div>

---

## Why OpenBoard?

- **Your data stays yours.** Run it on your own infrastructure with one command. Even the person who administers the installation cannot read your organization's work unless they are given access.
- **Self-hosted first, Cloud ready.** The same codebase serves a single team on one server and thousands of organizations in the cloud.
- **Access you can explain.** Two levels, organization and workspace, and nothing hidden in between. "Why can't she see this project?" has a one-line answer.
- **Built for engineering teams.** Work items, cycles, boards and docs in one place, with the rigor engineers expect: an OpenAPI contract, tested migrations and strict quality gates.
- **Open source.** Audit it, extend it, shape where it goes.

## Status

> [!NOTE]
> OpenBoard is in **early development** and not ready for production use yet. The foundations are in place; the first product features are next. Feedback and ideas are welcome.

## What's already here

- **One-command stack.** `docker compose up` starts PostgreSQL, applies the migrations and runs the API and the web app.
- **Production-grade API foundation.** Go service with health and readiness probes, consistent [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) error responses, structured logs and graceful shutdown.
- **Web app foundation.** React and TypeScript app organized by feature, with its own design system and light/dark themes.
- **Safe database evolution.** Versioned Liquibase migrations, and CI proves every one of them can be rolled back.
- **Quality on every change.** CI runs linting, complexity limits, test coverage, vulnerability scanning, API contract checks and image builds on every pull request.
- **A designed domain.** Organizations, workspaces, teams and projects are modeled up front, with a shared glossary and recorded decisions.

## Roadmap

The next major features, in order:

1. **Organizations, Workspaces and Teams.** Multi-tenant structure with open and restricted workspaces, so each area of a company gets its own space under one administration.
2. **Projects and Work Items.** Epics, features, tasks and bugs with rich descriptions, comments, attachments and relationships such as blocks, depends on and subtasks.
3. **Views.** List, Kanban board, calendar and timeline, with saved and shared filters.
4. **Eras.** Time-boxed cycles (like sprints) with goals, progress tracking, and burndown and burnup charts.
5. **Pages.** Documentation that lives next to the work, where an idea on a page becomes a work item in one click.

The full product vision is in [`docs/product-vision.md`](docs/product-vision.md).

## Quick start

Requirements: [Docker](https://docs.docker.com/get-docker/).

```bash
git clone https://github.com/elyosemite/openboard.git
cd openboard
docker compose up -d --build
```

- Web app: http://localhost:3000
- API: http://localhost:8080 (`/healthz`, `/readyz`)

## Built with

Go · React · TypeScript · PostgreSQL · Liquibase · Docker · GitHub Actions

## For contributors

Development requires Go 1.26+, Node 24+, pnpm 10+ and [Task](https://taskfile.dev).

```bash
task check   # every quality gate: format, lint, complexity, tests with coverage
task --list  # all available tasks
```

| Path | What lives there |
|---|---|
| `backend/` | Go API (modular monolith) |
| `frontend/` | React + TypeScript web app |
| `database/` | Liquibase migrations and the image that applies them |
| `docs/adr/` | Architecture decision records |
| `docs/api/openapi.yaml` | API contract |
| `CONTEXT.md` | Domain glossary |

Standards: [API design](api.md) · [Database migrations](Database.md)

A license will be chosen before the first release.
