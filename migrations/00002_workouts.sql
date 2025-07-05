-- +goose Up 
-- +goose StatementBegin 
CREATE TABLE IF NOT EXISTS workouts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT, 
    duration INT NOT NULL,
    calories INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

);

-- +goose StatementEnd 

-- +goose Down 
-- +goose StatementBegin 
DROP TABLE IF EXISTS workouts;
-- +goose StatementEnd 