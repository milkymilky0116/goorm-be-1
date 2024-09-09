-- +goose Up
CREATE TYPE role AS ENUM ('student', 'admin');
-- +goose Down
DROP TYPE role;
