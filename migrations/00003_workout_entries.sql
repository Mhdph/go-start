-- +goose Up 
-- +goose StatementBegin 
CREATE TABLE IF NOT EXISTS workout_entries (
    id BIGSERIAL PRIMARY KEY,
    workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_name VARCHAR(255) NOT NULL,
    sets INT NOT NULL,
    reps INT NOT NULL,
    weight DECIMAL(10,2) NOT NULL,
    duration INT NOT NULL,
    notes TEXT,
    order_index INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    CONSTRAINT valid_workout_entry CHECK (sets > 0 AND reps > 0 AND weight >= 0 AND duration >= 0)
);

-- +goose StatementEnd 

-- +goose Down 
-- +goose StatementBegin 
DROP TABLE IF EXISTS workout_entries;
-- +goose StatementEnd 