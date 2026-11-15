# internal

- `platform/`: cross-cutting technical building blocks (config, logging, HTTP helpers, health probes, server lifecycle). It holds no business rules.
- `<module>/`: one directory per domain module (organizations, workspaces, ...), each split into `domain/`, `application/`, `infrastructure/` and `http/`. A module never imports another module's `infrastructure/`.
