-- +goose Up
create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    login text not null,
    email text not null,
    password text not null,
    notification_methods text[],
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose Down
drop table if exists users;