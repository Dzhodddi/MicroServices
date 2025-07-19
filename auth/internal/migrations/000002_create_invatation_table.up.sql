CREATE TABLE IF NOT EXISTS user_inventations (
    token bytea primary key,
    user_id bigint not null,
    expiry timestamp(0) with time zone not null
);