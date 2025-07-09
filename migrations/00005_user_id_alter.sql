-- +goose Up
ALTER TABLE workouts ADD COLUMN user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE workouts DROP COLUMN user_id;