-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS matches (
    id uuid PRIMARY KEY,
    created_at timestamptz NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS matches;
-- +goose StatementEnd
