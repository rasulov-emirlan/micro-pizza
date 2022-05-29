-- +goose Up
-- +goose StatementBegin
BEGIN

CREATE TABLE IF NOT EXISTS roles (
    id          integer primary key generated always as identity,
    name        text not null,
    description text not null
);

CREATE TABLE IF NOT EXISTS users (
    id                  bigint primary key generated always as identity,
    full_name           text not null,
    email               text,
    phone_number        text,
    password            text,
    birth_date          date,
    country             text not null,
    city                text not null,
    street              text not null,
    floor               integer not null,
    address_instrutions text,
    created_at          timestamptz not null,
    updated_at          timestamptz
);

CREATE TABLE IF NOT EXISTS users_roles (
    user_id bigint not null,
    role_id integer not null,
    CONSTRAINT fk_users_roles_uid FOREIGN KEY(user_id)
        REFERENCES users(id),
    CONSTRAINT fk_users_roles_rid FOREIGN KEY(role_id)
        REFERENCES roles(id)
);

COMMIT;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
BEGIN

DROP TABLE IF EXISTS users_roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

DROP CONSTRAINT fk_users_roles_uid;
DROP CONSTRAINT fk_users_roles_rid;

COMMIT;
-- +goose StatementEnd