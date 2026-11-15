# Database

All schema changes are Liquibase migrations under `migrations/`, organized by object type. The rules (naming, mandatory rollbacks, include order) live in [`Database.md`](../Database.md).

- Entry point: `migrations/changelog/db.changelog-master.yaml`.
- Locally, `docker compose up` applies the changelog through the `liquibase` service before the backend starts.
- In shared environments, migrations are applied only by the CD pipeline.
