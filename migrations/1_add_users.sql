-- +migrate Up
CREATE TABLE IF NOT EXISTS users
(
    id         serial  NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    first_name varchar NOT NULL,
    last_name  varchar NOT NULL,
    email      varchar NOT NULL,
    phone      varchar NOT NULL,
    username   varchar NOT NULL
);

ALTER TABLE users
    OWNER TO "user";

CREATE UNIQUE INDEX IF NOT EXISTS users_email_uindex
    ON users (email);

CREATE UNIQUE INDEX IF NOT EXISTS users_first_name_uindex
    ON users (first_name);

CREATE UNIQUE INDEX IF NOT EXISTS users_last_name_uindex
    ON users (last_name);

CREATE UNIQUE INDEX IF NOT EXISTS users_phone_uindex
    ON users (phone);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_uindex
    ON users (username);

-- +migrate Down
DROP TABLE IF EXISTS users;
