CREATE TABLE IF NOT EXISTS todos (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    description text NOT NULL,
    done boolean NOT NULL DEFAULT FALSE,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    version integer NOT NULL DEFAULT 1
);