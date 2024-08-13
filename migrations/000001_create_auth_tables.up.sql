CREATE TABLE IF NOT EXISTS users
(
    id    UUID PRIMARY KEY,
    email varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_sessions
(
    id            SERIAL PRIMARY KEY,
    user_id       UUID REFERENCES users (id) ON DELETE CASCADE,
    ip            VARCHAR(15),
    refresh_token VARCHAR(255)             NOT NULL,
    issued_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    expires_in    TIMESTAMP WITH TIME ZONE NOT NULL
);