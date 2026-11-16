--liquibase formatted sql

--changeset openboard:20261115221447
--comment: Username is unique in the Instance and always stored lowercase (3-32 of a-z 0-9 . _ -).
ALTER TABLE users
    ADD CONSTRAINT uq_users_username UNIQUE (username),
    ADD CONSTRAINT ck_users_username_format CHECK (username ~ '^[a-z0-9._-]{3,32}$'),
    ADD CONSTRAINT ck_users_name_not_blank CHECK (btrim(name) <> '');

--rollback ALTER TABLE users DROP CONSTRAINT ck_users_name_not_blank, DROP CONSTRAINT ck_users_username_format, DROP CONSTRAINT uq_users_username;
