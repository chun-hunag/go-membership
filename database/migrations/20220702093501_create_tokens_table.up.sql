CREATE TABLE IF NOT EXISTS tokens (
    id serial NOT null PRIMARY KEY,
    user_id int NOT null,
    token varchar(200),
    expired_at timestamp,
    created_at timestamp,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_tokens_token_expired_at ON tokens (token,expired_at);
