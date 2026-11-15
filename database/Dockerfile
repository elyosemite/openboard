# Liquibase runner for OpenBoard migrations (local compose and CD).
# Liquibase 5 no longer bundles JDBC drivers, so the PostgreSQL driver is
# installed explicitly and pinned.
FROM liquibase/liquibase:5.0.2
RUN lpm add postgresql@42.7.8 --global
COPY --chown=liquibase:liquibase migrations /liquibase/changelog/migrations
