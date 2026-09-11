CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
