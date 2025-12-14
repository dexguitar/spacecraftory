-- +goose Up
create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    login text unique not null,
    email text unique not null,
    password text not null,
    notification_methods text[],
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists notification_methods (
    id uuid primary key default gen_random_uuid(),
    user_uuid uuid not null,
    provider_name text not null,
    target text not null,
    foreign key (user_uuid) references users(id) on delete cascade
);

-- +goose Down
drop table if exists notification_methods;
drop table if exists users;