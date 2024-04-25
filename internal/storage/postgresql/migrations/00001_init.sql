-- +goose Up
create table if not exists users
(
    userid   uuid         not null primary key,
    login    varchar(255) not null unique,
    password varchar(255) not null,
    created  timestamp    not null default current_timestamp
);

-- +goose Down
drop table if exists users;