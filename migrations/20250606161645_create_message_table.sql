-- +goose Up
create table messages (
    id bigserial primary key,
    chat_id bigint not null references chats(id),
    from_user varchar(64) not null,
    text text not null,
    timestamp timestamp not null
);

-- +goose Down
drop table messages;
