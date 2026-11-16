--liquibase formatted sql

--changeset openboard:20261115221449
--comment: Revoking all sessions of a User and purging expired sessions.
CREATE INDEX ix_sessions_user_id ON sessions (user_id);
CREATE INDEX ix_sessions_expires_at ON sessions (expires_at);

--rollback DROP INDEX ix_sessions_expires_at;
--rollback DROP INDEX ix_sessions_user_id;
