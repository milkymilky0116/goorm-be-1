-- +goose Up
CREATE TABLE "user" (
    id bigserial PRIMARY KEY,
    email varchar NOT NULL UNIQUE,
    password varchar NOT NULL,
    role role NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP  
);
-- +goose Down
DROP TABLE "user";
