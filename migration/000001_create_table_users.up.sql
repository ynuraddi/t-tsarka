create table if not exists users (
    id serial primary key,
    first_name varchar not null,
    last_name varchar not null
);