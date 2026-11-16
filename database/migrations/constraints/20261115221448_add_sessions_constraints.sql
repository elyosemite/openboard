--liquibase formatted sql

--changeset openboard:20261115221448
--comment: A session belongs to one User; its id is a 32-byte SHA-256 digest; it cannot expire before it starts.
ALTER TABLE sessions
    ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    ADD CONSTRAINT ck_sessions_id_length CHECK (octet_length(id) = 32),
    ADD CONSTRAINT ck_sessions_expiry_after_creation CHECK (expires_at > created_at);

--rollback ALTER TABLE sessions DROP CONSTRAINT ck_sessions_expiry_after_creation, DROP CONSTRAINT ck_sessions_id_length, DROP CONSTRAINT fk_sessions_user;
