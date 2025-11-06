-- +goose Up
create table if not exists orders (
    id uuid primary key default gen_random_uuid(),
    user_uuid text not null,
    total_price decimal(10, 2) not null,
    status text not null,
    transaction_uuid text,
    payment_method text,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists order_parts (
    order_id uuid not null,
    part_id uuid not null,
    primary key (order_id, part_id),
    foreign key (order_id) references orders(id) on delete cascade
);

create index if not exists idx_order_parts_part_id on order_parts(part_id);

-- +goose Down
drop table if exists order_parts;
drop table if exists orders;