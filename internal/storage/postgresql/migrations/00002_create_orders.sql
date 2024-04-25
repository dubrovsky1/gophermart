-- +goose Up
create table if not exists orders
(
    orderid     bigint      not null,
    userid      uuid        not null,
    uploaded_at timestamp   not null default current_timestamp,
    type        varchar(20) not null,
    status      varchar(20) not null,
    accrual     int         not null default 0
);

-- +goose Down
drop table if exists orders;