-- +goose Up
BEGIN;

-- Таблица чатов
CREATE TABLE chats (
    id          BIGSERIAL PRIMARY KEY,
    owner_id    BIGINT      NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Участники чатов
CREATE TABLE chat_participants (
    chat_id    BIGINT      NOT NULL,
    user_id    BIGINT      NOT NULL,
    role       SMALLINT    NOT NULL DEFAULT 0, -- 0:user, 1:admin, 2:owner
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id),
    CONSTRAINT fk_chat_participants_chat
        FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    CONSTRAINT chk_chat_participants_role
        CHECK (role IN (0,1,2))
);

-- Сообщения
CREATE TABLE messages (
    id         BIGSERIAL PRIMARY KEY,
    chat_id    BIGINT      NOT NULL,
    sender_id  BIGINT      NOT NULL,
    text       TEXT        NOT NULL,
    is_pinned  BOOLEAN     NOT NULL DEFAULT FALSE,
    send_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_messages_chat
        FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);

-- Индексы (практичные для типовых запросов)
CREATE INDEX idx_chats_owner_id ON chats(owner_id);

CREATE INDEX idx_chat_participants_user_id ON chat_participants(user_id);

-- для выборки сообщений по чату с сортировкой по времени
CREATE INDEX idx_messages_chat_id_send_at ON messages(chat_id, send_at DESC);

-- для выборок "сообщения пользователя"
CREATE INDEX idx_messages_sender_id ON messages(sender_id);

COMMIT;

-- +goose Down
BEGIN;

DROP INDEX IF EXISTS idx_messages_sender_id;
DROP INDEX IF EXISTS idx_messages_chat_id_send_at;
DROP INDEX IF EXISTS idx_chat_participants_user_id;
DROP INDEX IF EXISTS idx_chats_owner_id;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chat_participants;
DROP TABLE IF EXISTS chats;

COMMIT;
