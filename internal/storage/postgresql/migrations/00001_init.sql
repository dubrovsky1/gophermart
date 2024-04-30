-- +goose Up
create table if not exists users
(
    userid      uuid         not null primary key,
    login       varchar(255) not null unique,
    password    varchar(255) not null,
    uploaded_at timestamp    not null default current_timestamp
);

create table if not exists orders
(
    orderid      bigint      not null,
    userid       uuid        not null,
    type         varchar(20) not null,
    status       varchar(20) not null,
    accrual      float       not null default 0,
    uploaded_at  timestamp   not null default current_timestamp,
    processed_at timestamp   not null default current_timestamp
);

-- +goose Down
drop table if exists users;
drop table if exists orders;