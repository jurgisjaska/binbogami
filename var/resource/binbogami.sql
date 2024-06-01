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

CREATE TABLE IF NOT EXISTS organizations
(
    id          CHAR(36)    NOT NULL
        PRIMARY KEY,
    name        VARCHAR(64) NOT NULL,
    description TEXT        NULL,
    created_by  CHAR(36)    NOT NULL,
    created_at  TIMESTAMP   NOT NULL,
    updated_at  TIMESTAMP   NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at  TIMESTAMP   NULL,
    CONSTRAINT organizations_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS books
(
    id              CHAR(36)     NOT NULL
        PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    description     TEXT         NULL,
    organization_id CHAR(36)     NOT NULL,
    created_by      CHAR(36)     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    updated_at      TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP    NULL,
    closed_at       TIMESTAMP    NULL,
    CONSTRAINT books_organization_id_name_uindex
        UNIQUE (organization_id, name),
    CONSTRAINT books_organizations_id_fk
        FOREIGN KEY (organization_id) REFERENCES organizations (id),
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
    CONSTRAINT categories_organization_id_name_uindex
        UNIQUE (organization_id, name),
    CONSTRAINT categories_organizations_id_fk
        FOREIGN KEY (organization_id) REFERENCES organizations (id),
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
    organization_id CHAR(36)     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    opened_at       TIMESTAMP    NULL,
    deleted_at      TIMESTAMP    NULL,
    expired_at      TIMESTAMP    NOT NULL,
    CONSTRAINT organizations_invitations_organizations_id_fk
        FOREIGN KEY (organization_id) REFERENCES organizations (id),
    CONSTRAINT organizations_invitations_users_id_fk
        FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS locations
(
    id              CHAR(36)     NOT NULL
        PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    description     TEXT         NULL,
    organization_id CHAR(36)     NOT NULL,
    created_by      CHAR(36)     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    updated_at      TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP    NULL,
    CONSTRAINT locations_organization_id_name_uindex
        UNIQUE (organization_id, name),
    CONSTRAINT locations_organizations_id_fk
        FOREIGN KEY (organization_id) REFERENCES organizations (id),
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

CREATE TABLE IF NOT EXISTS members
(
    id              INT AUTO_INCREMENT
        PRIMARY KEY,
    role            TINYINT(2) DEFAULT 1 NOT NULL,
    organization_id CHAR(36)             NOT NULL,
    user_id         CHAR(36)             NOT NULL,
    created_by      CHAR(36)             NULL,
    created_at      TIMESTAMP            NOT NULL,
    updated_at      TIMESTAMP            NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP            NULL,
    CONSTRAINT members_organization_id_user_id_uindex
        UNIQUE (organization_id, user_id),
    CONSTRAINT members_organizations_id_fk
        FOREIGN KEY (organization_id) REFERENCES organizations (id),
    CONSTRAINT members_users_id_fk
        FOREIGN KEY (user_id) REFERENCES users (id)
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

