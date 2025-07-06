-- +goose Up 
-- +goose StatementBegin
CREATE TABLE tokens (
    hash BYTEA PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    expiry TIMESTAMP(6) NOT NULL,
    scope TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;
-- +goose StatementEnd