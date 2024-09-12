-- +goose Up
CREATE TABLE refresh_token (
  id bigserial PRIMARY KEY,
  refresh_token TEXT NOT NULL,
  user_id bigserial REFERENCES "user"(id) ON DELETE CASCADE,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE refresh_token;

