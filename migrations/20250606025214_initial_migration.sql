-- +goose Up
create table chats (
    id serial primary key,
    created_at timestamp not null default current_timestamp
);

create table chat_participants (
    chat_id bigint not null,
    username varchar(255) not null,
    joined_at timestamp not null default current_timestamp,
    primary key (chat_id, username),
    foreign key (chat_id) references chats(id) on delete cascade
);

-- +goose Down
drop table chats;
drop table chat_participants;
