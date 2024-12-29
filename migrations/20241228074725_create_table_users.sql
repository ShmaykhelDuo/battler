-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    username varchar(64) NOT NULL UNIQUE,
    password_hash text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
