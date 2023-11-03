CREATE TABLE IF NOT EXISTS users
(
    id         CHAR(36)  NOT NULL
        PRIMARY KEY,
    name       INT       NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS organizations
(
    id                 CHAR(36)     NOT NULL
        PRIMARY KEY,
    name               VARCHAR(128) NOT NULL,
    description        TEXT         NULL,
    created_by_user_id CHAR(36)     NOT NULL,
    owner_by_user_id   CHAR(36)     NOT NULL,
    created_at         TIMESTAMP    NOT NULL,
    updated_at         TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at         TIMESTAMP    NULL,
    CONSTRAINT organizations_users_id_fk
        FOREIGN KEY (created_by_user_id) REFERENCES users (id),
    CONSTRAINT organizations_users_id_fk2
        FOREIGN KEY (owner_by_user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS books
(
    id              CHAR(36)     NOT NULL
        PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    description     TEXT         NULL,
    organization_id CHAR(36)     NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    updated_at      TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP    NULL,
    CONSTRAINT books_organizations_id_fk
        FOREIGN KEY (organization_id) REFERENCES organizations (id)
);

CREATE TABLE IF NOT EXISTS categories
(
    id          CHAR(36)     NOT NULL
        PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    description TEXT         NULL,
    book_id     CHAR(36)     NOT NULL,
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at  TIMESTAMP    NULL,
    CONSTRAINT categories_books_id_fk
        FOREIGN KEY (book_id) REFERENCES books (id)
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS categories_name_index
    ON categories (name);

CREATE TABLE IF NOT EXISTS locations
(
    id         CHAR(36)     NOT NULL
        PRIMARY KEY,
    name       VARCHAR(128) NOT NULL,
    book_id    CHAR(36)     NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at TIMESTAMP    NULL,
    CONSTRAINT locations_books_id_fk
        FOREIGN KEY (book_id) REFERENCES books (id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS entries
(
    id                 CHAR(36)  NOT NULL
        PRIMARY KEY,
    amount             FLOAT     NOT NULL,
    description        TEXT      NULL,
    book_id            CHAR(36)  NOT NULL,
    category_id        CHAR(36)  NOT NULL,
    location_id        CHAR(36)  NOT NULL,
    created_by_user_id CHAR(36)  NOT NULL,
    created_at         TIMESTAMP NOT NULL,
    updated_at         TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at         TIMESTAMP NULL,
    CONSTRAINT entries_books_id_fk
        FOREIGN KEY (book_id) REFERENCES books (id)
            ON DELETE CASCADE,
    CONSTRAINT entries_categories_id_fk
        FOREIGN KEY (category_id) REFERENCES categories (id),
    CONSTRAINT entries_locations_id_fk
        FOREIGN KEY (location_id) REFERENCES locations (id),
    CONSTRAINT entries_users_id_fk
        FOREIGN KEY (created_by_user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS organizations_users
(
    id              INT AUTO_INCREMENT
        PRIMARY KEY,
    organization_id CHAR(36)  NOT NULL,
    user_id         CHAR(36)  NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP(),
    deleted_at      TIMESTAMP NULL,
    CONSTRAINT organizations_users_organizations_id_fk
        FOREIGN KEY (organization_id) REFERENCES organizations (id),
    CONSTRAINT organizations_users_users_id_fk
        FOREIGN KEY (user_id) REFERENCES users (id)
);

