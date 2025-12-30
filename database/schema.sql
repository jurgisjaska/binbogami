SET FOREIGN_KEY_CHECKS = 0;

SET SESSION group_concat_max_len = 1000000;

SET @tables = NULL;
SELECT GROUP_CONCAT('`', table_name, '`') INTO @tables
FROM information_schema.tables
WHERE table_schema = 'binbogami';

SET @query = IF(@tables IS NOT NULL, CONCAT('DROP TABLE IF EXISTS ', @tables), 'SELECT 1');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE IF NOT EXISTS users
(
    id         CHAR(36)     NOT NULL
        PRIMARY KEY,
    email      VARCHAR(128) NOT NULL,
    name       VARCHAR(64)  NOT NULL,
    surname    VARCHAR(64)  NOT NULL,
    salt       CHAR(16)     NOT NULL,
    password   VARCHAR(256) NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at TIMESTAMP    NULL,
    CONSTRAINT users_email_uindex
        UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS books
(
    id              CHAR(36)     NOT NULL
        PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    description     TEXT         NULL,
    created_by      CHAR(36)     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    updated_at      TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP    NULL,
    closed_at       TIMESTAMP    NULL,
    CONSTRAINT books_name_uindex
        UNIQUE (name),
    CONSTRAINT books_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS categories
(
    id              CHAR(36)     NOT NULL
        PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    description     TEXT         NULL,
    organization_id CHAR(36)     NULL,
    created_by      CHAR(36)     NOT NULL,
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
    book_id     CHAR(36)  NOT NULL,
    category_id CHAR(36)  NOT NULL,
    created_by  CHAR(36)  NOT NULL,
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
    id              CHAR(36)     NOT NULL
        PRIMARY KEY,
    email           VARCHAR(128) NOT NULL,
    created_by      CHAR(36)     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    opened_at       TIMESTAMP    NULL,
    deleted_at      TIMESTAMP    NULL,
    expired_at      TIMESTAMP    NOT NULL,
    CONSTRAINT organizations_invitations_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS locations
(
    id              CHAR(36)     NOT NULL
        PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    description     TEXT         NULL,
    created_by      CHAR(36)     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    updated_at      TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP    NULL,
    CONSTRAINT locations_name_uindex
        UNIQUE (name),
    CONSTRAINT locations_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS books_locations
(
    id          INT AUTO_INCREMENT
        PRIMARY KEY,
    book_id     CHAR(36)  NOT NULL,
    location_id CHAR(36)  NOT NULL,
    created_by  CHAR(36)  NOT NULL,
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
    id          CHAR(36)  NOT NULL
        PRIMARY KEY,
    amount      FLOAT     NOT NULL,
    description TEXT      NULL,
    book_id     CHAR(36)  NOT NULL,
    category_id CHAR(36)  NOT NULL,
    location_id CHAR(36)  NOT NULL,
    created_by  CHAR(36)  NOT NULL,
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
    id            CHAR(36)  NOT NULL
        PRIMARY KEY,
    configuration INT       NOT NULL,
    value         TEXT      NULL,
    user_id       CHAR(36)  NOT NULL,
    created_at    TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP(),
    CONSTRAINT user_configurations_user_id_configuration_uindex
        UNIQUE (user_id, configuration),
    CONSTRAINT user_configuration_users_id_fk
        FOREIGN KEY (user_id) REFERENCES users (id)
);

