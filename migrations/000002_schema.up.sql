CREATE SCHEMA IF NOT EXISTS posts;

CREATE TYPE posts.role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS posts.user (
    id              UUID            NOT NULL        DEFAULT uuid_generate_v4() PRIMARY KEY,
    name            VARCHAR(256)    NOT NULL        UNIQUE,
    role            posts.role      NOT NULL,
    created_at      TIMESTAMP       NOT NULL        DEFAULT CURRENT_TIMESTAMP
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts.session (
    id              UUID            NOT NULL        DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id         UUID            NOT NULL        REFERENCES posts.user(id),
    created_at      TIMESTAMP       NOT NULL        DEFAULT CURRENT_TIMESTAMP,
    expires_at      TIMESTAMP       NOT NULL        
);