CREATE TABLE IF NOT EXISTS users (
    id integer NOT null PRIMARY KEY,
    name varchar(10) not null,
    email varchar(40) unique not null,
    remember_token varchar(100) not null,
    created_at timestamp,
    updated_at timestamp
);

CREATE INDEX IF NOT EXISTS idx_users_name ON users (name);