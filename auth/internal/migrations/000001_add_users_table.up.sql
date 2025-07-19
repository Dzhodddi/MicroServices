CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    username varchar(255) not null,
    email varchar(255) not null unique,
    password text not null,
    created_at timestamp(0) with time zone not null default now(),
    role_id int2 not null,
    is_active boolean
)