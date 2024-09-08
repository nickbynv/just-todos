CREATE TABLE users (
    id            SERIAL       PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE lists (
    id          SERIAL       PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE users_lists (
    id      SERIAL                                     PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    list_id INT REFERENCES lists(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE items (
    id          SERIAL       PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    done        BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE lists_items (
    id      SERIAL                                     PRIMARY KEY,
    item_id INT REFERENCES items(id) ON DELETE CASCADE NOT NULL,
    list_id INT REFERENCES lists(id) ON DELETE CASCADE NOT NULL
);