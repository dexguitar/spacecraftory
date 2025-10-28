-- +goose Up
create table orders (
    id uuid primary key default gen_random_uuid(),
    user_uuid text not null,
    part_uuids text[] not null,
    total_price decimal(10, 2) not null,
    status text not null default 'PENDING_PAYMENT',
    transaction_uuid text not null default '',
    payment_method text not null default 'UNKNOWN',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose Down
drop table orders;