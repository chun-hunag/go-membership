CREATE TABLE IF NOT EXISTS "user" (
    id serial NOT null PRIMARY KEY,
    name varchar(10) not null,
    email varchar(40) unique not null,
    password varchar(60) not null,
    remember_token varchar(100),
    created_at timestamp,
    updated_at timestamp
);

CREATE INDEX IF NOT EXISTS idx_user_name ON "user" (name);
CREATE INDEX IF NOT EXISTS idx_user_email ON "user" (email);