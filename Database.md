# Database Migration Standards — Liquibase

All database schema changes and data migrations in the OpenBoard project MUST be managed through **Liquibase**. Direct, unmanaged database modifications are prohibited.

The migration structure must be organized by database object type. Each type of object must have its own dedicated directory. Do not mix tables, indexes, triggers, functions, constraints, seeds, or other database objects in the same directory.

## 1. Directory Structure

Organize Liquibase migration files using the following structure:

```text
database/
└── migrations/
    ├── changelog/
    │   └── db.changelog-master.yaml
    │
    ├── tables/
    ├── indexes/
    ├── constraints/
    ├── functions/
    ├── procedures/
    ├── triggers/
    ├── types/
    ├── sequences/
    ├── views/
    ├── materialized-views/
    ├── seeds/
    ├── permissions/
    └── rollback/
```

Each directory must contain only migrations related to its respective database object category.

Additional directories may be created when a new database object category is introduced, provided that the separation principle is maintained.

## 2. Migration File Naming Convention

Every migration file MUST start with a timestamp prefix.

Use the following naming convention:

```text
YYYYMMDDHHMMSS_description.sql
```

Examples:

```text
20261115103000_create_organizations_table.sql
20261115103100_create_workspaces_table.sql
20261115103200_create_work_items_indexes.sql
20261115103300_create_user_defined_types.sql
20261115103400_create_audit_trigger.sql
20261115103500_insert_initial_permissions.sql
```

Requirements:

- Timestamps must use the `YYYYMMDDHHMMSS` format.
- Timestamps must represent the migration creation time.
- Every migration must have a unique timestamp prefix.
- Use lowercase snake_case for descriptive names.
- Names must clearly describe the operation being performed.
- Never rename a migration that has already been applied to a shared or production database.
- Never reuse an existing migration timestamp.

## 3. Liquibase Changelog Organization

The master changelog must be the single entry point for all migrations.

Use Liquibase `includeAll` or explicit `include` declarations to load migration files from their respective directories.

The loading order must be deterministic and must respect database dependencies.

For example:

1. User-defined types and sequences.
2. Tables.
3. Constraints.
4. Indexes.
5. Functions and procedures.
6. Triggers.
7. Views and materialized views.
8. Seeds and initial reference data.
9. Permissions and other database configuration.

Do not assume that alphabetical directory ordering is sufficient to guarantee dependency correctness. Explicitly define the execution order in the master changelog.

Use Liquibase SQL formatted changelogs where appropriate, ensuring that each migration is properly registered and identifiable by Liquibase.

## 4. Mandatory Rollback Instructions

**Every migration MUST have an explicitly defined rollback instruction.**

Rollback behavior must be intentionally designed by the developer and must never be left to assumptions or automatic generation.

For SQL formatted changelogs, use Liquibase rollback sections:

```sql
--liquibase formatted sql

--changeset openboard:20261115103000
CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--rollback DROP TABLE organizations;
```

For more complex migrations, rollback instructions may use multiple SQL statements.

Rollback requirements:

- Every changeset must have a corresponding rollback definition.
- Rollbacks must reverse the original migration as accurately as possible.
- Respect foreign key dependencies and object creation order.
- Do not use destructive rollback operations without understanding their impact.
- For data migrations, define how modified or inserted data will be reverted.
- If a migration cannot be safely reversed, explicitly document why and provide an approved recovery strategy.
- Never assume that Liquibase can automatically infer a safe rollback.

## 5. Database Object Separation

Follow these rules strictly:

- Table creation and table structure modifications belong in `tables/`.
- Index creation and removal belong in `indexes/`.
- Primary keys, foreign keys, unique constraints, check constraints, and other constraints belong in `constraints/`, unless they are inseparable from the initial table definition.
- Functions belong in `functions/`.
- Stored procedures belong in `procedures/`.
- Triggers belong in `triggers/`.
- Custom PostgreSQL types, enums, domains, and composite types belong in `types/`.
- Sequence definitions belong in `sequences/`.
- Views belong in `views/`.
- Materialized views belong in `materialized-views/`.
- Initial reference data and controlled seed data belong in `seeds/`.
- Database grants and permissions belong in `permissions/`.

Avoid mixing unrelated database objects in a single changeset or migration file.

When an operation necessarily involves multiple object types, separate the operations into appropriately categorized migration files and ensure their execution order is correct.

## 6. Development and Review Requirements

Before considering a database migration complete, the developer or AI agent MUST:

1. Verify that the migration is located in the correct directory.
2. Verify that its filename follows the timestamp convention.
3. Confirm that the timestamp is unique.
4. Ensure that the migration is registered in the master changelog.
5. Verify that the changeset has an explicit rollback.
6. Validate SQL syntax against the project's PostgreSQL version.
7. Verify object dependencies and execution order.
8. Test migration execution against a clean database.
9. Test rollback execution.
10. Verify that rollback does not leave unexpected database objects or inconsistent data.

Never modify previously applied migrations to accommodate new requirements. Create a new migration instead.

## 7. General Principles

- Liquibase is the single source of truth for database schema evolution.
- All schema changes must be version-controlled.
- Never make undocumented manual changes to shared environments.
- Migrations must be deterministic, auditable, and reproducible.
- Prefer small, focused migrations over large migrations combining unrelated changes.
- Preserve backward compatibility whenever possible, especially when application versions may run concurrently during deployment.
- Treat database migrations as production code, subject to the same review, testing, and quality standards as application code.

**These requirements are mandatory for all human developers and AI coding agents working on OpenBoard.**
