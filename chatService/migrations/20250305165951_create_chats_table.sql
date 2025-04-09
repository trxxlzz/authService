-- +goose Up
CREATE TABLE chats
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages
(
    id        BIGSERIAL PRIMARY KEY,
    chat_id   BIGINT       NOT NULL,
    from_user VARCHAR(255) NOT NULL,
    text      TEXT         NOT NULL,
    timestamp TIMESTAMP    NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
);

CREATE TABLE chat_users
(
    chat_id  BIGINT       NOT NULL,
    username VARCHAR(255) NOT NULL,
    PRIMARY KEY (chat_id, username),
    FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chat_users;
DROP TABLE messages;
DROP TABLE chats;
