SET FOREIGN_KEY_CHECKS = 0;

SET SESSION group_concat_max_len = 1000000;

SET @tables = NULL;
SELECT GROUP_CONCAT('`', table_name, '`')
INTO @tables
FROM information_schema.tables
WHERE table_schema = 'binbogami';

SET @query = IF(@tables IS NOT NULL, CONCAT('DROP TABLE IF EXISTS ', @tables), 'SELECT 1');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET FOREIGN_KEY_CHECKS = 1;

create table if not exists users
(
    id           uuid         not null
        primary key,
    email        varchar(128) not null,
    name         varchar(64)  not null,
    surname      varchar(64)  not null,
    salt         char(16)     not null,
    password     varchar(256) not null,
    role         int          not null default 1,
    created_at   timestamp    not null,
    updated_at   timestamp    null on update current_timestamp(),
    confirmed_at timestamp    null,
    deleted_at   timestamp    null,
    constraint users_email_uindex
        unique (email)
);

CREATE TABLE IF NOT EXISTS books
(
    id          uuid         NOT NULL
        PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    description TEXT         NULL,
    created_by  UUID     NOT NULL,
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at  TIMESTAMP    NULL,
    closed_at   TIMESTAMP    NULL,
    CONSTRAINT books_name_uindex
        UNIQUE (name),
    CONSTRAINT books_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS categories
(
    id              uuid         NOT NULL
        PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    description     TEXT         NULL,
    created_by      UUID     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    updated_at      TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP    NULL,
    CONSTRAINT categories_name_uindex
        UNIQUE (name),
    CONSTRAINT categories_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS books_categories
(
    id          INT AUTO_INCREMENT
        PRIMARY KEY,
    book_id     UUID      NOT NULL,
    category_id UUID      NOT NULL,
    created_by  UUID      NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT books_categories_books_id_fk
        FOREIGN KEY (book_id) REFERENCES books (id),
    CONSTRAINT books_categories_categories_id_fk
        FOREIGN KEY (category_id) REFERENCES categories (id),
    CONSTRAINT books_categories_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS categories_name_index
    ON categories (name);

CREATE TABLE IF NOT EXISTS invitations
(
    id         UUID     NOT NULL
        PRIMARY KEY,
    email      VARCHAR(128) NOT NULL,
    created_by UUID     NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    opened_at  TIMESTAMP    NULL,
    deleted_at TIMESTAMP    NULL,
    expired_at TIMESTAMP    NOT NULL,
    CONSTRAINT organizations_invitations_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS locations
(
    id          UUID     NOT NULL
        PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    description TEXT         NULL,
    created_by  UUID     NOT NULL,
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at  TIMESTAMP    NULL,
    CONSTRAINT locations_name_uindex
        UNIQUE (name),
    CONSTRAINT locations_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS books_locations
(
    id          INT AUTO_INCREMENT
        PRIMARY KEY,
    book_id     UUID  NOT NULL,
    location_id UUID  NOT NULL,
    created_by  UUID  NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT books_locations_books_id_fk
        FOREIGN KEY (book_id) REFERENCES books (id),
    CONSTRAINT books_locations_locations_id_fk
        FOREIGN KEY (location_id) REFERENCES locations (id),
    CONSTRAINT books_locations_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS entries
(
    id          UUID  NOT NULL
        PRIMARY KEY,
    amount      FLOAT     NOT NULL,
    description TEXT      NULL,
    book_id     UUID  NOT NULL,
    category_id UUID  NOT NULL,
    location_id UUID  NOT NULL,
    created_by  UUID  NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT entries_books_id_fk
        FOREIGN KEY (book_id) REFERENCES books (id)
            ON DELETE CASCADE,
    CONSTRAINT entries_categories_id_fk
        FOREIGN KEY (category_id) REFERENCES categories (id),
    CONSTRAINT entries_locations_id_fk
        FOREIGN KEY (location_id) REFERENCES locations (id),
    CONSTRAINT entries_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS user_configurations
(
    id            UUID  NOT NULL
        PRIMARY KEY,
    configuration INT       NOT NULL,
    value         TEXT      NULL,
    user_id       UUID  NOT NULL,
    created_at    TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP(),
    CONSTRAINT user_configurations_user_id_configuration_uindex
        UNIQUE (user_id, configuration),
    CONSTRAINT user_configuration_users_id_fk
        FOREIGN KEY (user_id) REFERENCES users (id)
);

